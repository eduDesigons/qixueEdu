package course

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	coursesvc "qixue-api/internal/service/course"
)

type Handler struct {
	svc *coursesvc.Service
}

func New(svc *coursesvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, coursesvc.ErrNotFound),
		errors.Is(err, coursesvc.ErrChapterNotFound):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) Categories(c *gin.Context) {
	list, err := h.svc.Categories(c.Request.Context())
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "课程分类", list)
}

type listFilter struct {
	CateID1  bind.Int `json:"cateId1" form:"cateId1"`
	CateID2  bind.Int `json:"cateId2" form:"cateId2"`
	Type     bind.Int `json:"type" form:"type"`
	SortType string   `json:"sortType" form:"sortType"`
	Title    string   `json:"title" form:"title"`
	Page     bind.Int `json:"page" form:"page"`
	Limit    bind.Int `json:"limit" form:"limit"`
}

func (h *Handler) GetCourseList(c *gin.Context) {
	var filter listFilter
	_ = c.ShouldBind(&filter)
	pageNum, limit := page.Clamp(filter.Page.Int(), filter.Limit.Int())
	payload, err := h.svc.List(c.Request.Context(), coursesvc.ListFilter{
		CateID1: filter.CateID1.Int(), CateID2: filter.CateID2.Int(),
		Type: filter.Type.Int(), SortType: filter.SortType, Title: filter.Title,
		Page: pageNum, Limit: limit,
	})
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", payload)
}

func (h *Handler) CourseDetail(c *gin.Context) {
	id := bind.Query(c, "id")
	refresh := c.Query("refresh") == "1"
	if id <= 0 {
		var body struct {
			ID      bind.Int `json:"id"`
			Refresh bind.Int `json:"refresh"`
		}
		_ = c.ShouldBindJSON(&body)
		id = body.ID.Int()
		if body.Refresh.Int() == 1 {
			refresh = true
		}
	}
	data, err := h.svc.Detail(c.Request.Context(), id, middleware.UID(c), refresh)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "课程详情", data)
}

func (h *Handler) GetCourseChapterList(c *gin.Context) {
	id := bind.Query(c, "course_id")
	if id <= 0 {
		var body struct {
			CourseID bind.Int `json:"course_id"`
		}
		_ = c.ShouldBindJSON(&body)
		id = body.CourseID.Int()
	}
	list, err := h.svc.ChapterList(c.Request.Context(), id, middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "课程章节列表", list)
}

func (h *Handler) ChapterDetail(c *gin.Context) {
	id := bind.Query(c, "id")
	refresh := c.Query("refresh") == "1"
	if id <= 0 {
		var body struct {
			ID      bind.Int `json:"id"`
			Refresh bind.Int `json:"refresh"`
		}
		_ = c.ShouldBindJSON(&body)
		id = body.ID.Int()
		if body.Refresh.Int() == 1 {
			refresh = true
		}
	}
	data, err := h.svc.ChapterDetail(c.Request.Context(), id, middleware.UID(c), refresh)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "课程章节详情", data)
}
