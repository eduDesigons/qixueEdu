package paper

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/async"
	"qixue-api/internal/platform/bg"
	webapprepo "qixue-api/internal/webapp/repo"
	"qixue-api/internal/service/coin"
	"qixue-api/internal/service/enroll"
	"qixue-api/internal/service/textutil"
)

var (
	ErrPaperNotFound   = errors.New("不存在的试卷")
	ErrAlreadyPassed   = errors.New("您已经通过过此试卷了")
	ErrAlreadyPassedOK = errors.New("您已经通过过此试卷了,无需再次测试")
	ErrLogNotFound     = errors.New("不存在的记录")
	ErrNoPermission    = errors.New("您没有权限查看此记录")
	ErrBadParam        = errors.New("参数错误")
)

type Service struct {
	papers webapprepo.PaperRepository
	db     *gorm.DB
	rdb    *redis.Client
}

func New(papers webapprepo.PaperRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{papers: papers, db: db, rdb: rdb}
}

type AnswerIn struct {
	QuestionID int
	Answer     any
}

// AlreadyPassedInfo is returned when GetPaperInfo finds a prior pass.
type AlreadyPassedInfo struct {
	UserPaperID int
	Paper       map[string]any
}

func (s *Service) GetPaperInfo(ctx context.Context, uid, paperID int) (map[string]any, *AlreadyPassedInfo, error) {
	paper, err := s.findPaper(ctx, paperID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, ErrPaperNotFound
	}
	if err != nil {
		return nil, nil, err
	}
	var (
		passedID int
		last     *webapprepo.UserPaperAttempt
	)
	if err := async.Run(
		func() error {
			var callErr error
			passedID, callErr = s.papers.FindPassedID(ctx, uid, paperID)
			return async.SkipNotFound(callErr)
		},
		func() error {
			var callErr error
			last, callErr = s.papers.FindLatestAttempt(ctx, uid, paperID)
			return async.SkipNotFound(callErr)
		},
	); err != nil {
		return nil, nil, err
	}
	if passedID > 0 {
		return nil, &AlreadyPassedInfo{UserPaperID: passedID, Paper: paper}, ErrAlreadyPassedOK
	}
	isPass, isPeopleCheck := false, false
	var userPaperID any
	if last != nil {
		isPass = last.IsPass == 1
		isPeopleCheck = last.IsPeopleCheck == 1
		userPaperID = last.ID
	}
	questions, err := s.loadQuestionsForTake(ctx, paperID)
	if err != nil {
		return nil, nil, err
	}
	return map[string]any{
		"paper": paper, "questions": questions,
		"is_pass": isPass, "is_people_check": isPeopleCheck, "user_paper_id": userPaperID,
	}, nil, nil
}

func (s *Service) GetPaperLogInfo(ctx context.Context, uid, logID int) (map[string]any, error) {
	userPassLog, err := s.papers.FindUserPaper(ctx, logID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLogNotFound
	}
	if err != nil {
		return nil, err
	}
	if userPassLog.UserID != uid {
		return nil, ErrNoPermission
	}
	var (
		paper     map[string]any
		hasPassed int64
		questions []map[string]any
	)
	if err := async.Run(
		func() error {
			var callErr error
			paper, callErr = s.findPaper(ctx, userPassLog.PaperID)
			if errors.Is(callErr, gorm.ErrRecordNotFound) {
				return ErrPaperNotFound
			}
			return callErr
		},
		func() error {
			hasPassed, _ = s.papers.CountPassed(ctx, uid, userPassLog.PaperID)
			return nil
		},
		func() error {
			var callErr error
			questions, callErr = s.loadQuestionsForLog(ctx, userPassLog.PaperID, userPassLog.ID, userPassLog.IsPass == 1)
			return callErr
		},
	); err != nil {
		return nil, err
	}
	return map[string]any{
		"paper": paper, "questions": questions,
		"is_pass": userPassLog.IsPass == 1, "has_passed": hasPassed > 0,
		"is_people_check": userPassLog.IsPeopleCheck == 1,
	}, nil
}

