package course

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	admincourse "qixue-api/internal/service/admin/course"
)

type Course struct{ svc *admincourse.Service }

func NewCourse(svc *admincourse.Service) *Course {
	return &Course{svc: svc}
}

func (h *Course) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var be *admincourse.BizError
	if errors.As(err, &be) {
		response.Fail(c, 1, be.Message)
		return true
	}
	return response.WriteErr(c, err)
}

func (h *Course) bindRaw(c *gin.Context) (map[string]any, bool) {
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

func (h *Course) GetCourseList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetCourseList(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Course) GetCourseDetail(c *gin.Context) {
	data, err := h.svc.GetCourseDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Course) CourseSave(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	meta, msg, err := h.svc.CourseSave(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, meta)
}

func (h *Course) CourseDelete(c *gin.Context) {
	err := h.svc.CourseDelete(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Course) GetCourseAttachList(c *gin.Context) {
	list, err := h.svc.GetCourseAttachList(c.Request.Context(), bind.Query(c, "courseId", "course_id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Course) CourseAttachSave(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	err := h.svc.CourseAttachSave(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "课件创建成功", []any{})
}

func (h *Course) CourseAttachDelete(c *gin.Context) {
	err := h.svc.CourseAttachDelete(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Course) GetCourseUserList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetCourseUserList(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Course) GetCourseUserSearch(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	list, err := h.svc.GetCourseUserSearch(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Course) CourseUserListSave(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	err := h.svc.CourseUserListSave(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "保存成功", []any{})
}

func (h *Course) CourseUserDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	err := h.svc.CourseUserDelete(c.Request.Context(), bind.RequestInt(c, raw, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Course) GetCourseUserApplyList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	data, err := h.svc.GetCourseUserApplyList(c.Request.Context(), middleware.UID(c), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Course) CourseUserApplyStatusChange(c *gin.Context) {
	raw, ok := h.bindRaw(c)
	if !ok {
		return
	}
	msg, err := h.svc.CourseUserApplyStatusChange(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, []any{})
}

func (h *Course) CourseUserApplyDelete(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	err := h.svc.CourseUserApplyDelete(c.Request.Context(), bind.RequestInt(c, raw, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", []any{})
}
