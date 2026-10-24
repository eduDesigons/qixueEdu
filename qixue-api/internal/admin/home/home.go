package home

import (
	"github.com/gin-gonic/gin"

	adminhome "qixue-api/internal/service/adminhome"
)

type Handler struct{ svc *adminhome.Service }

func New(svc *adminhome.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetVerse(c *gin.Context) { h.svc.GetVerse(c) }

func (h *Handler) GetCdnTraffic(c *gin.Context) { h.svc.GetCdnTraffic(c) }

func (h *Handler) GetPendingStats(c *gin.Context) { h.svc.GetPendingStats(c) }

func (h *Handler) GetTodayStats(c *gin.Context) { h.svc.GetTodayStats(c) }

func (h *Handler) GetStats(c *gin.Context) { h.svc.GetStats(c) }

func (h *Handler) GetChartData(c *gin.Context) { h.svc.GetChartData(c) }

func (h *Handler) GetCourseLearningRank(c *gin.Context) { h.svc.GetCourseLearningRank(c) }

func (h *Handler) GetChapterLearningRank(c *gin.Context) { h.svc.GetChapterLearningRank(c) }

func (h *Handler) GetChapterCompleteUserRank(c *gin.Context) { h.svc.GetChapterCompleteUserRank(c) }
