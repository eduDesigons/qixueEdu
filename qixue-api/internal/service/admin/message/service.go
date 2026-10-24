package adminmessage

import (
	"context"
	"errors"
	"html"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/scope"
)

const (
	typeNotice = "notice"
	typePaper  = "paper"
)

var (
	ErrEmptyTitle   = errors.New("请填写标题")
	ErrTitleTooLong = errors.New("标题不能超过80个字")
	ErrEmptyContent = errors.New("请填写内容")
	ErrContentLong  = errors.New("内容过长，请精简后再发送")
	ErrNoUsers      = errors.New("请选择学员")
	ErrUsersMissing = errors.New("部分学员不存在")
	ErrNotFound     = errors.New("记录不存在")
)

var tagRe = regexp.MustCompile(`(?s)<[^>]*>`)

type Service struct {
	messages repo.AdminMessageRepository
	db       *gorm.DB
	rdb      *redis.Client
}

func New(messages repo.AdminMessageRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{messages: messages, db: db, rdb: rdb}
}

type messageRow struct {
	ID          int        `gorm:"column:id;primaryKey"`
	UserID      int        `gorm:"column:user_id"`
	AdminID     int        `gorm:"column:admin_id"`
	Title       string     `gorm:"column:title"`
	Content     string     `gorm:"column:content"`
	Type        string     `gorm:"column:type"`
	RelatedType *string    `gorm:"column:related_type"`
	RelatedID   *int       `gorm:"column:related_id"`
	IsRead      int        `gorm:"column:is_read"`
	ReadAt      *time.Time `gorm:"column:read_at"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (messageRow) TableName() string { return "user_messages" }

type UserBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
}

type AdminBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nickName"`
}

type MessageView struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	AdminID     int         `json:"admin_id"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Type        string      `json:"type"`
	RelatedType *string     `json:"related_type"`
	RelatedID   *int        `json:"related_id"`
	IsRead      int         `json:"is_read"`
	ReadAt      *string     `json:"read_at"`
	CreatedAt   *string     `json:"created_at"`
	UpdatedAt   *string     `json:"updated_at"`
	User        *UserBrief  `json:"user"`
	Admin       *AdminBrief `json:"admin"`
}

type SendInput struct {
	Title       string
	Content     string
	Type        string
	RelatedType *string
	RelatedID   *int
	UserIDs     []int
}

type ListQuery struct {
	Title    string
	NickName string
	Type     string
	IsRead   *int
	Page     int
	Size     int
}

func (s *Service) Send(ctx context.Context, adminID int, in SendInput) (int, error) {
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	msgType := strings.TrimSpace(in.Type)
	if msgType == "" {
		msgType = typeNotice
	}

	if title == "" {
		return 0, ErrEmptyTitle
	}
	if utf8.RuneCountInString(title) > 80 {
		return 0, ErrTitleTooLong
	}
	if plainContent(content) == "" {
		return 0, ErrEmptyContent
	}
	if utf8.RuneCountInString(content) > 20000 {
		return 0, ErrContentLong
	}
	if msgType != typeNotice && msgType != typePaper {
		msgType = typeNotice
	}

	userIDs := uniquePositive(in.UserIDs)
	if len(userIDs) == 0 {
		return 0, ErrNoUsers
	}

	existIDs, err := s.messages.PluckExistingUserIDs(ctx, userIDs)
	if err != nil {
		return 0, err
	}
	if len(existIDs) != len(userIDs) {
		return 0, ErrUsersMissing
	}

	now := time.Now()
	rows := make([]map[string]any, 0, len(userIDs))
	for _, uid := range userIDs {
		row := map[string]any{
			"user_id": uid, "admin_id": adminID, "title": title, "content": content,
			"type": msgType, "is_read": 0, "read_at": nil,
			"created_at": now, "updated_at": now,
		}
		if in.RelatedType != nil {
			row["related_type"] = *in.RelatedType
		} else {
			row["related_type"] = nil
		}
		if in.RelatedID != nil {
			row["related_id"] = *in.RelatedID
		} else {
			row["related_id"] = nil
		}
		rows = append(rows, row)
	}
	if err := s.messages.CreateMessages(ctx, rows); err != nil {
		return 0, err
	}
	return len(rows), nil
}

func (s *Service) List(ctx context.Context, adminID int, q ListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}

	query := s.db.WithContext(ctx).Model(&messageRow{})
	if title := strings.TrimSpace(q.Title); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if nick := strings.TrimSpace(q.NickName); nick != "" {
		query = query.Where(
			"user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)",
			"%"+nick+"%",
		)
	}
	if t := strings.TrimSpace(q.Type); t != "" {
		query = query.Where("type = ?", t)
	}
	if q.IsRead != nil {
		query = query.Where("is_read = ?", *q.IsRead)
	}
	if limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []messageRow
	if err := query.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	list, err := s.decorate(ctx, rows)
	if err != nil {
		return page.Body{}, err
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	var row messageRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Where("id = ?", row.ID).Delete(&messageRow{}).Error
}

func (s *Service) decorate(ctx context.Context, rows []messageRow) ([]MessageView, error) {
	userIDs, adminIDs := []int{}, []int{}
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		if r.AdminID > 0 {
			adminIDs = append(adminIDs, r.AdminID)
		}
	}
	users, err := s.loadUsers(ctx, uniquePositive(userIDs))
	if err != nil {
		return nil, err
	}
	admins, err := s.loadAdmins(ctx, uniquePositive(adminIDs))
	if err != nil {
		return nil, err
	}
	list := make([]MessageView, 0, len(rows))
	for _, r := range rows {
		var user *UserBrief
		if u, ok := users[r.UserID]; ok {
			cp := u
			user = &cp
		}
		var admin *AdminBrief
		if a, ok := admins[r.AdminID]; ok {
			cp := a
			admin = &cp
		}
		list = append(list, MessageView{
			ID: r.ID, UserID: r.UserID, AdminID: r.AdminID, Title: r.Title, Content: r.Content,
			Type: r.Type, RelatedType: r.RelatedType, RelatedID: r.RelatedID, IsRead: r.IsRead,
			ReadAt: model.FormatTime(r.ReadAt), CreatedAt: model.FormatTime(r.CreatedAt),
			UpdatedAt: model.FormatTime(r.UpdatedAt), User: user, Admin: admin,
		})
	}
	return list, nil
}

func (s *Service) loadUsers(ctx context.Context, ids []int) (map[int]UserBrief, error) {
	out := map[int]UserBrief{}
	rows, err := s.messages.ListUsersBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = UserBrief{ID: r.ID, NickName: r.NickName, Avatar: r.Avatar}
	}
	return out, nil
}

func (s *Service) loadAdmins(ctx context.Context, ids []int) (map[int]AdminBrief, error) {
	out := map[int]AdminBrief{}
	rows, err := s.messages.ListAdminsBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = AdminBrief{ID: r.ID, NickName: r.NickName}
	}
	return out, nil
}

func plainContent(content string) string {
	text := tagRe.ReplaceAllString(content, "")
	text = html.UnescapeString(text)
	text = strings.TrimSpace(text)
	var b strings.Builder
	for _, r := range text {
		if unicode.IsSpace(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func uniquePositive(ids []int) []int {
	seen := map[int]struct{}{}
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
