package rank

import (
	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	ranksvc "qixue-api/internal/service/rank"
)

type Handler struct {
	svc *ranksvc.Service
}

func New(svc *ranksvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Boards(c *gin.Context) {
	uid := middleware.UID(c)
	data, err := h.svc.Boards(c.Request.Context(), uid)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "学习榜单", data)
}

func (h *Handler) History(c *gin.Context) {
	uid := middleware.UID(c)
	typ := c.Query("type")
	if typ == "" {
		typ = "duration"
	}
	period := c.Query("period")
	if period == "" {
		period = "week"
	}
	list, err := h.svc.History(c.Request.Context(), typ, period, uid)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "历史榜单", gin.H{"list": list})
}
