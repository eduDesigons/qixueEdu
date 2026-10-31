package webapp

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// FeedbackRepository covers learner feedback submissions.
type FeedbackRepository interface {
	CountSince(ctx context.Context, uid int, since string) (int64, error)
	Create(ctx context.Context, uid int, content, images string, status int, now time.Time) (id int, err error)
	ListByUser(ctx context.Context, uid, limit int) ([]FeedbackRow, error)
}

type FeedbackRow struct {
	ID        int             `gorm:"column:id"`
	Content   string          `gorm:"column:content"`
	Images    json.RawMessage `gorm:"column:images"`
	Status    int             `gorm:"column:status"`
	Reply     *string         `gorm:"column:reply"`
	RepliedAt *time.Time      `gorm:"column:replied_at"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
}

type feedbackRepo struct{ db *gorm.DB }

func NewFeedback(db *gorm.DB) FeedbackRepository {
	return &feedbackRepo{db: db}
}

func (r *feedbackRepo) CountSince(ctx context.Context, uid int, since string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_feedbacks").
		Where("user_id = ? AND created_at >= ?", uid, since).Count(&n).Error
	return n, err
}

func (r *feedbackRepo) Create(ctx context.Context, uid int, content, images string, status int, now time.Time) (int, error) {
	row := struct {
		ID        int       `gorm:"column:id;primaryKey"`
		UserID    int       `gorm:"column:user_id"`
		Content   string    `gorm:"column:content"`
		Images    string    `gorm:"column:images"`
		Status    int       `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}{
		UserID: uid, Content: content, Images: images,
		Status: status, CreatedAt: now, UpdatedAt: now,
	}
	err := r.db.WithContext(ctx).Table("user_feedbacks").Create(&row).Error
	return row.ID, err
}

func (r *feedbackRepo) ListByUser(ctx context.Context, uid, limit int) ([]FeedbackRow, error) {
	var rows []FeedbackRow
	err := r.db.WithContext(ctx).Table("user_feedbacks").
		Select("id, content, images, status, reply, replied_at, created_at").
		Where("user_id = ?", uid).Order("id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}
