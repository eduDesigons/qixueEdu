package feedback

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"qixue-api/internal/model"
	webapprepo "qixue-api/internal/webapp/repo"
)

const statusPending = 0

var (
	ErrEmptyContent = errors.New("请填写内容")
	ErrTooLong      = errors.New("内容不能超过500字")
	ErrTooManyImgs  = errors.New("最多上传3张图片")
	ErrDailyLimit   = errors.New("今天提交次数已达上限")
)

var imageExtRe = regexp.MustCompile(`(?i)\.(jpe?g|png|gif|webp)$`)
var httpRe = regexp.MustCompile(`(?i)^https?:`)

type Service struct {
	feedbacks webapprepo.FeedbackRepository
}

func New(feedbacks webapprepo.FeedbackRepository) *Service {
	return &Service{feedbacks: feedbacks}
}

func (s *Service) Submit(ctx context.Context, uid int, content string, images any) (int, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return 0, ErrEmptyContent
	}
	if utf8.RuneCountInString(content) > 500 {
		return 0, ErrTooLong
	}
	imgs, ok := normalizeImages(images)
	if !ok {
		return 0, ErrTooManyImgs
	}
	today := time.Now().Format("2006-01-02") + " 00:00:00"
	count, err := s.feedbacks.CountSince(ctx, uid, today)
	if err != nil {
		return 0, err
	}
	if count >= 5 {
		return 0, ErrDailyLimit
	}
	raw, _ := json.Marshal(imgs)
	now := time.Now()
	return s.feedbacks.Create(ctx, uid, content, string(raw), statusPending, now)
}

func (s *Service) List(ctx context.Context, uid int) ([]map[string]any, error) {
	rows, err := s.feedbacks.ListByUser(ctx, uid, 50)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		list = append(list, map[string]any{
			"id": r.ID, "content": r.Content, "images": decodeImages(r.Images),
			"status": r.Status, "reply": r.Reply,
			"replied_at": model.FormatTime(r.RepliedAt), "created_at": model.FormatTime(r.CreatedAt),
		})
	}
	return list, nil
}

func normalizeImages(images any) ([]string, bool) {
	if images == nil {
		return []string{}, true
	}
	switch v := images.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return []string{}, true
		}
		return nil, false
	case []any:
		out := make([]string, 0, len(v))
		seen := map[string]struct{}{}
		for _, item := range v {
			path := strings.TrimSpace(flatten(item))
			if path == "" {
				continue
			}
			if i := strings.IndexByte(path, '?'); i >= 0 {
				path = path[:i]
			}
			if !imageExtRe.MatchString(path) {
				continue
			}
			if strings.Contains(path, "..") || httpRe.MatchString(path) {
				continue
			}
			path = filepath.ToSlash(path)
			if _, ok := seen[path]; ok {
				continue
			}
			seen[path] = struct{}{}
			out = append(out, path)
		}
		if len(out) > 3 {
			return nil, false
		}
		return out, true
	case []string:
		anyList := make([]any, 0, len(v))
		for _, s := range v {
			anyList = append(anyList, s)
		}
		return normalizeImages(anyList)
	default:
		return nil, false
	}
}

func flatten(v any) string {
	switch n := v.(type) {
	case string:
		return n
	default:
		return ""
	}
}

func decodeImages(raw json.RawMessage) any {
	if len(raw) == 0 || string(raw) == "null" {
		return []any{}
	}
	var arr any
	if json.Unmarshal(raw, &arr) == nil {
		return arr
	}
	return []any{}
}
