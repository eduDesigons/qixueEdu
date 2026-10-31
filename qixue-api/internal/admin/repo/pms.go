package admin

import (
	"context"

	"gorm.io/gorm"
)

type VerseRow struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`
}

func (VerseRow) TableName() string { return "verse" }

type AdminVerseRepository interface {
	List(ctx context.Context, keyword string, offset, limit int) ([]VerseRow, int64, error)
	UpdateContent(ctx context.Context, id int, content string) (int64, error)
	Count(ctx context.Context, id int) (int64, error)
	Create(ctx context.Context, row *VerseRow) error
	Delete(ctx context.Context, id int) (int64, error)
}

type adminVerseRepo struct{ db *gorm.DB }

func NewAdminVerse(db *gorm.DB) AdminVerseRepository { return &adminVerseRepo{db: db} }

func (r *adminVerseRepo) List(ctx context.Context, keyword string, offset, limit int) ([]VerseRow, int64, error) {
	q := r.db.WithContext(ctx).Model(&VerseRow{})
	if keyword != "" {
		q = q.Where("content LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []VerseRow
	err := q.Order("id DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *adminVerseRepo) UpdateContent(ctx context.Context, id int, content string) (int64, error) {
	res := r.db.WithContext(ctx).Model(&VerseRow{}).Where("id = ?", id).Update("content", content)
	return res.RowsAffected, res.Error
}

func (r *adminVerseRepo) Count(ctx context.Context, id int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&VerseRow{}).Where("id = ?", id).Count(&n).Error
	return n, err
}

func (r *adminVerseRepo) Create(ctx context.Context, row *VerseRow) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *adminVerseRepo) Delete(ctx context.Context, id int) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&VerseRow{}, id)
	return res.RowsAffected, res.Error
}

type AdminAreaCommonRepository interface {
	ListVisible(ctx context.Context, cityID int, limited bool) ([]AreaItem, error)
	ListParents(ctx context.Context, cityID int, limited bool) ([]AreaItem, error)
	ListChildren(ctx context.Context, cityID int, limited bool) ([]AreaItem, error)
}

type AreaItem struct {
	ID       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	ParentID int    `json:"parent_id" gorm:"column:parent_id"`
}

type adminAreaCommonRepo struct{ db *gorm.DB }

func NewAdminAreaCommon(db *gorm.DB) AdminAreaCommonRepository {
	return &adminAreaCommonRepo{db: db}
}

func (r *adminAreaCommonRepo) ListVisible(ctx context.Context, cityID int, limited bool) ([]AreaItem, error) {
	q := r.db.WithContext(ctx).Table("areas").Where("is_show = 1 AND deleted_at IS NULL").Order("sort ASC, id ASC")
	if limited {
		q = q.Where("id = ? OR parent_id = ?", cityID, cityID)
	}
	var list []AreaItem
	err := q.Find(&list).Error
	return list, err
}

func (r *adminAreaCommonRepo) ListParents(ctx context.Context, cityID int, limited bool) ([]AreaItem, error) {
	q := r.db.WithContext(ctx).Table("areas").
		Select("id, name, parent_id").
		Where("parent_id = 0 AND is_show = 1 AND deleted_at IS NULL").
		Order("sort ASC, id ASC")
	if limited {
		q = q.Where("id = ?", cityID)
	}
	var parents []AreaItem
	err := q.Find(&parents).Error
	return parents, err
}

func (r *adminAreaCommonRepo) ListChildren(ctx context.Context, cityID int, limited bool) ([]AreaItem, error) {
	q := r.db.WithContext(ctx).Table("areas").
		Select("id, name, parent_id").
		Where("parent_id > 0 AND is_show = 1 AND deleted_at IS NULL").
		Order("sort ASC, id ASC")
	if limited {
		q = q.Where("parent_id = ?", cityID)
	}
	var children []AreaItem
	err := q.Find(&children).Error
	return children, err
}

type AdminRoleRepository interface {
	CountUsersByRole(ctx context.Context, roleID int) (int64, error)
	ListPermissionLinks(ctx context.Context, roleIDs []int) ([]RolePermLink, error)
}

type AdminStaffRepository interface {
	CountRole(ctx context.Context, roleID int) (int64, error)
	CountRootArea(ctx context.Context, id int) (int64, error)
	ListAreasByIDs(ctx context.Context, ids []int) ([]AreaItem, error)
}

type RolePermLink struct {
	RoleID       int `gorm:"column:roleId"`
	PermissionID int `gorm:"column:permissionId"`
}

type adminRoleRepo struct{ db *gorm.DB }

func NewAdminRole(db *gorm.DB) AdminRoleRepository { return &adminRoleRepo{db: db} }

func (r *adminRoleRepo) CountUsersByRole(ctx context.Context, roleID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("administrator_role_relation").Where("roleId = ?", roleID).Count(&n).Error
	return n, err
}

func (r *adminRoleRepo) ListPermissionLinks(ctx context.Context, roleIDs []int) ([]RolePermLink, error) {
	var links []RolePermLink
	if len(roleIDs) == 0 {
		return links, nil
	}
	err := r.db.WithContext(ctx).Table("administrator_permissions").Where("roleId IN ?", roleIDs).Find(&links).Error
	return links, err
}

type adminStaffRepo struct{ db *gorm.DB }

func NewAdminStaff(db *gorm.DB) AdminStaffRepository { return &adminStaffRepo{db: db} }

func (r *adminStaffRepo) CountRole(ctx context.Context, roleID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("administrator_roles").Where("id = ?", roleID).Count(&n).Error
	return n, err
}

func (r *adminStaffRepo) CountRootArea(ctx context.Context, id int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("areas").Where("id = ? AND parent_id = 0 AND deleted_at IS NULL", id).Count(&n).Error
	return n, err
}

func (r *adminStaffRepo) ListAreasByIDs(ctx context.Context, ids []int) ([]AreaItem, error) {
	var found []AreaItem
	if len(ids) == 0 {
		return found, nil
	}
	err := r.db.WithContext(ctx).Table("areas").Select("id, name").Where("id IN ?", ids).Find(&found).Error
	return found, err
}

type AdminMenuRepository interface {
	CountGranted(ctx context.Context, permissionID int) (int64, error)
}

type adminMenuRepo struct{ db *gorm.DB }

func NewAdminMenu(db *gorm.DB) AdminMenuRepository { return &adminMenuRepo{db: db} }

func (r *adminMenuRepo) CountGranted(ctx context.Context, permissionID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("administrator_permissions").Where("permissionId = ?", permissionID).Count(&n).Error
	return n, err
}

type AdminLoginRepository interface {
	CountRoles(ctx context.Context, adminID int) (int64, error)
	IsSuper(ctx context.Context, adminID int) (bool, error)
}

type adminLoginRepo struct{ db *gorm.DB }

func NewAdminLogin(db *gorm.DB) AdminLoginRepository { return &adminLoginRepo{db: db} }

func (r *adminLoginRepo) CountRoles(ctx context.Context, adminID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("administrator_role_relation").Where("administratorId = ?", adminID).Count(&n).Error
	return n, err
}

func (r *adminLoginRepo) IsSuper(ctx context.Context, adminID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("administrator_role_relation AS rel").
		Joins("JOIN administrator_roles role ON role.id = rel.roleId").
		Where("rel.administratorId = ? AND role.code = ?", adminID, "SUPER_ADMIN").
		Count(&count).Error
	return count > 0, err
}
