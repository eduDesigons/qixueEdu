package course

import (
	"github.com/gin-gonic/gin"

	admincomment "qixue-api/internal/service/admincomment"
)

type Comment struct{ svc *admincomment.Service }

func NewComment(svc *admincomment.Service) *Comment {
	return &Comment{svc: svc}
}

func (h *Comment) GetCourseCommentList(c *gin.Context) { h.svc.GetCourseCommentList(c) }

func (h *Comment) GetCourseChapterCommentList(c *gin.Context) { h.svc.GetCourseChapterCommentList(c) }

func (h *Comment) AddCourseCommentReply(c *gin.Context) { h.svc.AddCourseCommentReply(c) }

func (h *Comment) AddChapterCommentReply(c *gin.Context) { h.svc.AddChapterCommentReply(c) }

func (h *Comment) GetCourseCommentReplies(c *gin.Context) { h.svc.GetCourseCommentReplies(c) }

func (h *Comment) GetChapterCommentReplies(c *gin.Context) { h.svc.GetChapterCommentReplies(c) }

func (h *Comment) CourseCommentDelete(c *gin.Context) { h.svc.CourseCommentDelete(c) }

func (h *Comment) CourseChapterCommentDelete(c *gin.Context) { h.svc.CourseChapterCommentDelete(c) }

func (h *Comment) CourseCommentBatchDelete(c *gin.Context) { h.svc.CourseCommentBatchDelete(c) }

func (h *Comment) CourseChapterCommentBatchDelete(c *gin.Context) {
	h.svc.CourseChapterCommentBatchDelete(c)
}
