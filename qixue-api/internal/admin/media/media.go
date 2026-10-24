package media

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	adminmedia "qixue-api/internal/service/adminmedia"
)

type Media struct{ svc *adminmedia.Service }

func NewMedia(db *gorm.DB) *Media {
	return &Media{svc: adminmedia.New(db)}
}

func (h *Media) List(c *gin.Context) { h.svc.List(c) }

func (h *Media) Upload(c *gin.Context) { h.svc.Upload(c) }

func (h *Media) Destroy(c *gin.Context) { h.svc.Destroy(c) }

func (h *Media) SaveUploaded(c *gin.Context) { h.svc.SaveUploaded(c) }

func (h *Media) InitVideoDirectUpload(c *gin.Context) { h.svc.InitVideoDirectUpload(c) }

func (h *Media) CompleteVideoDirectUpload(c *gin.Context) { h.svc.CompleteVideoDirectUpload(c) }

func (h *Media) AbortVideoDirectUpload(c *gin.Context) { h.svc.AbortVideoDirectUpload(c) }
