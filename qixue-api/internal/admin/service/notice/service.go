package adminnotice

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
)

var (
	ErrNotFound  = errors.New("通知不存在")
	ErrBadTime   = errors.New("发布时间格式错误")
	ErrNeedID    = errors.New("ID不能为空")
	ErrBadStatus = errors.New("状态参数错误，1=已发布，2=未发布")
)

type Service struct {
	notices adminrepo.AdminNoticeRepository
}

func New(notices adminrepo.AdminNoticeRepository) *Service {
	return &Service{notices: notices}
}

type View struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Status      int     `json:"status"`
	PublishTime string  `json:"publish_time"`
	AdminID     int     `json:"admin_id"`
	Content     string  `json:"content"`
	ViewTimes   int     `json:"view_times"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

type SaveInput struct {
	ID          int
	Title       string
	PublishTime string
	Content     string
	Status      int
}

func (s *Service) List(ctx context.Context, status *int, keyword string, page, size int) ([]View, int64, error) {
	offset := (page - 1) * size
	rows, total, err := s.notices.List(ctx, status, strings.TrimSpace(keyword), offset, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]View, 0, len(rows))
	for _, item := range rows {
		list = append(list, toView(item))
	}
	return list, total, nil
}

func (s *Service) Detail(ctx context.Context, id int) (*View, error) {
	if id <= 0 {
		return nil, ErrNeedID
	}
	item, err := s.notices.Find(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	v := toView(*item)
	return &v, nil
}

func (s *Service) Save(ctx context.Context, in SaveInput) error {
	published, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(in.PublishTime), time.Local)
	if err != nil {
		return ErrBadTime
	}
	fields := map[string]any{
		"title": in.Title, "publish_time": published, "content": in.Content, "status": in.Status,
	}
	if in.ID > 0 {
		n, err := s.notices.Update(ctx, in.ID, fields)
		if err != nil {
			return err
		}
		if n == 0 {
			if _, err := s.notices.Find(ctx, in.ID); errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			} else if err != nil {
				return err
			}
		}
		return nil
	}
	row := &model.Notice{Title: in.Title, PublishTime: published, Content: in.Content, Status: in.Status}
	return s.notices.Create(ctx, row)
}

func (s *Service) Status(ctx context.Context, id, status int) error {
	if id <= 0 {
		return ErrNeedID
	}
	if status != 1 && status != 2 {
		return ErrBadStatus
	}
	n, err := s.notices.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrNeedID
	}
	n, err := s.notices.Delete(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func toView(item model.Notice) View {
	return View{
		ID: int(item.ID), Title: item.Title, Status: item.Status,
		PublishTime: item.PublishTime.Format("2006-01-02 15:04:05"),
		AdminID:     item.AdminID, Content: item.Content, ViewTimes: item.ViewTimes,
		CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
	}
}
