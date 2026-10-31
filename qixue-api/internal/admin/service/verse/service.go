package adminverse

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	adminrepo "qixue-api/internal/admin/repo"
)

var (
	ErrEmpty   = errors.New("请填写一言")
	ErrTooLong = errors.New("一言不能超过 500 个字")
	ErrBadID   = errors.New("参数错误")
	ErrMissing = errors.New("一言不存在")
)

type Service struct {
	verses adminrepo.AdminVerseRepository
}

func New(verses adminrepo.AdminVerseRepository) *Service {
	return &Service{verses: verses}
}

func (s *Service) List(ctx context.Context, keyword string, page, size int) ([]adminrepo.VerseRow, int64, error) {
	offset := (page - 1) * size
	list, total, err := s.verses.List(ctx, strings.TrimSpace(keyword), offset, size)
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []adminrepo.VerseRow{}
	}
	return list, total, nil
}

func (s *Service) Save(ctx context.Context, id int, content string) (int, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return 0, ErrEmpty
	}
	if utf8.RuneCountInString(content) > 500 {
		return 0, ErrTooLong
	}
	if id > 0 {
		n, err := s.verses.UpdateContent(ctx, id, content)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			count, err := s.verses.Count(ctx, id)
			if err != nil {
				return 0, err
			}
			if count == 0 {
				return 0, ErrMissing
			}
		}
		return id, nil
	}
	row := &adminrepo.VerseRow{Content: content}
	if err := s.verses.Create(ctx, row); err != nil {
		return 0, err
	}
	return row.ID, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrBadID
	}
	n, err := s.verses.Delete(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrMissing
	}
	return nil
}
