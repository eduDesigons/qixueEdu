package adminrole

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/service/menutree"
)

var (
	ErrBadParams      = errors.New("参数错误")
	ErrDuplicate      = errors.New("角色已存在（角色名和角色编码不能重复）")
	ErrNotFound       = errors.New("角色不存在")
	ErrSuperProtected = errors.New("超级管理员角色不能删除")
	ErrHasUsers       = errors.New("当前角色存在已授权的用户，不允许删除！")
)

type Service struct {
	roles adminrepo.AdminRoleRepository
	db    *gorm.DB
}

func New(roles adminrepo.AdminRoleRepository, db *gorm.DB) *Service {
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

type View struct {
	ID            int     `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Enable        int     `json:"enable"`
	CreatedAt     *string `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
	PermissionIDs []int   `json:"permissionIds"`
}

type SaveInput struct {
	ID            int
	Code          string
	Name          string
	Enable        int
	PermissionIDs []int
}

type ListQuery struct {
	Name   string
	Enable int
	Page   int
	Size   int
}

func (s *Service) Add(ctx context.Context, in SaveInput) error {
	if in.Name == "" || in.Code == "" {
		return ErrBadParams
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&row{}).Where("name = ? AND code = ?", in.Name, in.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrDuplicate
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		item := row{Code: in.Code, Name: in.Name, Enable: in.Enable}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		return attachPermissions(tx, item.ID, in.PermissionIDs)
	})
}

func (s *Service) Edit(ctx context.Context, in SaveInput) error {
	if in.ID <= 0 {
		return ErrBadParams
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&row{}).Where("id = ?", in.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	if err := s.db.WithContext(ctx).Model(&row{}).Where("name = ? AND code = ? AND id <> ?", in.Name, in.Code, in.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrDuplicate
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&row{}).Where("id = ?", in.ID).Updates(map[string]any{
			"code": in.Code, "name": in.Name, "enable": in.Enable,
		}).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM administrator_permissions WHERE roleId = ?", in.ID).Error; err != nil {
			return err
		}
		return attachPermissions(tx, in.ID, in.PermissionIDs)
	})
}

func (s *Service) Delete(ctx context.Context, id int) error {
	var item row
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if item.Code == "SUPER_ADMIN" {
		return ErrSuperProtected
	}
	users, err := s.roles.CountUsersByRole(ctx, id)
	if err != nil {
		return err
	}
	if users > 0 {
		return ErrHasUsers
	}
	return s.db.WithContext(ctx).Delete(&row{}, id).Error
}

func (s *Service) List(ctx context.Context, q ListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	query := s.db.WithContext(ctx).Model(&row{})
	if q.Name != "" {
		query = query.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.Enable != 0 {
		query = query.Where("enable = ?", q.Enable)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []row
	if err := query.Order("created_at DESC, id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	ids := make([]int, 0, len(rows))
	for _, item := range rows {
		ids = append(ids, item.ID)
	}
	perms := map[int][]int{}
	if len(ids) > 0 {
		links, err := s.roles.ListPermissionLinks(ctx, ids)
		if err != nil {
			return page.Body{}, err
		}
		for _, link := range links {
			perms[link.RoleID] = append(perms[link.RoleID], link.PermissionID)
		}
	}
	list := make([]View, 0, len(rows))
	for _, item := range rows {
		pids := perms[item.ID]
		if pids == nil {
			pids = []int{}
		}
		list = append(list, View{
			ID: item.ID, Code: item.Code, Name: item.Name, Enable: item.Enable,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
			PermissionIDs: pids,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) MenuTree(ctx context.Context) ([]menutree.Node, error) {
	var rows []model.AdminMenu
	if err := s.db.WithContext(ctx).Order("`order` ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	kept := rows[:0]
	for _, row := range rows {
		if row.Code == "OpsInspect" {
			continue
		}
		kept = append(kept, row)
	}
	return menutree.Tree(kept), nil
}

func (s *Service) All(ctx context.Context) ([]View, error) {
	var rows []row
	if err := s.db.WithContext(ctx).Order("created_at DESC, id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]View, 0, len(rows))
	for _, item := range rows {
		list = append(list, View{
			ID: item.ID, Code: item.Code, Name: item.Name, Enable: item.Enable,
			CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
		})
	}
	return list, nil
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
