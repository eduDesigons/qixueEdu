package decoration

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/adminnotice"
)

type Notice struct{ svc *adminnotice.Service }

func NewNotice(svc *adminnotice.Service) *Notice { return &Notice{svc: svc} }

func (h *Notice) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminnotice.ErrNotFound),
		errors.Is(err, adminnotice.ErrBadTime),
		errors.Is(err, adminnotice.ErrNeedID),
		errors.Is(err, adminnotice.ErrBadStatus):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Notice) List(c *gin.Context) {
	var body struct {
		Status  any      `json:"status"`
		Keyword string   `json:"keyword"`
		Page    bind.Int `json:"page"`
		Size    bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	var status *int
	if bind.Present(body.Status) {
		n := bind.AsInt(body.Status)
		status = &n
	}
	list, total, err := h.svc.List(c.Request.Context(), status, body.Keyword, current, size)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "通知列表", page.New(current, size, total, list))
}

func (h *Notice) Detail(c *gin.Context) {
	view, err := h.svc.Detail(c.Request.Context(), bind.Atoi(c.Query("id")))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "通知详情", view)
}

func (h *Notice) Save(c *gin.Context) {
	var body struct {
		ID          bind.Int `json:"id"`
		Title       string   `json:"title"`
		PublishTime string   `json:"publish_time"`
		Content     string   `json:"content"`
		Status      bind.Int `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	err := h.svc.Save(c.Request.Context(), adminnotice.SaveInput{
		ID: body.ID.Int(), Title: body.Title, PublishTime: body.PublishTime,
		Content: body.Content, Status: body.Status.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	if body.ID.Int() > 0 {
		response.OK(c, "更新成功", []any{})
		return
	}
	response.OK(c, "添加成功", []any{})
}

func (h *Notice) Status(c *gin.Context) {
	var body struct {
		ID     bind.Int `json:"id"`
		Status bind.Int `json:"status"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Status(c.Request.Context(), body.ID.Int(), body.Status.Int())) {
		return
	}
	response.OK(c, "状态更新成功", []any{})
}

func (h *Notice) Delete(c *gin.Context) {
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), bind.Atoi(c.Query("id")))) {
		return
	}
	response.OK(c, "删除成功", []any{})
}
