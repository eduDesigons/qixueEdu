package common

import (
	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/admin/service/common"
)

type Handler struct{ svc *admincommon.Service }

func New(svc *admincommon.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Handler) Tree(c *gin.Context) {
	list, err := h.svc.Tree(c.Request.Context(), middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", list)
}
