package feedback

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	feedbacksvc "qixue-api/internal/service/feedback"
)

type Handler struct {
	svc *feedbacksvc.Service
}

func New(svc *feedbacksvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, feedbacksvc.ErrEmptyContent),
		errors.Is(err, feedbacksvc.ErrTooLong),
		errors.Is(err, feedbacksvc.ErrTooManyImgs),
		errors.Is(err, feedbacksvc.ErrDailyLimit):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) Submit(c *gin.Context) {
	var body struct {
		Content string `json:"content"`
		Images  any    `json:"images"`
	}
	_ = c.ShouldBindJSON(&body)
	id, err := h.svc.Submit(c.Request.Context(), middleware.UID(c), body.Content, body.Images)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "提交成功", gin.H{"id": id})
}

func (h *Handler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", list)
}
