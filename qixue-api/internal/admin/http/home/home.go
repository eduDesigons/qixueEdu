package home

import (
	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/admin/service/home"
)

type Handler struct{ svc *adminhome.Service }

func New(svc *adminhome.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetVerse(c *gin.Context) {
	data, err := h.svc.GetVerse(c.Request.Context())
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "随机一言", data)
}

func (h *Handler) GetCdnTraffic(c *gin.Context) {
	payload, err := h.svc.GetCdnTraffic(c.Request.Context())
	if err != nil {
		response.Fail(c, 1, "CDN 流量获取失败："+err.Error())
		return
	}
	response.OK(c, "CDN 流量", payload)
}

func (h *Handler) GetPendingStats(c *gin.Context) {
	data, err := h.svc.GetPendingStats(c.Request.Context(), middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "待办统计", data)
}

func (h *Handler) GetTodayStats(c *gin.Context) {
	data, err := h.svc.GetTodayStats(c.Request.Context(), middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "今日学情", data)
}

func (h *Handler) GetStats(c *gin.Context) {
	data, err := h.svc.GetStats(c.Request.Context(), middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "看板统计", data)
}

func (h *Handler) GetChartData(c *gin.Context) {
	data, err := h.svc.GetChartData(c.Request.Context(), middleware.UID(c))
	if err != nil {
		response.Fail(c, 1, "图表数据获取失败："+err.Error())
		return
	}
	response.OK(c, "图表数据", data)
}

func (h *Handler) GetCourseLearningRank(c *gin.Context) {
	response.OK(c, "课程学习中排行", h.svc.GetCourseLearningRank(c.Request.Context(), middleware.UID(c)))
}

func (h *Handler) GetChapterLearningRank(c *gin.Context) {
	response.OK(c, "近30天完成课节排行", h.svc.GetChapterLearningRank(c.Request.Context(), middleware.UID(c)))
}

func (h *Handler) GetChapterCompleteUserRank(c *gin.Context) {
	response.OK(c, "近两周积分排行", h.svc.GetChapterCompleteUserRank(c.Request.Context(), middleware.UID(c)))
}
