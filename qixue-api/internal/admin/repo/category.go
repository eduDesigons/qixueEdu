package admin

import (
	"context"

	"gorm.io/gorm"
)

type AdminCategoryRepository interface {
	CountPathsByCategory(ctx context.Context, categoryID int) (int64, error)
	CountCoursesByCategory(ctx context.Context, categoryID int) (int64, error)
	CategoryItemCounts(ctx context.Context, table string) (map[int]int, error)
}

type adminCategoryRepo struct{ db *gorm.DB }

func NewAdminCategory(db *gorm.DB) AdminCategoryRepository {
	return &adminCategoryRepo{db: db}
}

func (r *adminCategoryRepo) CountPathsByCategory(ctx context.Context, categoryID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("paths").
		Where("category_id = ? AND deleted_at IS NULL", categoryID).Count(&n).Error
	return n, err
}

func (r *adminCategoryRepo) CountCoursesByCategory(ctx context.Context, categoryID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("courses").
		Where("category_id = ? AND deleted_at IS NULL", categoryID).Count(&n).Error
	return n, err
}

func (r *adminCategoryRepo) CategoryItemCounts(ctx context.Context, table string) (map[int]int, error) {
	var rows []struct {
		CategoryID int `gorm:"column:category_id"`
		Total      int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table(table).
		Select("category_id, COUNT(*) AS total").
		Where("deleted_at IS NULL").
		Group("category_id").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := map[int]int{}
	for _, row := range rows {
		out[row.CategoryID] = row.Total
	}
	return out, nil
}
