package member

import (
	"github.com/gin-gonic/gin"

	admincert "qixue-api/internal/service/admincert"
)

type Cert struct{ svc *admincert.Service }

func NewCert(svc *admincert.Service) *Cert {
	return &Cert{svc: svc}
}

func (h *Cert) List(c *gin.Context) { h.svc.List(c) }

func (h *Cert) Detail(c *gin.Context) { h.svc.Detail(c) }

func (h *Cert) Options(c *gin.Context) { h.svc.Options(c) }

func (h *Cert) SetShow(c *gin.Context) { h.svc.SetShow(c) }

func (h *Cert) Save(c *gin.Context) { h.svc.Save(c) }

func (h *Cert) Delete(c *gin.Context) { h.svc.Delete(c) }

func (h *Cert) Issue(c *gin.Context) { h.svc.Issue(c) }

func (h *Cert) Backfill(c *gin.Context) { h.svc.Backfill(c) }

func (h *Cert) IssuedList(c *gin.Context) { h.svc.IssuedList(c) }

func (h *Cert) Revoke(c *gin.Context) { h.svc.Revoke(c) }
