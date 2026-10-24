package member

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/admin/feedback"
)

type Feedback struct{ svc *adminfeedback.Service }

func NewFeedback(svc *adminfeedback.Service) *Feedback {
	return &Feedback{svc: svc}
}

func (h *Feedback) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminfeedback.ErrNotFound),
		errors.Is(err, adminfeedback.ErrEmptyReply),
		errors.Is(err, adminfeedback.ErrReplyTooLong):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Feedback) List(c *gin.Context) {
	var body struct {
		NickName string   `json:"nick_name"`
		Status   any      `json:"status"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	q := adminfeedback.ListQuery{
		NickName: body.NickName,
		Page:     body.Page.Int(),
		Size:     body.Size.Int(),
	}
	if bind.Present(body.Status) {
		n := bind.AsInt(body.Status)
		q.Status = &n
	}
	data, err := h.svc.List(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Feedback) Reply(c *gin.Context) {
	var body struct {
		ID    bind.Int `json:"id"`
		Reply string   `json:"reply"`
	}
	_ = c.ShouldBindJSON(&body)
	err := h.svc.Reply(c.Request.Context(), middleware.UID(c), body.ID.Int(), body.Reply)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "回复成功", nil)
}
