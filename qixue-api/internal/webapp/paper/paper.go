package paper

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	papersvc "qixue-api/internal/service/webapp/paper"
)

type Handler struct {
	svc *papersvc.Service
}

func New(svc *papersvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, papersvc.ErrAlreadyPassedOK):
		return false
	case errors.Is(err, papersvc.ErrPaperNotFound),
		errors.Is(err, papersvc.ErrAlreadyPassed),
		errors.Is(err, papersvc.ErrLogNotFound),
		errors.Is(err, papersvc.ErrNoPermission),
		errors.Is(err, papersvc.ErrBadParam):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) GetPaperInfo(c *gin.Context) {
	data, passed, err := h.svc.GetPaperInfo(c.Request.Context(), middleware.UID(c), bind.Query(c, "id"))
	if errors.Is(err, papersvc.ErrAlreadyPassedOK) && passed != nil {
		response.OK(c, err.Error(), gin.H{
			"user_paper_id": passed.UserPaperID, "paper": passed.Paper, "is_pass": true,
		})
		return
	}
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "paper", data)
}

func (h *Handler) GetPaperLogInfo(c *gin.Context) {
	data, err := h.svc.GetPaperLogInfo(c.Request.Context(), middleware.UID(c), bind.Query(c, "id"))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "paper", data)
}

func (h *Handler) SubmitPaper(c *gin.Context) {
	var body struct {
		ID      bind.Int `json:"id"`
		Answers []struct {
			QuestionID bind.Int `json:"question_id"`
			Answer     any      `json:"answer"`
		} `json:"answers"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	answers := make([]papersvc.AnswerIn, 0, len(body.Answers))
	for _, a := range body.Answers {
		answers = append(answers, papersvc.AnswerIn{QuestionID: a.QuestionID.Int(), Answer: a.Answer})
	}
	data, err := h.svc.SubmitPaper(c.Request.Context(), middleware.UID(c), body.ID.Int(), answers)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "提交成功", data)
}
