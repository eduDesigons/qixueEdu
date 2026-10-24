package ops

import (
	"github.com/gin-gonic/gin"

	adminops "qixue-api/internal/service/adminops"
)

type Handler struct{ svc *adminops.Service }

func New(svc *adminops.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Chart(c *gin.Context) { h.svc.Chart(c) }

func (h *Handler) CourseEnrollList(c *gin.Context) { h.svc.CourseEnrollList(c) }

func (h *Handler) PathEnrollList(c *gin.Context) { h.svc.PathEnrollList(c) }

func (h *Handler) LearnList(c *gin.Context) { h.svc.LearnList(c) }

func (h *Handler) CreditList(c *gin.Context) { h.svc.CreditList(c) }

func (h *Handler) CoinsList(c *gin.Context) { h.svc.CoinsList(c) }
