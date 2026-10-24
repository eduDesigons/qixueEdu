package media

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/admin/mediacate"
)

type Cate struct{ svc *adminmediacate.Service }

func NewCate(db *gorm.DB) *Cate {
	return &Cate{svc: adminmediacate.New(db)}
}

func (h *Cate) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminmediacate.ErrBadParams),
		errors.Is(err, adminmediacate.ErrNotFound):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Cate) List(c *gin.Context) {
	var body struct {
		Type string `json:"type"`
	}
	_ = c.ShouldBindJSON(&body)
	list, err := h.svc.List(c.Request.Context(), body.Type)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "分类列表", list)
}

func (h *Cate) Destroy(c *gin.Context) {
	var body struct {
		ID any `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Destroy(c.Request.Context(), collectIDs(body.ID))) {
		return
	}
	response.OK(c, "删除成功", nil)
}

func collectIDs(v any) []int {
	seen := map[int]struct{}{}
	var out []int
	add := func(n int) {
		if n <= 0 {
			return
		}
		if _, ok := seen[n]; ok {
			return
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	switch x := v.(type) {
	case []any:
		for _, item := range x {
			for _, id := range collectIDs(item) {
				add(id)
			}
		}
	case []int:
		for _, id := range x {
			add(id)
		}
	default:
		add(bind.AsInt(v))
	}
	return out
}
