package ops

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/admin/ops"
)

type Handler struct{ svc *adminops.Service }

func New(svc *adminops.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, adminops.ErrBadScene) {
		response.Fail(c, 1, err.Error())
		return true
	}
	return response.WriteErr(c, err)
}

func (h *Handler) Chart(c *gin.Context) {
	var body struct {
		Scene string `json:"scene"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.Chart(c.Request.Context(), middleware.UID(c), body.Scene)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Handler) CourseEnrollList(c *gin.Context) {
	q, ok := h.bindEnrollList(c)
	if !ok {
		return
	}
	data, err := h.svc.CourseEnrollList(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Handler) PathEnrollList(c *gin.Context) {
	q, ok := h.bindEnrollList(c)
	if !ok {
		return
	}
	data, err := h.svc.PathEnrollList(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Handler) LearnList(c *gin.Context) {
	q, ok := h.bindEnrollList(c)
	if !ok {
		return
	}
	data, err := h.svc.LearnList(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Handler) CreditList(c *gin.Context) {
	q, ok := h.bindAmountList(c)
	if !ok {
		return
	}
	data, err := h.svc.CreditList(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Handler) CoinsList(c *gin.Context) {
	q, ok := h.bindAmountList(c)
	if !ok {
		return
	}
	data, err := h.svc.CoinsList(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Handler) bindEnrollList(c *gin.Context) (adminops.EnrollListQuery, bool) {
	var body struct {
		NickName string   `json:"nick_name"`
		Title    string   `json:"title"`
		IsEnd    any      `json:"is_end"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	q := adminops.EnrollListQuery{
		NickName: body.NickName, Title: body.Title,
		Page: body.Page.Int(), Size: body.Size.Int(),
	}
	if bind.Present(body.IsEnd) {
		n := bind.AsInt(body.IsEnd)
		q.IsEnd = &n
	}
	return q, true
}

func (h *Handler) bindAmountList(c *gin.Context) (adminops.AmountListQuery, bool) {
	var body struct {
		NickName string   `json:"nick_name"`
		Field    string   `json:"field"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	return adminops.AmountListQuery{
		NickName: body.NickName, Field: body.Field,
		Page: body.Page.Int(), Size: body.Size.Int(),
	}, true
}
