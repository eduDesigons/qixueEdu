package adminbanner

import (
	"context"
	"errors"

	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
)

var (
	ErrNotFound = errors.New("Banner不存在")
	ErrNeedIDs  = errors.New("请传入排序后的 id 列表")
)

type Service struct {
	banners adminrepo.AdminBannerRepository
}

func New(banners adminrepo.AdminBannerRepository) *Service {
	return &Service{banners: banners}
}

type View struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Sort      int     `json:"sort"`
	ImgURL    string  `json:"img_url"`
	LinkURL   *string `json:"link_url"`
	Type      int     `json:"type"`
	Status    int     `json:"status"`
	Content   *string `json:"content"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type SaveInput struct {
	ID      int
	Name    string
	ImgURL  string
	Sort    int
	LinkURL *string
	Type    int
	Status  int
	Content *string
}

func (s *Service) List(ctx context.Context, typ, status *int) ([]View, error) {
	rows, err := s.banners.List(ctx, typ, status)
	if err != nil {
		return nil, err
	}
	list := make([]View, 0, len(rows))
	for _, item := range rows {
		list = append(list, toView(item))
	}
	return list, nil
}

func (s *Service) Save(ctx context.Context, in SaveInput) error {
	fields := map[string]any{
		"name": in.Name, "img_url": in.ImgURL, "sort": in.Sort,
		"link_url": in.LinkURL, "type": in.Type, "status": in.Status, "content": in.Content,
	}
	if in.ID > 0 {
		n, err := s.banners.Update(ctx, in.ID, fields)
		if err != nil {
			return err
		}
		if n == 0 {
			count, err := s.banners.Count(ctx, in.ID)
			if err != nil {
				return err
			}
			if count == 0 {
				return ErrNotFound
			}
		}
		return nil
	}
	row := &model.Banner{
		Name: in.Name, ImgURL: in.ImgURL, Sort: in.Sort, LinkURL: in.LinkURL,
		Type: in.Type, Status: in.Status, Content: in.Content,
	}
	return s.banners.Create(ctx, row)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrNotFound
	}
	n, err := s.banners.Delete(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Service) Sort(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return ErrNeedIDs
	}
	return s.banners.UpdateSorts(ctx, ids)
}

func toView(item model.Banner) View {
	return View{
		ID: int(item.ID), Name: item.Name, Sort: item.Sort, ImgURL: item.ImgURL, LinkURL: item.LinkURL,
		Type: item.Type, Status: item.Status, Content: item.Content,
		CreatedAt: model.FormatTime(item.CreatedAt), UpdatedAt: model.FormatTime(item.UpdatedAt),
	}
}
