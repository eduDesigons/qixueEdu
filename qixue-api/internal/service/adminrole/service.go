package adminrole

import (
	"time"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/menutree"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

type Service struct {
	roles repo.AdminRoleRepository
	db    *gorm.DB
}

func New(roles repo.AdminRoleRepository, db *gorm.DB) *Service {
	return &Service{roles: roles, db: db}
}

type row struct {
	ID        int        `gorm:"column:id;primaryKey" json:"id"`
	Code      string     `gorm:"column:code" json:"code"`
	Name      string     `gorm:"column:name" json:"name"`
	Enable    int        `gorm:"column:enable" json:"enable"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"-"`
}

func (row) TableName() string { return "administrator_roles" }

type view struct {
	ID            int     `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Enable        int     `json:"enable"`
	CreatedAt     *string `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
	PermissionIDs []int   `json:"permissionIds"`
}

func (s *Service) Add(c *gin.Context) {
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
	ctx := c.Request.Context()
	var count int64
	if err := s.db.WithContext(ctx).Model(&row{}).Where("name = ? AND code = ?", body.Name, body.Code).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count > 0 {
		response.Fail(c, 1, "角色已存在（角色名和角色编码不能重复）")
		return
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		item := row{Code: body.Code, Name: body.Name, Enable: body.Enable.Int()}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		return attachPermissions(tx, item.ID, bind.Ints(body.PermissionIDs))
	})
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "添加成功", []any{})
}

func (s *Service) Edit(c *gin.Context) {
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
	ctx := c.Request.Context()
	var count int64
	if err := s.db.WithContext(ctx).Model(&row{}).Where("id = ?", body.ID.Int()).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count == 0 {
		response.Fail(c, 1, "角色不存在")
		return
	}
	if err := s.db.WithContext(ctx).Model(&row{}).Where("name = ? AND code = ? AND id <> ?", body.Name, body.Code, body.ID.Int()).Count(&count).Error; err != nil {
		response.WriteErr(c, err)
		return
	}
	if count > 0 {
		response.Fail(c, 1, "角色已存在（角色名和角色编码不能重复）")
		return
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&row{}).Where("id = ?", body.ID.Int()).Updates(map[string]any{
			"code": body.Code, "name": body.Name, "enable": body.Enable.Int(),
		}).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM administrator_permissions WHERE roleId = ?", body.ID.Int()).Error; err != nil {
			return err
		}
		return attachPermissions(tx, body.ID.Int(), bind.Ints(body.PermissionIDs))
	})
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "编辑成功", []any{})
}

func (s *Service) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	var item row
	err := s.db.WithContext(c.Request.Context()).Where("id = ?", body.ID.Int()).Take(&item).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "角色不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if item.Code == "SUPER_ADMIN" {
		response.Fail(c, 1, "超级管理员角色不能删除")
		return
	}
	users, err := s.roles.CountUsersByRole(c.Request.Context(), body.ID.Int())
	if response.WriteErr(c, err) {
		return
	}
	if users > 0 {
		response.Fail(c, 1, "当前角色存在已授权的用户，不允许删除！")
		return
	}
	if response.WriteErr(c, s.db.WithContext(c.Request.Context()).Delete(&row{}, body.ID.Int()).Error) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (s *Service) List(c *gin.Context) {
	var body struct {
		Name   string   `json:"name"`
		Enable bind.Int `json:"enable"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	query := s.db.WithContext(c.Request.Context()).Model(&row{})
	if body.Name != "" {
		query = query.Where("name LIKE ?", "%"+body.Name+"%")
	}
	if body.Enable.Int() != 0 {
		query = query.Where("enable = ?", body.Enable.Int())
	}
	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []row
	err := query.Order("created_at DESC, id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error
	if response.WriteErr(c, err) {
		return
	}
	ids := make([]int, 0, len(rows))
	for _, item := range rows {
		ids = append(ids, item.ID)
	}
	perms := map[int][]int{}
	if len(ids) > 0 {
		links, err := s.roles.ListPermissionLinks(c.Request.Context(), ids)
		if response.WriteErr(c, err) {
			return
		}
		for _, link := range links {
			perms[link.RoleID] = append(perms[link.RoleID], link.PermissionID)
		}
	}
	list := make([]view, 0, len(rows))
	for _, item := range rows {
		ids := perms[item.ID]
		if ids == nil {
			ids = []int{}
		}
		list = append(list, view{
			ID: item.ID, Code: item.Code, Name: item.Name, Enable: item.Enable,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
			PermissionIDs: ids,
		})
	}
	response.OK(c, "角色列表", page.New(current, size, total, list))
}

func (s *Service) MenuTree(c *gin.Context) {
	var rows []model.AdminMenu
	err := s.db.WithContext(c.Request.Context()).Order("`order` ASC, id ASC").Find(&rows).Error
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", menutree.Tree(rows))
}

func (s *Service) All(c *gin.Context) {
	var rows []row
	err := s.db.WithContext(c.Request.Context()).Order("created_at DESC, id DESC").Find(&rows).Error
	if response.WriteErr(c, err) {
		return
	}
	list := make([]view, 0, len(rows))
	for _, item := range rows {
		list = append(list, view{
			ID: item.ID, Code: item.Code, Name: item.Name, Enable: item.Enable,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
		})
	}
	response.OK(c, "角色列表", list)
}

func attachPermissions(tx *gorm.DB, roleID int, ids []int) error {
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if err := tx.Exec("INSERT INTO administrator_permissions (roleId, permissionId) VALUES (?, ?)", roleID, id).Error; err != nil {
			return err
		}
	}
	return nil
}
