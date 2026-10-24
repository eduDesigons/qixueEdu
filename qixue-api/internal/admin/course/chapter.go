package course

import (
	"github.com/gin-gonic/gin"

	adminchapter "qixue-api/internal/service/adminchapter"
)

type Chapter struct{ svc *adminchapter.Service }

func NewChapter(svc *adminchapter.Service) *Chapter {
	return &Chapter{svc: svc}
}

func (h *Chapter) GetChapterList(c *gin.Context) { h.svc.GetChapterList(c) }

func (h *Chapter) GetChapterDetail(c *gin.Context) { h.svc.GetChapterDetail(c) }

func (h *Chapter) ChapterSave(c *gin.Context) { h.svc.ChapterSave(c) }

func (h *Chapter) ChapterDelete(c *gin.Context) { h.svc.ChapterDelete(c) }

func (h *Chapter) ChapterSort(c *gin.Context) { h.svc.ChapterSort(c) }

func (h *Chapter) GetChapterUserList(c *gin.Context) { h.svc.GetChapterUserList(c) }

func (h *Chapter) UnlockChapterUser(c *gin.Context) { h.svc.UnlockChapterUser(c) }
