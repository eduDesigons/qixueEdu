package comment

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	commentsvc "qixue-api/internal/webapp/service/comment"
)

type Handler struct {
	svc *commentsvc.Service
}

func New(svc *commentsvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, commentsvc.ErrEmptyContent),
		errors.Is(err, commentsvc.ErrNeedChapter),
		errors.Is(err, commentsvc.ErrNeedCourse),
		errors.Is(err, commentsvc.ErrNeedJoinCourse):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) GetChapterCommentList(c *gin.Context) {
	chapterID, pageNum, limit := listParams(c, true)
	data, err := h.svc.ListChapter(c.Request.Context(), chapterID, pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "章节讨论列表", data)
}

func (h *Handler) GetCourseCommentList(c *gin.Context) {
	courseID, pageNum, limit := listParams(c, false)
	data, err := h.svc.ListCourse(c.Request.Context(), courseID, pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "课程讨论列表", data)
}

func listParams(c *gin.Context, withChapter bool) (targetID, pageNum, limit int) {
	var body struct {
		ChapterID bind.Int `json:"chapter_id"`
		CourseID  bind.Int `json:"course_id"`
		Page      bind.Int `json:"page"`
		Limit     bind.Int `json:"limit"`
	}
	_ = c.ShouldBindJSON(&body)
	if withChapter {
		targetID = body.ChapterID.Int()
		if targetID == 0 {
			targetID = bind.Query(c, "chapter_id")
		}
	} else {
		targetID = body.CourseID.Int()
		if targetID == 0 {
			targetID = bind.Query(c, "course_id")
		}
	}
	pageNum = body.Page.Int()
	limit = body.Limit.Int()
	if pageNum == 0 {
		pageNum = bind.Query(c, "page")
	}
	if limit == 0 {
		limit = bind.Query(c, "limit")
	}
	return targetID, pageNum, limit
}

func (h *Handler) ChapterComment(c *gin.Context) {
	uid := middleware.UID(c)
	var body struct {
		ChapterID bind.Int `json:"chapter_id"`
		CourseID  bind.Int `json:"course_id"`
		Content   string   `json:"content"`
		PID       bind.Int `json:"pid"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.ChapterComment(c.Request.Context(), uid,
		body.ChapterID.Int(), body.CourseID.Int(), body.PID.Int(), body.Content)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "评论成功", data)
}

func (h *Handler) CourseComment(c *gin.Context) {
	uid := middleware.UID(c)
	var body struct {
		CourseID bind.Int `json:"course_id"`
		Content  string   `json:"content"`
		PID      bind.Int `json:"pid"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.CourseComment(c.Request.Context(), uid,
		body.CourseID.Int(), body.PID.Int(), body.Content)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "评论成功", data)
}
