package course

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/admin/category"
)

type Category struct{ svc *admincategory.Service }

func NewCategory(svc *admincategory.Service) *Category {
	return &Category{svc: svc}
}

func (h *Category) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, admincategory.ErrBadParams),
		errors.Is(err, admincategory.ErrEmptyName),
		errors.Is(err, admincategory.ErrNameTooLong),
		errors.Is(err, admincategory.ErrNotFound),
		errors.Is(err, admincategory.ErrHasPaths),
		errors.Is(err, admincategory.ErrHasChildren),
		errors.Is(err, admincategory.ErrHasCourses),
		errors.Is(err, admincategory.ErrIconTooLong),
		errors.Is(err, admincategory.ErrParentMissing),
		errors.Is(err, admincategory.ErrParentHasCourses),
		errors.Is(err, admincategory.ErrDupName),
		errors.Is(err, admincategory.ErrSelfParent),
		errors.Is(err, admincategory.ErrHasKidsCantDemote),
		errors.Is(err, admincategory.ErrHasCoursesCantPromote),
		errors.Is(err, admincategory.ErrPathDupName),
		errors.Is(err, admincategory.ErrMinCreditNegative),
		errors.Is(err, admincategory.ErrBaseCreditNegative):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Category) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), c.Query("type"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (h *Category) Save(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	in := admincategory.SaveInput{
		Type:     asString(raw["type"]),
		ID:       bind.AsInt(raw["id"]),
		Name:     asString(raw["name"]),
		Sort:     bind.AsInt(raw["sort"]),
		ParentID: bind.AsInt(raw["parent_id"]),
		IsShow:   1,
		Icon:     asString(raw["icon"]),
	}
	if _, ok := raw["is_show"]; ok && bind.AsInt(raw["is_show"]) != 1 {
		in.IsShow = 0
	}
	if _, ok := raw["min_need_credit"]; ok {
		n := bind.AsInt(raw["min_need_credit"])
		in.MinNeedCredit = &n
	}
	if _, ok := raw["base_credit"]; ok {
		n := bind.AsInt(raw["base_credit"])
		in.BaseCredit = &n
	}
	result, err := h.svc.Save(c.Request.Context(), in)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "保存成功", gin.H{"id": result.ID})
}

func (h *Category) Delete(c *gin.Context) {
	var body struct {
		Type string   `json:"type"`
		ID   bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), body.Type, body.ID.Int())) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Category) Options(c *gin.Context) {
	list, err := h.svc.Options(c.Request.Context())
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func asString(value any) string {
	text, _ := value.(string)
	return strings.TrimSpace(text)
}
