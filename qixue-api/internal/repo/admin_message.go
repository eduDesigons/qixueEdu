package repo

import (
	"context"

	"gorm.io/gorm"
)

// AdminMessageRepository covers admin message Table-based access.
type AdminMessageRepository interface {
	PluckExistingUserIDs(ctx context.Context, ids []int) ([]int, error)
	CreateMessages(ctx context.Context, rows []map[string]any) error
	ListUsersBrief(ctx context.Context, ids []int) ([]AdminMsgUserBrief, error)
	ListAdminsBrief(ctx context.Context, ids []int) ([]AdminMsgAdminBrief, error)
}

type AdminMsgUserBrief struct {
	ID       int     `gorm:"column:id" json:"id"`
	NickName *string `gorm:"column:nick_name" json:"nick_name"`
	Avatar   string  `gorm:"column:avatar" json:"avatar"`
}

type AdminMsgAdminBrief struct {
	ID       int     `gorm:"column:id" json:"id"`
	NickName *string `gorm:"column:nickName" json:"nickName"`
}

type adminMessageRepo struct{ db *gorm.DB }

func NewAdminMessage(db *gorm.DB) AdminMessageRepository {
	return &adminMessageRepo{db: db}
}

func (r *adminMessageRepo) PluckExistingUserIDs(ctx context.Context, ids []int) ([]int, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var existIDs []int
	err := r.db.WithContext(ctx).Table("users").Where("id IN ?", ids).Pluck("id", &existIDs).Error
	return existIDs, err
}

func (r *adminMessageRepo) CreateMessages(ctx context.Context, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Table("user_messages").Create(&rows).Error
}

func (r *adminMessageRepo) ListUsersBrief(ctx context.Context, ids []int) ([]AdminMsgUserBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminMsgUserBrief
	err := r.db.WithContext(ctx).Table("users").Select("id, nick_name, avatar").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminMessageRepo) ListAdminsBrief(ctx context.Context, ids []int) ([]AdminMsgAdminBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminMsgAdminBrief
	err := r.db.WithContext(ctx).Table("administrator").Select("id, nickName").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}
