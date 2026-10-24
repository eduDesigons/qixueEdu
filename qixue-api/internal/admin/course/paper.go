package course

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	adminpaper "qixue-api/internal/service/admin/paper"
)

type Paper struct{ svc *adminpaper.Service }

func NewPaper(svc *adminpaper.Service) *Paper {
	return &Paper{svc: svc}
}

func (h *Paper) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var be *adminpaper.BizError
	if errors.As(err, &be) {
		response.Fail(c, 1, be.Message)
		return true
	}
	return response.WriteErr(c, err)
}

func (h *Paper) GetUserPaperDetail(c *gin.Context) {
	data, err := h.svc.GetUserPaperDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Paper) PaperGrade(c *gin.Context) {
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
	in := adminpaper.GradeInput{UserPaperID: body.UserPaperID.Int()}
	for _, item := range body.QuestionScores {
		in.QuestionScores = append(in.QuestionScores, adminpaper.QuestionScoreInput{
			QuestionID: item.QuestionID.Int(), Score: item.Score, GradeRemark: item.GradeRemark,
		})
	}
	err := h.svc.PaperGrade(c.Request.Context(), in)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "批改成功", []any{})
}

func (h *Paper) GetPaperDetail(c *gin.Context) {
	paper, err := h.svc.GetPaperDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", paper)
}

func (h *Paper) PaperSave(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	msg, err := h.svc.PaperSave(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, []any{})
}

func (h *Paper) PaperDelete(c *gin.Context) {
	err := h.svc.PaperDelete(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "成功", []any{})
}

func (h *Paper) GetExamList(c *gin.Context) {
	q, ok := h.bindUserPaperList(c)
	if !ok {
		return
	}
	data, err := h.svc.GetExamList(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Paper) GetAssessmentList(c *gin.Context) {
	q, ok := h.bindUserPaperList(c)
	if !ok {
		return
	}
	data, err := h.svc.GetAssessmentList(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Paper) bindUserPaperList(c *gin.Context) (adminpaper.UserPaperListQuery, bool) {
	var body struct {
		CourseID      bind.Int `json:"course_id"`
		ChapterID     bind.Int `json:"chapter_id"`
		Title         string   `json:"title"`
		IsPass        any      `json:"is_pass"`
		IsPeopleCheck any      `json:"is_people_check"`
		NickName      string   `json:"nick_name"`
		Sort          string   `json:"sort"`
		Order         string   `json:"order"`
		Page          bind.Int `json:"page"`
		Size          bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	return adminpaper.UserPaperListQuery{
		CourseID: body.CourseID.Int(), ChapterID: body.ChapterID.Int(),
		Title: body.Title, IsPass: body.IsPass, IsPeopleCheck: body.IsPeopleCheck,
		NickName: body.NickName, Sort: body.Sort, Order: body.Order,
		Page: body.Page.Int(), Size: body.Size.Int(),
	}, true
}

func (h *Paper) GetPaperLogs(c *gin.Context) {
	var body struct {
		PaperID  bind.Int `json:"paper_id"`
		NickName string   `json:"nick_name"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.GetPaperLogs(c.Request.Context(), middleware.UID(c), adminpaper.PaperLogsQuery{
		PaperID: body.PaperID.Int(), NickName: body.NickName,
		Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}