func (s *Service) SubmitPaper(ctx context.Context, uid, paperID int, answers []AnswerIn) (map[string]any, error) {
	if paperID <= 0 {
		return nil, ErrBadParam
	}
	passed, err := s.papers.CountPassed(ctx, uid, paperID)
	if err != nil {
		return nil, err
	}
	if passed > 0 {
		return nil, ErrAlreadyPassed
	}
	paperMeta, err := s.papers.FindMeta(ctx, paperID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPaperNotFound
	}
	if err != nil {
		return nil, err
	}
	questions, err := s.papers.ListQuestions(ctx, paperID)
	if err != nil {
		return nil, err
	}
	qMap := map[int64]struct {
		ID     int64
		Type   string
		Answer any
		Score  int
	}{}
	for _, q := range questions {
		qMap[q.ID] = struct {
			ID     int64
			Type   string
			Answer any
			Score  int
		}{q.ID, q.Type, decodeJSON(q.Answer), q.Score}
	}
	now := time.Now()
	userPaperID, err := s.papers.CreateUserPaper(ctx, uid, paperID, now)
	if err != nil {
		return nil, err
	}
	needPeopleCheck := false
	totalScore := 0
	logs := make([]map[string]any, 0)
	logRows := make([]map[string]any, 0)
	for _, ans := range answers {
		q, ok := qMap[int64(ans.QuestionID)]
		if !ok {
			continue
		}
		score := 0
		needPeople := false
		isCorrect := false
		if !isEmptyAnswer(ans.Answer) {
			switch q.Type {
			case "single", "multiple", "completion":
				isCorrect = answersEqual(q.Answer, ans.Answer)
				if isCorrect {
					score = q.Score
				}
			case "text":
				sim := textutil.AverageSimilarity(q.Answer, ans.Answer)
				if sim >= 40 {
					isCorrect = true
					score = q.Score
				} else if sim > 2 {
					isCorrect = true
					score = int(sim*float64(q.Score)/100) * 2
					if score <= 0 {
						score = 1
					}
				}
			default:
				needPeopleCheck = true
				needPeople = true
			}
		}
		answerJSON, _ := json.Marshal(ans.Answer)
		correctFlag := 0
		if isCorrect {
			correctFlag = 1
		}
		peopleFlag := 0
		if needPeople {
			peopleFlag = 1
		}
		logRows = append(logRows, map[string]any{
			"user_id": uid, "user_paper_id": userPaperID, "question_id": q.ID,
			"answer_content": string(answerJSON), "score": score,
			"is_correct": correctFlag, "need_people_check": peopleFlag,
			"created_at": now, "updated_at": now,
		})
		totalScore += score
		logs = append(logs, map[string]any{
			"question_id": q.ID, "score": score, "is_correct": isCorrect,
		})
	}
	if err := s.papers.CreateQuestionLogs(ctx, logRows); err != nil {
		return nil, err
	}
	isPass := 0
	isPeople := 0
	if needPeopleCheck {
		isPeople = 1
	}
	if totalScore >= paperMeta.PassScore {
		isPass = 1
		isPeople = 0
	}
	if err := s.papers.UpdateUserPaper(ctx, userPaperID, map[string]any{
		"score": totalScore, "is_pass": isPass, "is_people_check": isPeople, "updated_at": now,
	}); err != nil {
		return nil, err
	}
	if isPass == 1 {
		s.onPass(ctx, uid, paperID, paperMeta)
	}
	return map[string]any{
		"id": userPaperID, "user_id": uid, "paper_id": paperID,
		"score": totalScore, "is_pass": isPass == 1, "is_people_check": isPeople == 1,
		"created_at": model.FormatTime(&now), "updated_at": model.FormatTime(&now),
		"logs": logs,
	}, nil
}

