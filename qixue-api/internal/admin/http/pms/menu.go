package pms

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/menutree"
)

type Menu struct {
	db    *gorm.DB
	menus adminrepo.AdminMenuRepository
}

func NewMenu(db *gorm.DB) *Menu {
	return &Menu{db: db, menus: adminrepo.NewAdminMenu(db)}
}

func (h *Menu) Tree(c *gin.Context) {
	var rows []model.AdminMenu
	err := h.db.WithContext(c.Request.Context()).Where("type = ?", "MENU").Order("`order` ASC, id ASC").Find(&rows).Error
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", menutree.Tree(rows))
}

func (h *Menu) Buttons(c *gin.Context) {
	var rows []model.AdminMenu
	err := h.db.WithContext(c.Request.Context()).
		Where("parentId = ? AND type = ?", c.Query("parentId"), "BUTTON").
		Order("`order` ASC, id ASC").Find(&rows).Error
	if response.WriteErr(c, err) {
		return
	}
	list := make([]menutree.Node, 0, len(rows))
	for _, row := range rows {
		list = append(list, menutree.ToNode(row))
	}
	response.OK(c, "", list)
}

func (h *Menu) Add(c *gin.Context) {
	values, err := menuValues(c, false)
	if err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	now := time.Now()
	values["created_at"] = now
	values["updated_at"] = now
	if response.WriteErr(c, h.db.WithContext(c.Request.Context()).Model(&model.AdminMenu{}).Create(values).Error) {
		return
	}
	response.OK(c, "", []any{})
}

func (h *Menu) Save(c *gin.Context) {
	values, err := menuValues(c, true)
	if err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	id, _ := values["id"].(int)
	delete(values, "id")
	values["updated_at"] = time.Now()
	if id <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	if response.WriteErr(c, h.db.WithContext(c.Request.Context()).Model(&model.AdminMenu{}).Where("id = ?", id).Updates(values).Error) {
		return
	}
	response.OK(c, "", []any{})
}

func (h *Menu) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID.Int() <= 0 {
		response.Fail(c, 1, "菜单不存在")
		return
	}
	ctx := c.Request.Context()
	var count int64
	if err := h.db.WithContext(ctx).Model(&model.AdminMenu{}).Where("id = ?", body.ID.Int()).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count == 0 {
		response.Fail(c, 1, "菜单不存在")
		return
	}
	granted, err := h.menus.CountGranted(ctx, body.ID.Int())
	if response.WriteErr(c, err) {
		return
	}
	if granted > 0 {
		response.Fail(c, 1, "当前权限存在已授权的角色，不允许删除！")
		return
	}
	if response.WriteErr(c, h.db.WithContext(ctx).Delete(&model.AdminMenu{}, body.ID.Int()).Error) {
		return
	}
	response.OK(c, "", []any{})
}

func menuValues(c *gin.Context, keepID bool) (map[string]any, error) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		return nil, err
	}
	delete(raw, "children")
	allowed := map[string]bool{
		"name": true, "code": true, "type": true, "parentId": true, "path": true,
		"redirect": true, "icon": true, "component": true, "layout": true,
		"keepAlive": true, "method": true, "description": true, "show": true,
		"enable": true, "order": true,
	}
	if keepID {
		allowed["id"] = true
	}
	intKeys := map[string]bool{"id": true, "parentId": true, "keepAlive": true, "show": true, "enable": true, "order": true}
	out := map[string]any{}
	for key, value := range raw {
		if !allowed[key] {
			continue
		}
		if intKeys[key] {
			if value == nil {
				out[key] = nil
			} else {
				out[key] = bind.AsInt(value)
			}
			continue
		}
		out[key] = value
	}
	return out, nil
}
