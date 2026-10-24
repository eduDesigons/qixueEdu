package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	usersvc "qixue-api/internal/service/webapp/user"
)

type Handler struct {
	svc *usersvc.Service
}

func New(svc *usersvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, usersvc.ErrNotFound),
		errors.Is(err, usersvc.ErrBadParam),
		errors.Is(err, usersvc.ErrPwdRequired),
		errors.Is(err, usersvc.ErrPwdLength),
		errors.Is(err, usersvc.ErrPwdMismatch):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) Info(c *gin.Context) {
	row, err := h.svc.Info(c.Request.Context(), middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "用户信息", row)
}

func (h *Handler) Edit(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	if h.writeSvcErr(c, h.svc.Edit(c.Request.Context(), middleware.UID(c), body)) {
		return
	}
	response.OK(c, "修改成功", nil)
}

func (h *Handler) SetPwd(c *gin.Context) {
	var body struct {
		OldPwd string `json:"old_pwd"`
		NewPwd string `json:"new_pwd"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeSvcErr(c, h.svc.SetPwd(c.Request.Context(), middleware.UID(c), body.OldPwd, body.NewPwd)) {
		return
	}
	response.OK(c, "修改成功", nil)
}

func (h *Handler) MyCount(c *gin.Context) {
	data, err := h.svc.MyCount(c.Request.Context(), middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "用户数据汇总", data)
}

func (h *Handler) MyMessages(c *gin.Context) {
	pageNum, limit := pageLimit(c)
	total, rows, err := h.svc.MyMessages(c.Request.Context(), middleware.UID(c), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, rows))
}

func (h *Handler) UnreadMessageCount(c *gin.Context) {
	count, err := h.svc.UnreadMessageCount(c.Request.Context(), middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "未读数量", gin.H{"count": count})
}

func (h *Handler) ReadMessage(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	id := body.ID.Int()
	if id == 0 {
		id = bind.Query(c, "id")
	}
	n, err := h.svc.ReadMessage(c.Request.Context(), middleware.UID(c), id)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "已读", gin.H{"count": n})
}

func (h *Handler) MyLikeCourse(c *gin.Context) {
	pageNum, limit := pageLimit(c)
	total, list, err := h.svc.MyLikeCourse(c.Request.Context(), middleware.UID(c), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, list))
}

func (h *Handler) MyLikeChapter(c *gin.Context) {
	pageNum, limit := pageLimit(c)
	total, list, err := h.svc.MyLikeChapter(c.Request.Context(), middleware.UID(c), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, list))
}

func (h *Handler) MyCourseComment(c *gin.Context) {
	h.pageMyComments(c, "course_comments", false)
}

func (h *Handler) MyChapterComment(c *gin.Context) {
	h.pageMyComments(c, "course_chapter_comments", true)
}

func (h *Handler) pageMyComments(c *gin.Context, table string, withChapter bool) {
	pageNum, limit := pageLimit(c)
	total, list, err := h.svc.MyComments(c.Request.Context(), table, middleware.UID(c), pageNum, limit, withChapter)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, list))
}

func (h *Handler) MyExam(c *gin.Context) {
	h.pageUserPapers(c, "course")
}

func (h *Handler) MyTest(c *gin.Context) {
	h.pageUserPapers(c, "chapter")
}

func (h *Handler) pageUserPapers(c *gin.Context, paperType string) {
	pageNum, limit := pageLimit(c)
	total, list, err := h.svc.MyPapers(c.Request.Context(), middleware.UID(c), paperType, pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, list))
}

func (h *Handler) MyCoinsLog(c *gin.Context) {
	h.pageRecords(c, "user_coins_records")
}

func (h *Handler) MyCreditLog(c *gin.Context) {
	h.pageRecords(c, "user_credit_records")
}

func (h *Handler) pageRecords(c *gin.Context, table string) {
	pageNum, limit := pageLimit(c)
	total, rows, err := h.svc.MyRecords(c.Request.Context(), table, middleware.UID(c), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, rows))
}

func (h *Handler) MyPath(c *gin.Context) {
	pageNum, limit := pageLimit(c)
	total, list, err := h.svc.MyPath(c.Request.Context(), middleware.UID(c), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, list))
}

func (h *Handler) MyCert(c *gin.Context) {
	pageNum := bind.Query(c, "page")
	limit := bind.Query(c, "limit")
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	total, list, err := h.svc.MyCert(c.Request.Context(), middleware.UID(c), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", page.NewMore(pageNum, limit, total, list))
}

func pageLimit(c *gin.Context) (int, int) {
	pageNum := bind.Query(c, "page")
	limit := bind.Query(c, "limit")
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 200 {
		limit = 200
	}
	return pageNum, limit
}
