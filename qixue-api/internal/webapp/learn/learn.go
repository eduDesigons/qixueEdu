package learn

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	learnsvc "qixue-api/internal/service/webapp/learn"
)

type Handler struct {
	svc *learnsvc.Service
}

func New(svc *learnsvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	code := 1
	switch {
	case errors.Is(err, learnsvc.ErrNeedUnlock):
		code = 400
		response.Fail(c, code, err.Error())
		return true
	case errors.Is(err, learnsvc.ErrAlreadyJoined),
		errors.Is(err, learnsvc.ErrCourseNotFound),
		errors.Is(err, learnsvc.ErrNoChapters),
		errors.Is(err, learnsvc.ErrCreditShort),
		errors.Is(err, learnsvc.ErrJoinFailed),
		errors.Is(err, learnsvc.ErrReasonRequired),
		errors.Is(err, learnsvc.ErrApplyNoChapters),
		errors.Is(err, learnsvc.ErrAlreadyApplied),
		errors.Is(err, learnsvc.ErrChapterNotFound),
		errors.Is(err, learnsvc.ErrNeedJoinCourse),
		errors.Is(err, learnsvc.ErrAlreadyEnded),
		errors.Is(err, learnsvc.ErrAlreadyUnlocked),
		errors.Is(err, learnsvc.ErrNeedPrevChapter),
		errors.Is(err, learnsvc.ErrCourseAlreadyEnded),
		errors.Is(err, learnsvc.ErrCourseNeedExam):
		response.Fail(c, code, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) MyCourse(c *gin.Context) {
	data, err := h.svc.MyCourse(c.Request.Context(), middleware.UID(c),
		bind.Query(c, "type"), bind.Query(c, "page"), bind.Query(c, "limit"))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Handler) MyTask(c *gin.Context) {
	list, err := h.svc.MyTask(c.Request.Context(), middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "待完成任务", list)
}

func (h *Handler) MyHistory(c *gin.Context) {
	list, total, hasMore, err := h.svc.MyHistory(c.Request.Context(), middleware.UID(c),
		bind.Query(c, "page"), bind.Query(c, "limit"))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "学习轨迹", gin.H{"list": list, "total": total, "hasMore": hasMore})
}

func (h *Handler) SoonExpire(c *gin.Context) {
	pageNum := bind.Query(c, "page")
	limit := bind.Query(c, "limit")
	pageNum, limit = page.Clamp(pageNum, limit)
	if limit > 50 {
		limit = 50
	}
	list, total, err := h.svc.SoonExpire(c.Request.Context(), middleware.UID(c), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, list))
}

func (h *Handler) LearnNum(c *gin.Context) {
	data, err := h.svc.LearnNum(c.Request.Context(), middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "学习统计", data)
}

func (h *Handler) Collect(c *gin.Context) {
	var body struct {
		CourseID bind.Int `json:"course_id" form:"course_id"`
	}
	_ = c.ShouldBind(&body)
	courseID := body.CourseID.Int()
	if courseID <= 0 {
		courseID = bind.Query(c, "course_id")
	}
	if h.writeSvcErr(c, h.svc.Collect(c.Request.Context(), middleware.UID(c), courseID)) {
		return
	}
	response.OK(c, "操作成功", nil)
}

func (h *Handler) CollectChapter(c *gin.Context) {
	var body struct {
		ChapterID bind.Int `json:"chapter_id" form:"chapter_id"`
	}
	_ = c.ShouldBind(&body)
	chapterID := body.ChapterID.Int()
	if chapterID <= 0 {
		chapterID = bind.Query(c, "chapter_id")
	}
	if h.writeSvcErr(c, h.svc.CollectChapter(c.Request.Context(), middleware.UID(c), chapterID)) {
		return
	}
	response.OK(c, "操作成功", nil)
}

func (h *Handler) JoinCourse(c *gin.Context) {
	var body struct {
		CourseID bind.Int `json:"course_id" form:"course_id"`
	}
	_ = c.ShouldBind(&body)
	courseID := body.CourseID.Int()
	if courseID <= 0 {
		courseID = bind.Query(c, "course_id")
	}
	if h.writeSvcErr(c, h.svc.JoinCourse(c.Request.Context(), middleware.UID(c), courseID)) {
		return
	}
	response.OK(c, "加入成功", nil)
}

func (h *Handler) ApplyCourse(c *gin.Context) {
	var body struct {
		ID     bind.Int `json:"id" form:"id"`
		Reason string   `json:"reason" form:"reason"`
	}
	_ = c.ShouldBind(&body)
	id := body.ID.Int()
	if id <= 0 {
		id = bind.Query(c, "id")
	}
	if h.writeSvcErr(c, h.svc.ApplyCourse(c.Request.Context(), middleware.UID(c), id, body.Reason)) {
		return
	}
	response.OK(c, "提交申请成功", nil)
}

func (h *Handler) UnlockChapter(c *gin.Context) {
	var body struct {
		ChapterID bind.Int `json:"chapter_id" form:"chapter_id"`
	}
	_ = c.ShouldBind(&body)
	chapterID := body.ChapterID.Int()
	if chapterID <= 0 {
		chapterID = bind.Query(c, "chapter_id")
	}
	if chapterID <= 0 {
		response.Fail(c, 1, "缺少章节参数")
		return
	}
	if h.writeSvcErr(c, h.svc.UnlockChapter(c.Request.Context(), middleware.UID(c), chapterID)) {
		return
	}
	response.OK(c, "解锁成功", nil)
}

func (h *Handler) LearnEnd(c *gin.Context) {
	var body struct {
		ChapterID bind.Int `json:"chapter_id" form:"chapter_id"`
	}
	_ = c.ShouldBind(&body)
	chapterID := body.ChapterID.Int()
	if chapterID <= 0 {
		chapterID = bind.Query(c, "chapter_id")
	}
	row, err := h.svc.LearnEnd(c.Request.Context(), middleware.UID(c), chapterID)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "学习完成上报成功", row)
}

func (h *Handler) EndCourse(c *gin.Context) {
	var body struct {
		CourseID bind.Int `json:"course_id" form:"course_id"`
	}
	_ = c.ShouldBind(&body)
	courseID := body.CourseID.Int()
	if courseID <= 0 {
		courseID = bind.Query(c, "course_id")
	}
	if h.writeSvcErr(c, h.svc.EndCourse(c.Request.Context(), middleware.UID(c), courseID)) {
		return
	}
	response.OK(c, "课程结束成功", nil)
}
