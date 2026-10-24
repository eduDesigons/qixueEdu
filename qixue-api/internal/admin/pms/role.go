package pms

import (
	"github.com/gin-gonic/gin"

	adminrole "qixue-api/internal/service/adminrole"
)

type Role struct{ svc *adminrole.Service }

func NewRole(svc *adminrole.Service) *Role {
	return &Role{svc: svc}
}

func (h *Role) Add(c *gin.Context) { h.svc.Add(c) }

func (h *Role) Edit(c *gin.Context) { h.svc.Edit(c) }

func (h *Role) Delete(c *gin.Context) { h.svc.Delete(c) }

func (h *Role) List(c *gin.Context) { h.svc.List(c) }

func (h *Role) MenuTree(c *gin.Context) { h.svc.MenuTree(c) }

func (h *Role) All(c *gin.Context) { h.svc.All(c) }
