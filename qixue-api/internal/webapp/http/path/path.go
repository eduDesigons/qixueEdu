package path

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	pathsvc "qixue-api/internal/webapp/service/path"
)

type Handler struct {
	svc *pathsvc.Service
}

func New(svc *pathsvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) writeSvcErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, pathsvc.ErrNotFound),
		errors.Is(err, pathsvc.ErrAlreadyJoined),
		errors.Is(err, pathsvc.ErrCreditLow),
		errors.Is(err, pathsvc.ErrReasonEmpty),
		errors.Is(err, pathsvc.ErrApplyExists):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Handler) Categories(c *gin.Context) {
	list, err := h.svc.Categories(c.Request.Context())
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "学习路径分类", list)
}

func (h *Handler) GetPathList(c *gin.Context) {
	var q struct {
		CID   bind.Int `form:"cid" json:"cid"`
		Page  bind.Int `form:"page" json:"page"`
		Limit bind.Int `form:"limit" json:"limit"`
	}
	_ = c.ShouldBind(&q)
	pageNum, limit := page.Clamp(q.Page.Int(), q.Limit.Int())
	payload, err := h.svc.List(c.Request.Context(), q.CID.Int(), pageNum, limit)
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "ok", payload)
}

func (h *Handler) GetPathDetail(c *gin.Context) {
	id := bind.Query(c, "id")
	if id <= 0 {
		var body struct {
			ID bind.Int `json:"id"`
		}
		_ = c.ShouldBindJSON(&body)
		id = body.ID.Int()
	}
	data, err := h.svc.Detail(c.Request.Context(), id, middleware.UID(c))
	if h.writeSvcErr(c, err) {
		return
	}
	response.OK(c, "学习路径详情", data)
}

func (h *Handler) JoinPath(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id" form:"id"`
	}
	_ = c.ShouldBind(&body)
	id := body.ID.Int()
	if id <= 0 {
		id = bind.Query(c, "id")
	}
	if h.writeSvcErr(c, h.svc.Join(c.Request.Context(), id, middleware.UID(c))) {
		return
	}
	response.OK(c, "加入成功", nil)
}

func (h *Handler) ApplyPath(c *gin.Context) {
	var body struct {
		ID     bind.Int `json:"id" form:"id"`
		Reason string   `json:"reason" form:"reason"`
	}
	_ = c.ShouldBind(&body)
	if h.writeSvcErr(c, h.svc.Apply(c.Request.Context(), body.ID.Int(), middleware.UID(c), body.Reason)) {
		return
	}
	response.OK(c, "提交申请成功", nil)
}
