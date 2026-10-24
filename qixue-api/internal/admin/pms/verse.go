package pms

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/adminverse"
)

type Verse struct{ svc *adminverse.Service }

func NewVerse(svc *adminverse.Service) *Verse { return &Verse{svc: svc} }

func (h *Verse) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminverse.ErrEmpty),
		errors.Is(err, adminverse.ErrTooLong),
		errors.Is(err, adminverse.ErrBadID),
		errors.Is(err, adminverse.ErrMissing):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Verse) List(c *gin.Context) {
	var body struct {
		Keyword string `json:"keyword"`
		Page    int    `json:"page"`
		Size    int    `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page, body.Size)
	list, total, err := h.svc.List(c.Request.Context(), body.Keyword, current, size)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (h *Verse) Save(c *gin.Context) {
	var body struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "请填写一言")
		return
	}
	id, err := h.svc.Save(c.Request.Context(), body.ID, body.Content)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "保存成功", gin.H{"id": id})
}

func (h *Verse) Delete(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), body.ID)) {
		return
	}
	response.OK(c, "删除成功", []any{})
}
