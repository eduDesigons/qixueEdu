package pms

import (
	"github.com/gin-gonic/gin"

	adminstaff "qixue-api/internal/service/adminstaff"
)

type Staff struct{ svc *adminstaff.Service }

func NewStaff(svc *adminstaff.Service) *Staff {
	return &Staff{svc: svc}
}

func (h *Staff) List(c *gin.Context) { h.svc.List(c) }

func (h *Staff) AddRole(c *gin.Context) { h.svc.AddRole(c) }

func (h *Staff) RemoveRole(c *gin.Context) { h.svc.RemoveRole(c) }

func (h *Staff) Add(c *gin.Context) { h.svc.Add(c) }

func (h *Staff) Edit(c *gin.Context) { h.svc.Edit(c) }

func (h *Staff) Delete(c *gin.Context) { h.svc.Delete(c) }

func (h *Staff) ResetPassword(c *gin.Context) { h.svc.ResetPassword(c) }