func (s *Service) onPass(ctx context.Context, uid, paperID int, paperMeta *webapprepo.PaperMeta) {
	if paperMeta.ChapterID != nil && *paperMeta.ChapterID > 0 {
		userChapter, err := s.papers.FindUserChapter(ctx, uid, *paperMeta.ChapterID)
		if err != nil {
			return
		}
		_ = s.papers.MarkChapterEnded(ctx, userChapter.ID)
		if userChapter.CourseID != nil {
			courseID := *userChapter.CourseID
			_ = s.papers.IncrCourseHasEnd(ctx, uid, courseID)
			enroll.FinishCourseIfReady(ctx, s.db, s.rdb, uid, courseID)
		}
		enroll.ForgetListCache(ctx, s.rdb, uid)
		chapterID := *paperMeta.ChapterID
		bg.Go(ctx, func(ctx context.Context) {
			chapterTitle, _ := s.papers.ChapterTitle(ctx, chapterID)
			if chapterTitle == "" {
				chapterTitle = "课节测评"
			}
			coin.OnPassTest(ctx, s.db, s.rdb, uid, paperID, chapterTitle)
		})
		return
	}
	if paperMeta.CourseID != nil && *paperMeta.CourseID > 0 {
		courseID := *paperMeta.CourseID
		_ = s.papers.MarkCourseEnded(ctx, uid, courseID)
		enroll.ForgetListCache(ctx, s.rdb, uid)
		title, credit, _ := s.papers.FindCourseCredit(ctx, courseID)
		enroll.GrantCourseEnd(ctx, s.db, s.rdb, uid, courseID, title, credit)
	}
}

func (s *Service) findPaper(ctx context.Context, id int) (map[string]any, error) {
	meta, err := s.papers.FindMeta(ctx, id)
	if err != nil {
		return nil, err
	}
	var courseID int
	var courseTitle string
	var chapterID int
	var chapterTitle string
	if meta.CourseID != nil && *meta.CourseID > 0 {
		courseID, courseTitle, _ = s.papers.FindCourseBrief(ctx, *meta.CourseID)
	}
	if meta.ChapterID != nil && *meta.ChapterID > 0 {
		chapterID, chapterTitle, _ = s.papers.FindChapterBrief(ctx, *meta.ChapterID)
	}
	title := ""
	switch {
	case chapterTitle != "" && courseTitle != "":
		title = courseTitle + " · " + chapterTitle
	case chapterTitle != "":
		title = chapterTitle
	case courseTitle != "":
		title = courseTitle
	case meta.Type == "course":
		title = "结课考试"
	default:
		title = "章节测评"
	}
	out := map[string]any{
		"id": meta.ID, "score": meta.Score, "pass_score": meta.PassScore,
		"type": meta.Type, "course_id": meta.CourseID, "chapter_id": meta.ChapterID, "title": title,
	}
	if courseID > 0 {
		out["course"] = map[string]any{"id": courseID, "title": courseTitle}
	}
	if chapterID > 0 {
		out["chapter"] = map[string]any{"id": chapterID, "title": chapterTitle}
	}
	return out, nil
}

func (s *Service) loadQuestionsForTake(ctx context.Context, paperID int) ([]map[string]any, error) {
	rows, err := s.papers.ListQuestions(ctx, paperID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, q := range rows {
		item := map[string]any{
			"id": q.ID, "paper_id": q.PaperID, "sort": q.Sort, "type": q.Type,
			"title": q.Title, "score": q.Score, "options": decodeJSON(q.Options),
			"created_at": model.FormatTime(q.CreatedAt), "updated_at": model.FormatTime(q.UpdatedAt),
		}
		answer := decodeJSON(q.Answer)
		if q.Type == "completion" {
			keyWords, answerCount := completionOptions(answer, q.Title)
			item["options"] = keyWords
			item["answer_count"] = answerCount
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *Service) loadQuestionsForLog(ctx context.Context, paperID, userPaperID int, showAnswer bool) ([]map[string]any, error) {
	rows, err := s.papers.ListQuestions(ctx, paperID)
	if err != nil {
		return nil, err
	}
	logs, _ := s.papers.ListQuestionLogs(ctx, userPaperID)
	byQ := map[int64]map[string]any{}
	for _, log := range logs {
		qid := int64(bind.AsInt(log["question_id"]))
		if _, ok := byQ[qid]; ok {
			continue
		}
		formatLogTimes(log)
		byQ[qid] = log
	}
	out := make([]map[string]any, 0, len(rows))
	for _, q := range rows {
		item := map[string]any{
			"id": q.ID, "paper_id": q.PaperID, "sort": q.Sort, "type": q.Type,
			"title": q.Title, "score": q.Score, "options": decodeJSON(q.Options),
			"created_at": model.FormatTime(q.CreatedAt), "updated_at": model.FormatTime(q.UpdatedAt),
		}
		if ua, ok := byQ[q.ID]; ok {
			item["user_answer"] = ua
		}
		if showAnswer {
			item["answer"] = decodeJSON(q.Answer)
		}
		out = append(out, item)
	}
	return out, nil
}
