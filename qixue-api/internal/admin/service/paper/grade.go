package adminpaper

import (
	"context"
	"encoding/json"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"

	"qixue-api/internal/platform/bg"
	"qixue-api/internal/service/coin"
)

type QuestionScoreInput struct {
	QuestionID  int
	Score       *int
	GradeRemark *string
}

type GradeInput struct {
	UserPaperID    int
	QuestionScores []QuestionScoreInput
}

func (s *Service) GetUserPaperDetail(ctx context.Context, id int) (map[string]any, error) {
	userPaper, err := s.papers.FindUserPaper(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return nil, biz("记录不存在")
	}
	if err != nil {
		return nil, err
	}
	var paper paperRow
	err = s.db.WithContext(ctx).Where("id = ?", userPaper.PaperID).Take(&paper).Error
	if err == gorm.ErrRecordNotFound {
		return nil, biz("试卷不存在")
	}
	if err != nil {
		return nil, err
	}
	list, err := s.decorateUserPapers(ctx, []userPaperRow{{
		ID: userPaper.ID, UserID: userPaper.UserID, PaperID: userPaper.PaperID,
		Score: userPaper.Score, IsPass: userPaper.IsPass, IsPeopleCheck: userPaper.IsPeopleCheck,
		SubmitCount: userPaper.SubmitCount, CreatedAt: userPaper.CreatedAt,
	}}, true)
	if err != nil {
		return nil, err
	}
	var questions []questionRow
	if err := s.db.WithContext(ctx).Where("paper_id = ?", paper.ID).Order("sort ASC, id ASC").Find(&questions).Error; err != nil {
		return nil, err
	}
	logs, err := s.papers.ListQuestionLogsMap(ctx, id)
	if err != nil {
		return nil, err
	}
	logByQ := map[int64]map[string]any{}
	for _, log := range logs {
		qid := int64(asInt(log["question_id"]))
		logByQ[qid] = log
	}
	outQuestions := make([]map[string]any, 0, len(questions))
	for _, q := range questions {
		item := map[string]any{
			"id": q.ID, "paper_id": q.PaperID, "sort": q.Sort, "type": q.Type,
			"title": q.Title, "answer": decodeJSON(q.Answer), "score": q.Score, "options": decodeJSON(q.Options),
			"user_answer": logByQ[q.ID],
		}
		outQuestions = append(outQuestions, item)
	}
	return map[string]any{
		"user_paper": list[0],
		"paper":      paper,
		"questions":  outQuestions,
	}, nil
}

func (s *Service) PaperGrade(ctx context.Context, in GradeInput) error {
	if in.UserPaperID <= 0 {
		return biz("参数错误")
	}
	userPaperID := in.UserPaperID
	_, userID, paperID, _, wasPassRaw, err := s.papers.FindUserPaperForGrade(ctx, userPaperID)
	if err == gorm.ErrRecordNotFound {
		return biz("记录不存在")
	}
	if err != nil {
		return err
	}
	paper, err := s.papers.FindPaperForGrade(ctx, paperID)
	if err == gorm.ErrRecordNotFound {
		return biz("试卷不存在")
	}
	if err != nil {
		return err
	}
	scoreMap := map[int]struct {
		Score       *int
		GradeRemark *string
	}{}
	for _, item := range in.QuestionScores {
		scoreMap[item.QuestionID] = struct {
			Score       *int
			GradeRemark *string
		}{item.Score, item.GradeRemark}
	}
	logs, err := s.papers.ListQuestionLogs(ctx, userPaperID)
	if err != nil {
		return err
	}
	totalScore := 0
	for _, log := range logs {
		item, ok := scoreMap[log.QuestionID]
		updates := map[string]any{}
		if ok && item.Score != nil {
			score := *item.Score
			if log.Score != score {
				updates["score"] = score
				updates["is_correct"] = score > 0
				log.Score = score
			}
		}
		if ok && item.GradeRemark != nil {
			remark := strings.TrimSpace(*item.GradeRemark)
			if utf8.RuneCountInString(remark) > 500 {
				remark = string([]rune(remark)[:500])
			}
			if log.GradeRemark != remark {
				updates["grade_remark"] = remark
			}
		}
		if len(updates) > 0 {
			if err := s.papers.UpdateQuestionLog(ctx, log.ID, updates); err != nil {
				return err
			}
		}
		totalScore += log.Score
	}
	wasPassed := wasPassRaw == 1
	isPass := 0
	if totalScore >= paper.PassScore {
		isPass = 1
	}
	if err := s.papers.UpdateUserPaperGrade(ctx, userPaperID, totalScore, isPass); err != nil {
		return err
	}
	if isPass == 1 && !wasPassed {
		if paper.ChapterID != nil && *paper.ChapterID > 0 {
			userChapter, err := s.papers.FindUserChapter(ctx, userID, *paper.ChapterID)
			if err == nil {
				_ = s.papers.MarkUserChapterEnd(ctx, userChapter.ID)
				if userChapter.CourseID != nil {
					_ = s.papers.IncHasEndChapter(ctx, userID, *userChapter.CourseID)
				}
				chapterID, paperID := *paper.ChapterID, paper.ID
				bg.Go(ctx, func(ctx context.Context) {
					chapterTitle, _ := s.papers.ChapterTitle(ctx, chapterID)
					coin.OnPassTest(ctx, s.db, s.rdb, userID, paperID, chapterTitle)
				})
			}
		} else if paper.CourseID != nil && *paper.CourseID > 0 {
			courseID := *paper.CourseID
			_ = s.papers.MarkUserCourseEnd(ctx, userID, courseID)
			bg.Go(ctx, func(ctx context.Context) {
				course, err := s.papers.FindCourseCredit(ctx, courseID)
				if err != nil {
					return
				}
				if course.Credit > 0 {
					_ = coin.AddCredit(ctx, s.db, userID, "课程发放", course.Credit, course.Title+"课程学习结束发放")
				}
				coin.OnCourseEnd(ctx, s.db, s.rdb, userID, courseID, course.Title)
			})
		}
	}
	return nil
}

func decodeJSON(raw []byte) any {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return string(raw)
	}
	return value
}
