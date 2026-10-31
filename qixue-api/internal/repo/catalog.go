package repo

import (
	"context"

	"gorm.io/gorm"
)

// CatalogRepository loads public course/path category trees.
type CatalogRepository interface {
	ListShownCourseParents(ctx context.Context) ([]CourseCategoryRow, error)
	ListShownCourseChildren(ctx context.Context) ([]CourseCategoryRow, error)
	ListPathRoots(ctx context.Context) ([]PathCategoryRow, error)
}

type CourseCategoryRow struct {
	ID       int     `gorm:"column:id"`
	Name     string  `gorm:"column:name"`
	Icon     *string `gorm:"column:icon"`
	ParentID int     `gorm:"column:parent_id"`
	Sort     int     `gorm:"column:sort"`
}

type PathCategoryRow struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

type catalogRepo struct{ db *gorm.DB }

func NewCatalog(db *gorm.DB) CatalogRepository {
	return &catalogRepo{db: db}
}

func (r *catalogRepo) ListShownCourseParents(ctx context.Context) ([]CourseCategoryRow, error) {
	var rows []CourseCategoryRow
	err := r.db.WithContext(ctx).Table("course_categories").
		Select("id, name, icon, parent_id, sort").
		Where("parent_id = 0 AND is_show = 1").
		Order("sort ASC, id ASC").Find(&rows).Error
	return rows, err
}

func (r *catalogRepo) ListShownCourseChildren(ctx context.Context) ([]CourseCategoryRow, error) {
	var rows []CourseCategoryRow
	err := r.db.WithContext(ctx).Table("course_categories").
		Select("id, name, icon, parent_id, sort").
		Where("parent_id > 0 AND is_show = 1").
		Order("sort ASC, id ASC").Find(&rows).Error
	return rows, err
}

func (r *catalogRepo) ListPathRoots(ctx context.Context) ([]PathCategoryRow, error) {
	var rows []PathCategoryRow
	err := r.db.WithContext(ctx).Table("path_categories").
		Select("id, name").
		Where("parent_id = 0").
		Order("sort ASC, id ASC").Find(&rows).Error
	return rows, err
}
