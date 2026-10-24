package course

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/admin/comment"
)

type Comment struct{ svc *admincomment.Service }

func NewComment(svc *admincomment.Service) *Comment {
	return &Comment{svc: svc}
}

func (h *Comment) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, admincomment.ErrBadParams),
		errors.Is(err, admincomment.ErrNeedCommentID),
		errors.Is(err, admincomment.ErrEmptyReply),
		errors.Is(err, admincomment.ErrNotFound),
		errors.Is(err, admincomment.ErrNeedIDs):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Comment) GetCourseCommentList(c *gin.Context) {
	var body struct {
		CourseID    any      `json:"course_id"`
		CourseTitle string   `json:"course_title"`
		NickName    string   `json:"nick_name"`
		Content     string   `json:"content"`
		StartTime   string   `json:"start_time"`
		EndTime     string   `json:"end_time"`
		Page        bind.Int `json:"page"`
		Size        bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.ListCourse(c.Request.Context(), middleware.UID(c), admincomment.CourseListQuery{
		CourseID: bind.AsInt(body.CourseID), CourseTitle: body.CourseTitle, NickName: body.NickName,
		Content: body.Content, StartTime: body.StartTime, EndTime: body.EndTime,
		Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Comment) GetCourseChapterCommentList(c *gin.Context) {
	var body struct {
		CourseID     any      `json:"course_id"`
		ChapterID    any      `json:"chapter_id"`
		CourseTitle  string   `json:"course_title"`
		ChapterTitle string   `json:"chapter_title"`
		NickName     string   `json:"nick_name"`
		Content      string   `json:"content"`
		StartTime    string   `json:"start_time"`
		EndTime      string   `json:"end_time"`
		Page         bind.Int `json:"page"`
		Size         bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.ListChapter(c.Request.Context(), middleware.UID(c), admincomment.ChapterListQuery{
		CourseID: bind.AsInt(body.CourseID), ChapterID: bind.AsInt(body.ChapterID),
		CourseTitle: body.CourseTitle, ChapterTitle: body.ChapterTitle, NickName: body.NickName,
		Content: body.Content, StartTime: body.StartTime, EndTime: body.EndTime,
		Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", data)
}

func (h *Comment) AddCourseCommentReply(c *gin.Context) {
	var body struct {
		CommentID bind.Int `json:"comment_id"`
		Content   string   `json:"content"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.ReplyCourse(c.Request.Context(), middleware.UID(c), admincomment.ReplyInput{
		CommentID: body.CommentID.Int(), Content: body.Content,
	})) {
		return
	}
	response.OK(c, "回复成功", []any{})
}

func (h *Comment) AddChapterCommentReply(c *gin.Context) {
	var body struct {
		CommentID bind.Int `json:"comment_id"`
		Content   string   `json:"content"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.ReplyChapter(c.Request.Context(), middleware.UID(c), admincomment.ReplyInput{
		CommentID: body.CommentID.Int(), Content: body.Content,
	})) {
		return
	}
	response.OK(c, "回复成功", []any{})
}

func (h *Comment) GetCourseCommentReplies(c *gin.Context) {
	list, err := h.svc.CourseReplies(c.Request.Context(), bind.Query(c, "comment_id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Comment) GetChapterCommentReplies(c *gin.Context) {
	list, err := h.svc.ChapterReplies(c.Request.Context(), bind.Query(c, "comment_id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Comment) CourseCommentDelete(c *gin.Context) {
	if h.writeErr(c, h.svc.DeleteCourse(c.Request.Context(), bind.Query(c, "id"))) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Comment) CourseChapterCommentDelete(c *gin.Context) {
	if h.writeErr(c, h.svc.DeleteChapter(c.Request.Context(), bind.Query(c, "id"))) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Comment) CourseCommentBatchDelete(c *gin.Context) {
	var body struct {
		IDs []bind.Int `json:"ids"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.BatchDeleteCourse(c.Request.Context(), bind.Ints(body.IDs))) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Comment) CourseChapterCommentBatchDelete(c *gin.Context) {
	var body struct {
		IDs []bind.Int `json:"ids"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.BatchDeleteChapter(c.Request.Context(), bind.Ints(body.IDs))) {
		return
	}
	response.OK(c, "删除成功", []any{})
}
