package admin

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// AdminFeedbackRepository covers admin feedback Table-based access.
type AdminFeedbackRepository interface {
	ListUsersBrief(ctx context.Context, ids []int) ([]AdminFeedbackUserBrief, error)
	ListAdminsBrief(ctx context.Context, ids []int) ([]AdminFeedbackAdminBrief, error)
	CreateNotice(ctx context.Context, userID, adminID int, title, content, relatedType string, relatedID int, now time.Time) error
}

type AdminFeedbackUserBrief struct {
	ID       int     `gorm:"column:id" json:"id"`
	NickName *string `gorm:"column:nick_name" json:"nick_name"`
	Avatar   string  `gorm:"column:avatar" json:"avatar"`
}

type AdminFeedbackAdminBrief struct {
	ID       int     `gorm:"column:id" json:"id"`
	NickName *string `gorm:"column:nickName" json:"nickName"`
}

type adminFeedbackRepo struct{ db *gorm.DB }

func NewAdminFeedback(db *gorm.DB) AdminFeedbackRepository {
	return &adminFeedbackRepo{db: db}
}

func (r *adminFeedbackRepo) ListUsersBrief(ctx context.Context, ids []int) ([]AdminFeedbackUserBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminFeedbackUserBrief
	err := r.db.WithContext(ctx).Table("users").Select("id, nick_name, avatar").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminFeedbackRepo) ListAdminsBrief(ctx context.Context, ids []int) ([]AdminFeedbackAdminBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminFeedbackAdminBrief
	err := r.db.WithContext(ctx).Table("administrator").Select("id, nickName").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminFeedbackRepo) CreateNotice(ctx context.Context, userID, adminID int, title, content, relatedType string, relatedID int, now time.Time) error {
	if userID <= 0 {
		return nil
	}
	return r.db.WithContext(ctx).Table("user_messages").Create(map[string]any{
		"user_id": userID, "admin_id": adminID, "title": title, "content": content,
		"type": "notice", "related_type": relatedType, "related_id": relatedID,
		"is_read": 0, "created_at": now, "updated_at": now,
	}).Error
}
