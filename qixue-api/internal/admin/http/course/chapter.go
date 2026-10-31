package course

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	adminchapter "qixue-api/internal/admin/service/chapter"
)

type Chapter struct{ svc *adminchapter.Service }

func NewChapter(svc *adminchapter.Service) *Chapter {
	return &Chapter{svc: svc}
}

func (h *Chapter) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var be *adminchapter.BizError
	if errors.As(err, &be) {
		response.Fail(c, 1, be.Message)
		return true
	}
	return response.WriteErr(c, err)
}

func (h *Chapter) bindRaw(c *gin.Context) (map[string]any, bool) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return nil, false
	}
	if raw == nil {
		raw = map[string]any{}
	}
	return raw, true
}

func (h *Chapter) GetChapterList(c *gin.Context) {
	list, err := h.svc.GetChapterList(c.Request.Context(), bind.Query(c, "courseId", "course_id"), c.Query("title"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Chapter) GetChapterDetail(c *gin.Context) {
	data, err := h.svc.GetChapterDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Chapter) ChapterSave(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	msg, err := h.svc.ChapterSave(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, []any{})
}

func (h *Chapter) ChapterDelete(c *gin.Context) {
	err := h.svc.ChapterDelete(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Chapter) ChapterSort(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	err := h.svc.ChapterSort(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "成功", []any{})
}

func (h *Chapter) GetChapterUserList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetChapterUserList(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Chapter) UnlockChapterUser(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	err := h.svc.UnlockChapterUser(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "解锁成功", []any{})
}
