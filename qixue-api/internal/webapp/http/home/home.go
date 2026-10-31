package home

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	homesvc "qixue-api/internal/webapp/service/home"
)

type Handler struct {
	svc           *homesvc.Service
	continueStudy func(context.Context, int) (map[string]any, error)
}

func New(svc *homesvc.Service) *Handler {
	return &Handler{svc: svc}
}

// SetContinue 在公共首页数据之外，给当前登录用户附上最近一门在学课程。不写入公共缓存。
func (h *Handler) SetContinue(fn func(context.Context, int) (map[string]any, error)) {
	h.continueStudy = fn
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, homesvc.ErrBannerNotFound),
		errors.Is(err, homesvc.ErrNoticeNotFound):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) GetHomeInfo(c *gin.Context) {
	data, err := h.svc.GetHomeInfo(c.Request.Context())
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "homeInfo", data)
}

func (h *Handler) Boot(c *gin.Context) {
	payload, err := h.svc.Boot(c.Request.Context())
	if h.writeSvcErr(c, err) {
		return
	}
	if payload == nil {
		payload = map[string]any{}
	}
	if uid := middleware.UID(c); uid > 0 && h.continueStudy != nil {
		if row, err := h.continueStudy(c.Request.Context(), uid); err == nil && row != nil {
			payload["continueStudy"] = row
		}
	}
	response.OK(c, "homeBoot", payload)
}

func (h *Handler) GetBannerDetail(c *gin.Context) {
	data, err := h.svc.BannerDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "banner", data)
}

func (h *Handler) GetNoticeDetail(c *gin.Context) {
	data, err := h.svc.NoticeDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "notice", data)
}

func (h *Handler) GetHotClassListByCateId(c *gin.Context) {
	list, err := h.svc.HotByCate(c.Request.Context(), bind.Query(c, "id"), 10)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "首页热学榜课程", list)
}

func (h *Handler) GetHotPathList(c *gin.Context) {
	list, err := h.svc.HotPaths(c.Request.Context(), 10)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "首页热学榜路径", list)
}

func (h *Handler) GetHotCourseList(c *gin.Context) {
	limit := bind.Query(c, "limit")
	pageNum := bind.Query(c, "page")
	pageNum, limit = page.Clamp(pageNum, limit)
	list, _, total, err := h.svc.HotCoursesList(c.Request.Context(), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "首页热学榜课程", page.NewMore(pageNum, limit, total, list))
}
