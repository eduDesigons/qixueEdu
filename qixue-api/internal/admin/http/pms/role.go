package pms

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/admin/service/role"
)

type Role struct{ svc *adminrole.Service }

func NewRole(svc *adminrole.Service) *Role {
	return &Role{svc: svc}
}

func (h *Role) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminrole.ErrBadParams),
		errors.Is(err, adminrole.ErrDuplicate),
		errors.Is(err, adminrole.ErrNotFound),
		errors.Is(err, adminrole.ErrSuperProtected),
		errors.Is(err, adminrole.ErrHasUsers):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Role) Add(c *gin.Context) {
	var body struct {
		Code          string     `json:"code"`
		Name          string     `json:"name"`
		Enable        bind.Int   `json:"enable"`
		PermissionIDs []bind.Int `json:"permissionIds"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" || body.Code == "" {
		response.Fail(c, 1, "参数错误")
		return
	}
	if h.writeErr(c, h.svc.Add(c.Request.Context(), adminrole.SaveInput{
		Code: body.Code, Name: body.Name, Enable: body.Enable.Int(),
		PermissionIDs: bind.Ints(body.PermissionIDs),
	})) {
		return
	}
	response.OK(c, "添加成功", []any{})
}

func (h *Role) Edit(c *gin.Context) {
	var body struct {
		ID            bind.Int   `json:"id"`
		Code          string     `json:"code"`
		Name          string     `json:"name"`
		Enable        bind.Int   `json:"enable"`
		PermissionIDs []bind.Int `json:"permissionIds"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ID.Int() <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	if h.writeErr(c, h.svc.Edit(c.Request.Context(), adminrole.SaveInput{
		ID: body.ID.Int(), Code: body.Code, Name: body.Name, Enable: body.Enable.Int(),
		PermissionIDs: bind.Ints(body.PermissionIDs),
	})) {
		return
	}
	response.OK(c, "编辑成功", []any{})
}

func (h *Role) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), body.ID.Int())) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *Role) List(c *gin.Context) {
	var body struct {
		Name   string   `json:"name"`
		Enable bind.Int `json:"enable"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.List(c.Request.Context(), adminrole.ListQuery{
		Name: body.Name, Enable: body.Enable.Int(), Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "角色列表", data)
}

func (h *Role) MenuTree(c *gin.Context) {
	tree, err := h.svc.MenuTree(c.Request.Context())
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", tree)
}

func (h *Role) All(c *gin.Context) {
	list, err := h.svc.All(c.Request.Context())
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "角色列表", list)
}
