package pms

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/admin/service/staff"
)

type Staff struct{ svc *adminstaff.Service }

func NewStaff(svc *adminstaff.Service) *Staff {
	return &Staff{svc: svc}
}

func (h *Staff) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminstaff.ErrBadParams),
		errors.Is(err, adminstaff.ErrRoleNotFound),
		errors.Is(err, adminstaff.ErrUserExists),
		errors.Is(err, adminstaff.ErrUserNotFound),
		errors.Is(err, adminstaff.ErrNoCityAssign),
		errors.Is(err, adminstaff.ErrInvalidCity),
		errors.Is(err, adminstaff.ErrSuperProtected):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Staff) List(c *gin.Context) {
	var body struct {
		Username string   `json:"username"`
		Enable   bind.Int `json:"enable"`
		CityID   bind.Int `json:"city_id"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.List(c.Request.Context(), middleware.UID(c), adminstaff.ListQuery{
		Username: body.Username, Enable: body.Enable.Int(), CityID: body.CityID.Int(),
		Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "管理员列表", data)
}

func (h *Staff) AddRole(c *gin.Context) {
	h.changeRole(c, true)
}

func (h *Staff) RemoveRole(c *gin.Context) {
	h.changeRole(c, false)
}

func (h *Staff) changeRole(c *gin.Context, attach bool) {
	var body struct {
		UserIDs []bind.Int `json:"userIds"`
		RoleID  bind.Int   `json:"roleId"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.ChangeRole(c.Request.Context(), bind.Ints(body.UserIDs), body.RoleID.Int(), attach)) {
		return
	}
	if attach {
		response.OK(c, "添加成功", []any{})
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Staff) Add(c *gin.Context) {
	var body struct {
		Username string     `json:"username"`
		Password string     `json:"password"`
		Enable   bind.Int   `json:"enable"`
		RoleIDs  []bind.Int `json:"roleIds"`
		CityID   any        `json:"city_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	in := adminstaff.AddInput{
		Username: body.Username, Password: body.Password, Enable: body.Enable.Int(),
		RoleIDs: bind.Ints(body.RoleIDs),
	}
	if bind.Present(body.CityID) {
		in.HasCity = true
		cid := bind.AsInt(body.CityID)
		in.CityID = &cid
	}
	if h.writeErr(c, h.svc.Add(c.Request.Context(), middleware.UID(c), in)) {
		return
	}
	response.OK(c, "添加成功", []any{})
}

func (h *Staff) Edit(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	in := adminstaff.EditInput{ID: bind.AsInt(raw["id"])}
	if name, ok := raw["username"].(string); ok && name != "" {
		in.Username = &name
	}
	if _, ok := raw["enable"]; ok && raw["enable"] != nil && raw["enable"] != "" {
		n := bind.AsInt(raw["enable"])
		in.Enable = &n
	}
	if _, ok := raw["city_id"]; ok {
		in.CityProvided = true
		if raw["city_id"] == nil || raw["city_id"] == "" {
			in.ClearCity = true
		} else {
			in.CityID = bind.AsInt(raw["city_id"])
		}
	}
	if _, ok := raw["roleIds"]; ok {
		in.RoleIDsProvided = true
		in.RoleIDs = asInts(raw["roleIds"])
	}
	if h.writeErr(c, h.svc.Edit(c.Request.Context(), middleware.UID(c), in)) {
		return
	}
	response.OK(c, "编辑成功", []any{})
}

func (h *Staff) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), body.ID.Int())) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Staff) ResetPassword(c *gin.Context) {
	var body struct {
		ID       bind.Int `json:"id"`
		Password string   `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
		response.Fail(c, 1, "参数错误")
		return
	}
	if h.writeErr(c, h.svc.ResetPassword(c.Request.Context(), body.ID.Int(), body.Password)) {
		return
	}
	response.OK(c, "重置成功", []any{})
}

func asInts(value any) []int {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	out := make([]int, 0, len(items))
	for _, item := range items {
		out = append(out, bind.AsInt(item))
	}
	return out
}
