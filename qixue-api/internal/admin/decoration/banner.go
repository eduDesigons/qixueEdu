package decoration

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/adminbanner"
)

type Banner struct{ svc *adminbanner.Service }

func NewBanner(svc *adminbanner.Service) *Banner { return &Banner{svc: svc} }

func (h *Banner) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminbanner.ErrNotFound),
		errors.Is(err, adminbanner.ErrNeedIDs):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Banner) List(c *gin.Context) {
	var typ, status *int
	if v := strings.TrimSpace(c.Query("type")); v != "" {
		n := bind.Atoi(v)
		typ = &n
	}
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		n := bind.Atoi(v)
		status = &n
	}
	list, err := h.svc.List(c.Request.Context(), typ, status)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "Banner列表", list)
}

func (h *Banner) Save(c *gin.Context) {
	var body struct {
		ID      bind.Int `json:"id"`
		Name    string   `json:"name"`
		ImgURL  string   `json:"img_url"`
		Sort    bind.Int `json:"sort"`
		LinkURL *string  `json:"link_url"`
		Type    bind.Int `json:"type"`
		Status  bind.Int `json:"status"`
		Content *string  `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	err := h.svc.Save(c.Request.Context(), adminbanner.SaveInput{
		ID: body.ID.Int(), Name: body.Name, ImgURL: body.ImgURL, Sort: body.Sort.Int(),
		LinkURL: body.LinkURL, Type: body.Type.Int(), Status: body.Status.Int(), Content: body.Content,
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

func (h *Banner) Delete(c *gin.Context) {
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), bind.Query(c, "id"))) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Banner) Sort(c *gin.Context) {
	var body struct {
		IDs []bind.Int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "请传入排序后的 id 列表")
		return
	}
	if h.writeErr(c, h.svc.Sort(c.Request.Context(), bind.Ints(body.IDs))) {
		return
	}
	response.OK(c, "排序已更新", []any{})
}
