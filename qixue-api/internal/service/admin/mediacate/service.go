package adminmediacate

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrBadParams = errors.New("参数错误")
	ErrNotFound  = errors.New("分类不存在")
)

type Service struct{ db *gorm.DB }

func New(db *gorm.DB) *Service { return &Service{db: db} }

type CateView struct {
	ID     int    `json:"id" gorm:"column:id;primaryKey"`
	Name   string `json:"name" gorm:"column:name"`
	Sort   int    `json:"sort" gorm:"column:sort"`
	IsInit int    `json:"is_init" gorm:"column:is_init"`
	Type   string `json:"-" gorm:"column:type"`
}

func (CateView) TableName() string { return "media_cate" }

func (s *Service) List(ctx context.Context, typ string) ([]CateView, error) {
	if typ == "" {
		typ = "image"
	}
	var list []CateView
	err := s.db.WithContext(ctx).Where("type = ?", typ).Order("sort ASC, id ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []CateView{}
	}
	return list, nil
}

func (s *Service) Destroy(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return ErrBadParams
	}
	res := s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&CateView{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
