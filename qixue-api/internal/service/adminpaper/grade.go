package adminpaper

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/coin"
)

func (s *Service) GetUserPaperDetail(c *gin.Context) {
	id := bind.Query(c, "id")
	ctx := c.Request.Context()
	userPaper, err := s.papers.FindUserPaper(ctx, id)
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "记录不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	var paper paperRow
	err = s.db.WithContext(ctx).Where("id = ?", userPaper.PaperID).Take(&paper).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "试卷不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	list, err := s.decorateUserPapers(c, []userPaperRow{{
		ID: userPaper.ID, UserID: userPaper.UserID, PaperID: userPaper.PaperID,
		Score: userPaper.Score, IsPass: userPaper.IsPass, IsPeopleCheck: userPaper.IsPeopleCheck,
		SubmitCount: userPaper.SubmitCount, CreatedAt: userPaper.CreatedAt,
	}}, true)
	if response.WriteErr(c, err) {
		return
	}
	var questions []questionRow
	if response.WriteErr(c, s.db.WithContext(ctx).Where("paper_id = ?", paper.ID).Order("sort ASC, id ASC").Find(&questions).Error) {
		return
	}
	logs, err := s.papers.ListQuestionLogsMap(ctx, id)
	if response.WriteErr(c, err) {
		return
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
	response.OK(c, "", gin.H{
		"user_paper": list[0],
		"paper":      paper,
		"questions":  outQuestions,
	})
}

func (s *Service) PaperGrade(c *gin.Context) {
	var body struct {
		UserPaperID    bind.Int `json:"user_paper_id"`
		QuestionScores []struct {
			QuestionID  bind.Int `json:"question_id"`
			Score       *int     `json:"score"`
			GradeRemark *string  `json:"grade_remark"`
		} `json:"question_scores"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.UserPaperID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	userPaperID := body.UserPaperID.Int()
	ctx := c.Request.Context()
	_, userID, paperID, _, wasPassRaw, err := s.papers.FindUserPaperForGrade(ctx, userPaperID)
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "记录不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	paper, err := s.papers.FindPaperForGrade(ctx, paperID)
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "试卷不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	scoreMap := map[int]struct {
		Score       *int
		GradeRemark *string
	}{}
	for _, item := range body.QuestionScores {
		scoreMap[item.QuestionID.Int()] = struct {
			Score       *int
			GradeRemark *string
		}{item.Score, item.GradeRemark}
	}
	logs, err := s.papers.ListQuestionLogs(ctx, userPaperID)
	if response.WriteErr(c, err) {
		return
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
			if response.WriteErr(c, s.papers.UpdateQuestionLog(ctx, log.ID, updates)) {
				return
			}
		}
		totalScore += log.Score
	}
	wasPassed := wasPassRaw == 1
	isPass := 0
	if totalScore >= paper.PassScore {
		isPass = 1
	}
	if response.WriteErr(c, s.papers.UpdateUserPaperGrade(ctx, userPaperID, totalScore, isPass)) {
		return
	}
	if isPass == 1 && !wasPassed {
		if paper.ChapterID != nil && *paper.ChapterID > 0 {
			userChapter, err := s.papers.FindUserChapter(ctx, userID, *paper.ChapterID)
			if err == nil {
				_ = s.papers.MarkUserChapterEnd(ctx, userChapter.ID)
				if userChapter.CourseID != nil {
					_ = s.papers.IncHasEndChapter(ctx, userID, *userChapter.CourseID)
				}
				chapterTitle, _ := s.papers.ChapterTitle(ctx, *paper.ChapterID)
				coin.OnPassTest(ctx, s.db, s.rdb, userID, paper.ID, chapterTitle)
			}
		} else if paper.CourseID != nil && *paper.CourseID > 0 {
			courseID := *paper.CourseID
			_ = s.papers.MarkUserCourseEnd(ctx, userID, courseID)
			course, err := s.papers.FindCourseCredit(ctx, courseID)
			if err == nil {
				if course.Credit > 0 {
					_ = coin.AddCredit(ctx, s.db, userID, "课程发放", course.Credit, course.Title+"课程学习结束发放")
				}
				coin.OnCourseEnd(ctx, s.db, s.rdb, userID, courseID, course.Title)
			}
		}
	}
	response.OK(c, "批改成功", []any{})
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
