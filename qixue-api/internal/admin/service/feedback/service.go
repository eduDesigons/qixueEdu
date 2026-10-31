package adminfeedback

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/scope"
)

const statusReplied = 1

var (
	ErrNotFound     = errors.New("记录不存在")
	ErrEmptyReply   = errors.New("请填写回复")
	ErrReplyTooLong = errors.New("回复不能超过1000字")
)

type Service struct {
	feedbacks adminrepo.AdminFeedbackRepository
	db        *gorm.DB
	rdb       *redis.Client
}

func New(feedbacks adminrepo.AdminFeedbackRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{feedbacks: feedbacks, db: db, rdb: rdb}
}

type feedbackRow struct {
	ID           int             `gorm:"column:id;primaryKey"`
	UserID       int             `gorm:"column:user_id"`
	Content      string          `gorm:"column:content"`
	Images       json.RawMessage `gorm:"column:images"`
	Status       int             `gorm:"column:status"`
	Reply        *string         `gorm:"column:reply"`
	ReplyAdminID int             `gorm:"column:reply_admin_id"`
	RepliedAt    *time.Time      `gorm:"column:replied_at"`
	CreatedAt    *time.Time      `gorm:"column:created_at"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at"`
}

func (feedbackRow) TableName() string { return "user_feedbacks" }

type UserBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
}

type AdminBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nickName"`
}

type FeedbackView struct {
	ID           int         `json:"id"`
	UserID       int         `json:"user_id"`
	Content      string      `json:"content"`
	Images       any         `json:"images"`
	Status       int         `json:"status"`
	Reply        *string     `json:"reply"`
	ReplyAdminID int         `json:"reply_admin_id"`
	RepliedAt    *string     `json:"replied_at"`
	CreatedAt    *string     `json:"created_at"`
	UpdatedAt    *string     `json:"updated_at"`
	User         *UserBrief  `json:"user"`
	Admin        *AdminBrief `json:"admin"`
}

type ListQuery struct {
	NickName string
	Status   *int
	CityID   int
	Page     int
	Size     int
}

func (s *Service) resolveCity(ctx context.Context, adminID, requestCityID int) (int, error) {
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return 0, err
	}
	if limited {
		return cityID, nil
	}
	if requestCityID > 0 {
		return requestCityID, nil
	}
	return 0, nil
}

func (s *Service) lookupUserIDs(ctx context.Context, nick string, cityID int) ([]int, error) {
	query := s.db.WithContext(ctx).Table("users").Select("id")
	if nick != "" {
		query = query.Where("nick_name LIKE ?", "%"+nick+"%")
	}
	if cityID > 0 {
		query = query.Where("city_id = ?", cityID)
	}
	var ids []int
	err := query.Limit(2000).Scan(&ids).Error
	return ids, err
}

func (s *Service) List(ctx context.Context, adminID int, q ListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, err := s.resolveCity(ctx, adminID, q.CityID)
	if err != nil {
		return page.Body{}, err
	}

	query := s.db.WithContext(ctx).Table("user_feedbacks AS f")
	nick := strings.TrimSpace(q.NickName)
	if nick != "" {
		userIDs, err := s.lookupUserIDs(ctx, nick, cityID)
		if err != nil {
			return page.Body{}, err
		}
		if len(userIDs) == 0 {
			return page.New(current, size, 0, []FeedbackView{}), nil
		}
		query = query.Where("f.user_id IN ?", userIDs)
	} else if cityID > 0 {
		query = query.Joins("INNER JOIN users u ON u.id = f.user_id AND u.city_id = ?", cityID)
	}
	if q.Status != nil {
		query = query.Where("f.status = ?", *q.Status)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	if total == 0 {
		return page.New(current, size, 0, []FeedbackView{}), nil
	}

	var ids []int
	if err := query.Session(&gorm.Session{}).
		Select("f.id").
		Order("f.id DESC").
		Offset((current - 1) * size).
		Limit(size).
		Scan(&ids).Error; err != nil {
		return page.Body{}, err
	}
	if len(ids) == 0 {
		return page.New(current, size, total, []FeedbackView{}), nil
	}

	var rows []feedbackRow
	if err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	byID := make(map[int]feedbackRow, len(rows))
	for _, row := range rows {
		byID[row.ID] = row
	}
	ordered := make([]feedbackRow, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			ordered = append(ordered, row)
		}
	}
	list, err := s.decorate(ctx, ordered)
	if err != nil {
		return page.Body{}, err
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) Reply(ctx context.Context, adminID, id int, reply string) error {
	reply = strings.TrimSpace(reply)
	var row feedbackRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if reply == "" {
		return ErrEmptyReply
	}
	if utf8.RuneCountInString(reply) > 1000 {
		return ErrReplyTooLong
	}

	now := time.Now()
	updates := map[string]any{
		"reply":          reply,
		"status":         statusReplied,
		"reply_admin_id": adminID,
		"replied_at":     now,
		"updated_at":     now,
	}
	if err := s.db.WithContext(ctx).Model(&feedbackRow{}).Where("id = ?", row.ID).Updates(updates).Error; err != nil {
		return err
	}

	safe := html.EscapeString(reply)
	content := "<p>" + strings.ReplaceAll(safe, "\n", "<br />\n") + "</p>"
	return s.feedbacks.CreateNotice(ctx, row.UserID, adminID, "意见建议已回复", content, "feedback", row.ID, now)
}

func (s *Service) decorate(ctx context.Context, rows []feedbackRow) ([]FeedbackView, error) {
	userIDs, adminIDs := []int{}, []int{}
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
		if r.ReplyAdminID > 0 {
			adminIDs = append(adminIDs, r.ReplyAdminID)
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
	list := make([]FeedbackView, 0, len(rows))
	for _, r := range rows {
		var user *UserBrief
		if u, ok := users[r.UserID]; ok {
			cp := u
			user = &cp
		}
		var admin *AdminBrief
		if a, ok := admins[r.ReplyAdminID]; ok {
			cp := a
			admin = &cp
		}
		list = append(list, FeedbackView{
			ID: r.ID, UserID: r.UserID, Content: r.Content, Images: decodeImages(r.Images),
			Status: r.Status, Reply: r.Reply, ReplyAdminID: r.ReplyAdminID,
			RepliedAt: model.FormatTime(r.RepliedAt), CreatedAt: model.FormatTime(r.CreatedAt),
			UpdatedAt: model.FormatTime(r.UpdatedAt), User: user, Admin: admin,
		})
	}
	return list, nil
}

func (s *Service) loadUsers(ctx context.Context, ids []int) (map[int]UserBrief, error) {
	out := map[int]UserBrief{}
	rows, err := s.feedbacks.ListUsersBrief(ctx, ids)
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
	rows, err := s.feedbacks.ListAdminsBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = AdminBrief{ID: r.ID, NickName: r.NickName}
	}
	return out, nil
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
