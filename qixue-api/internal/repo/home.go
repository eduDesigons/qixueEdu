package repo

import (
	"context"

	"gorm.io/gorm"

	"qixue-api/internal/model"
)

// HomeRepository covers home boot / hot-list / banner / notice reads.
type HomeRepository interface {
	LatestNoticeMeta(ctx context.Context) (*model.Notice, error)
	FindNotice(ctx context.Context, id int) (*model.Notice, error)
	ListActiveBannersMeta(ctx context.Context) ([]model.Banner, error)
	FindBanner(ctx context.Context, id int) (*model.Banner, error)
	ListHotPaths(ctx context.Context, limit int) ([]HotPathRow, error)
	ListHotCoursesByCate(ctx context.Context, cateID, limit int) ([]HotCourseByCateRow, error)
	ListHotCoursesPage(ctx context.Context, pageNum, limit int) ([]HotCoursePageRow, error)
}

type HotPathRow struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
	Thumb string `gorm:"column:thumb"`
	Desc  string `gorm:"column:desc"`
}

type HotCourseByCateRow struct {
	ID            int    `gorm:"column:id"`
	Title         string `gorm:"column:title"`
	Thumb         string `gorm:"column:thumb"`
	Desc          string `gorm:"column:desc"`
	Type          int    `gorm:"column:type"`
	Credit        int    `gorm:"column:credit"`
	MinNeedCredit int    `gorm:"column:min_need_credit"`
}

type HotCoursePageRow struct {
	ID         int     `gorm:"column:id"`
	Title      string  `gorm:"column:title"`
	Thumb      string  `gorm:"column:thumb"`
	Desc       string  `gorm:"column:desc"`
	CategoryID int     `gorm:"column:category_id"`
	CateName   *string `gorm:"column:cate_name"`
}

type homeRepo struct{ db *gorm.DB }

func NewHome(db *gorm.DB) HomeRepository {
	return &homeRepo{db: db}
}

func (r *homeRepo) LatestNoticeMeta(ctx context.Context) (*model.Notice, error) {
	var row model.Notice
	err := r.db.WithContext(ctx).Table("notices").
		Select("id, title, status, publish_time, admin_id, view_times, created_at, updated_at").
		Where("status = 1 AND deleted_at IS NULL").
		Order("publish_time DESC").Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *homeRepo) FindNotice(ctx context.Context, id int) (*model.Notice, error) {
	var row model.Notice
	err := r.db.WithContext(ctx).Table("notices").
		Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *homeRepo) ListActiveBannersMeta(ctx context.Context) ([]model.Banner, error) {
	var rows []model.Banner
	err := r.db.WithContext(ctx).Table("banners").
		Select("id, name, sort, img_url, link_url, type, status, created_at, updated_at").
		Where("status = 1 AND deleted_at IS NULL").
		Order("sort ASC").Find(&rows).Error
	return rows, err
}

func (r *homeRepo) FindBanner(ctx context.Context, id int) (*model.Banner, error) {
	var row model.Banner
	err := r.db.WithContext(ctx).Table("banners").
		Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *homeRepo) ListHotPaths(ctx context.Context, limit int) ([]HotPathRow, error) {
	var rows []HotPathRow
	err := r.db.WithContext(ctx).Table("paths").
		Select("id, title, thumb, `desc`").
		Where("is_show = 1 AND is_rec = 1 AND deleted_at IS NULL").
		Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *homeRepo) ListHotCoursesByCate(ctx context.Context, cateID, limit int) ([]HotCourseByCateRow, error) {
	if cateID <= 0 {
		return []HotCourseByCateRow{}, nil
	}
	var rows []HotCourseByCateRow
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb, `desc`, type, credit, min_need_credit").
		Where("category_id = ? AND is_show = 1 AND is_rec = 1 AND deleted_at IS NULL", cateID).
		Order("view_count DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *homeRepo) ListHotCoursesPage(ctx context.Context, pageNum, limit int) ([]HotCoursePageRow, error) {
	var rows []HotCoursePageRow
	err := r.db.WithContext(ctx).Table("courses AS c").
		Select("c.id, c.title, c.thumb, c.`desc`, c.category_id, cat.name AS cate_name").
		Joins("LEFT JOIN course_categories cat ON cat.id = c.category_id").
		Where("c.is_show = 1 AND c.is_rec = 1 AND c.deleted_at IS NULL").
		Order("c.id DESC").
		Offset((pageNum - 1) * limit).Limit(limit).
		Find(&rows).Error
	return rows, err
}
