package course

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	adminpath "qixue-api/internal/admin/service/path"
)

type Path struct{ svc *adminpath.Service }

func NewPath(svc *adminpath.Service) *Path {
	return &Path{svc: svc}
}

func (h *Path) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var be *adminpath.BizError
	if errors.As(err, &be) {
		response.Fail(c, 1, be.Message)
		return true
	}
	return response.WriteErr(c, err)
}

func (h *Path) bindRaw(c *gin.Context) (map[string]any, bool) {
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

func (h *Path) GetPathCourseSearch(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetPathCourseSearch(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Path) GetPathCourseList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	list, err := h.svc.GetPathCourseList(c.Request.Context(), bind.RequestInt(c, raw, "pathId", "path_id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Path) PathCourseListSave(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	err := h.svc.PathCourseListSave(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "课节创建成功", []any{})
}

func (h *Path) PathCourseDelete(c *gin.Context) {
	err := h.svc.PathCourseDelete(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Path) GetPathCategories(c *gin.Context) {
	list, err := h.svc.GetPathCategories(c.Request.Context())
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Path) GetPathList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetPathList(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Path) GetPathDetail(c *gin.Context) {
	data, err := h.svc.GetPathDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Path) PathSave(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	meta, msg, err := h.svc.PathSave(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, meta)
}

func (h *Path) PathDelete(c *gin.Context) {
	err := h.svc.PathDelete(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Path) GetUserPathList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetUserPathList(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Path) GetUserPathSearch(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	list, err := h.svc.GetUserPathSearch(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Path) UserPathListSave(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	err := h.svc.UserPathListSave(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "保存成功", []any{})
}

func (h *Path) UserPathDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	err := h.svc.UserPathDelete(c.Request.Context(), bind.RequestInt(c, raw, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Path) GetUserPathApplyList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetUserPathApplyList(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Path) UserPathApplyStatusChange(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	msg, err := h.svc.UserPathApplyStatusChange(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, []any{})
}

func (h *Path) UserPathApplyDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	err := h.svc.UserPathApplyDelete(c.Request.Context(), bind.RequestInt(c, raw, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", []any{})
}
