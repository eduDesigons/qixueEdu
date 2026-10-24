package media

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	adminmediacate "qixue-api/internal/service/adminmediacate"
)

type Cate struct{ svc *adminmediacate.Service }

func NewCate(db *gorm.DB) *Cate {
	return &Cate{svc: adminmediacate.New(db)}
}

func (h *Cate) List(c *gin.Context) { h.svc.List(c) }

func (h *Cate) Destroy(c *gin.Context) { h.svc.Destroy(c) }
