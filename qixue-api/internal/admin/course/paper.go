package course

import (
	"github.com/gin-gonic/gin"

	adminpaper "qixue-api/internal/service/adminpaper"
)

type Paper struct{ svc *adminpaper.Service }

func NewPaper(svc *adminpaper.Service) *Paper {
	return &Paper{svc: svc}
}

func (h *Paper) GetUserPaperDetail(c *gin.Context) { h.svc.GetUserPaperDetail(c) }

func (h *Paper) PaperGrade(c *gin.Context) { h.svc.PaperGrade(c) }

func (h *Paper) GetPaperDetail(c *gin.Context) { h.svc.GetPaperDetail(c) }

func (h *Paper) PaperSave(c *gin.Context) { h.svc.PaperSave(c) }

func (h *Paper) PaperDelete(c *gin.Context) { h.svc.PaperDelete(c) }

func (h *Paper) GetExamList(c *gin.Context) { h.svc.GetExamList(c) }

func (h *Paper) GetAssessmentList(c *gin.Context) { h.svc.GetAssessmentList(c) }

func (h *Paper) GetPaperLogs(c *gin.Context) { h.svc.GetPaperLogs(c) }
