package repo

import (
	"context"

	"gorm.io/gorm"

	"qixue-api/internal/model"
)

type AdminBannerRepository interface {
	List(ctx context.Context, typ, status *int) ([]model.Banner, error)
	Update(ctx context.Context, id int, fields map[string]any) (int64, error)
	Count(ctx context.Context, id int) (int64, error)
	Create(ctx context.Context, row *model.Banner) error
	Delete(ctx context.Context, id int) (int64, error)
	UpdateSorts(ctx context.Context, ids []int) error
}

type AdminNoticeRepository interface {
	List(ctx context.Context, status *int, keyword string, offset, limit int) ([]model.Notice, int64, error)
	Find(ctx context.Context, id int) (*model.Notice, error)
	Update(ctx context.Context, id int, fields map[string]any) (int64, error)
	Create(ctx context.Context, row *model.Notice) error
	UpdateStatus(ctx context.Context, id, status int) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
}

type adminBannerRepo struct{ db *gorm.DB }
type adminNoticeRepo struct{ db *gorm.DB }

func NewAdminBanner(db *gorm.DB) AdminBannerRepository { return &adminBannerRepo{db: db} }
func NewAdminNotice(db *gorm.DB) AdminNoticeRepository { return &adminNoticeRepo{db: db} }

func (r *adminBannerRepo) List(ctx context.Context, typ, status *int) ([]model.Banner, error) {
	q := r.db.WithContext(ctx).Model(&model.Banner{}).Order("sort ASC, id DESC")
	if typ != nil {
		q = q.Where("type = ?", *typ)
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var rows []model.Banner
	err := q.Find(&rows).Error
	return rows, err
}

func (r *adminBannerRepo) Update(ctx context.Context, id int, fields map[string]any) (int64, error) {
	res := r.db.WithContext(ctx).Model(&model.Banner{}).Where("id = ?", id).Updates(fields)
	return res.RowsAffected, res.Error
}

func (r *adminBannerRepo) Count(ctx context.Context, id int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Banner{}).Where("id = ?", id).Count(&n).Error
	return n, err
}

func (r *adminBannerRepo) Create(ctx context.Context, row *model.Banner) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *adminBannerRepo) Delete(ctx context.Context, id int) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&model.Banner{}, id)
	return res.RowsAffected, res.Error
}

func (r *adminBannerRepo) UpdateSorts(ctx context.Context, ids []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for sort, id := range ids {
			if err := tx.Model(&model.Banner{}).Where("id = ?", id).Update("sort", sort).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *adminNoticeRepo) List(ctx context.Context, status *int, keyword string, offset, limit int) ([]model.Notice, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Notice{})
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.Notice
	err := q.Order("publish_time DESC, id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminNoticeRepo) Find(ctx context.Context, id int) (*model.Notice, error) {
	var row model.Notice
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminNoticeRepo) Update(ctx context.Context, id int, fields map[string]any) (int64, error) {
	res := r.db.WithContext(ctx).Model(&model.Notice{}).Where("id = ?", id).Updates(fields)
	return res.RowsAffected, res.Error
}

func (r *adminNoticeRepo) Create(ctx context.Context, row *model.Notice) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *adminNoticeRepo) UpdateStatus(ctx context.Context, id, status int) (int64, error) {
	res := r.db.WithContext(ctx).Model(&model.Notice{}).Where("id = ?", id).Update("status", status)
	return res.RowsAffected, res.Error
}

func (r *adminNoticeRepo) Delete(ctx context.Context, id int) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&model.Notice{}, id)
	return res.RowsAffected, res.Error
}
