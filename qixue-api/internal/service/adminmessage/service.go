package adminmessage

import (
	"context"
	"html"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/scope"
)

const (
	typeNotice = "notice"
	typePaper  = "paper"
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

type userBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
}

type adminBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nickName"`
}

type messageView struct {
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
	User        *userBrief  `json:"user"`
	Admin       *adminBrief `json:"admin"`
}

func (s *Service) Send(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Type        string `json:"type"`
		RelatedType any    `json:"related_type"`
		RelatedID   any    `json:"related_id"`
		UserIDs     any    `json:"user_ids"`
		UserID      any    `json:"user_id"`
	}
	_ = c.ShouldBindJSON(&body)
	ctx := c.Request.Context()

	title := strings.TrimSpace(body.Title)
	content := strings.TrimSpace(body.Content)
	msgType := strings.TrimSpace(body.Type)
	if msgType == "" {
		msgType = typeNotice
	}

	if title == "" {
		response.Fail(c, 1, "请填写标题")
		return
	}
	if utf8.RuneCountInString(title) > 80 {
		response.Fail(c, 1, "标题不能超过80个字")
		return
	}
	if plainContent(content) == "" {
		response.Fail(c, 1, "请填写内容")
		return
	}
	if utf8.RuneCountInString(content) > 20000 {
		response.Fail(c, 1, "内容过长，请精简后再发送")
		return
	}
	if msgType != typeNotice && msgType != typePaper {
		msgType = typeNotice
	}

	userIDs := asInts(body.UserIDs)
	if len(userIDs) == 0 {
		if uid := asInt(body.UserID); uid > 0 {
			userIDs = []int{uid}
		}
	}
	userIDs = uniquePositive(userIDs)
	if len(userIDs) == 0 {
		response.Fail(c, 1, "请选择学员")
		return
	}

	existIDs, err := s.messages.PluckExistingUserIDs(ctx, userIDs)
	if response.WriteErr(c, err) {
		return
	}
	if len(existIDs) != len(userIDs) {
		response.Fail(c, 1, "部分学员不存在")
		return
	}

	adminID := middleware.UID(c)
	now := time.Now()
	relatedType := asOptionalString(body.RelatedType)
	relatedID := asOptionalRelatedID(body.RelatedID)
	rows := make([]map[string]any, 0, len(userIDs))
	for _, uid := range userIDs {
		row := map[string]any{
			"user_id": uid, "admin_id": adminID, "title": title, "content": content,
			"type": msgType, "is_read": 0, "read_at": nil,
			"created_at": now, "updated_at": now,
		}
		if relatedType != nil {
			row["related_type"] = *relatedType
		} else {
			row["related_type"] = nil
		}
		if relatedID != nil {
			row["related_id"] = *relatedID
		} else {
			row["related_id"] = nil
		}
		rows = append(rows, row)
	}
	if response.WriteErr(c, s.messages.CreateMessages(ctx, rows)) {
		return
	}
	response.OK(c, "发送成功", gin.H{"count": len(rows)})
}

func (s *Service) List(c *gin.Context) {
	var body struct {
		Title    string   `json:"title"`
		NickName string   `json:"nick_name"`
		Type     string   `json:"type"`
		IsRead   any      `json:"is_read"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}

	query := s.db.WithContext(ctx).Model(&messageRow{})
	if title := strings.TrimSpace(body.Title); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if nick := strings.TrimSpace(body.NickName); nick != "" {
		query = query.Where(
			"user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)",
			"%"+nick+"%",
		)
	}
	if t := strings.TrimSpace(body.Type); t != "" {
		query = query.Where("type = ?", t)
	}
	if isRead, ok := asOptionalInt(body.IsRead); ok {
		query = query.Where("is_read = ?", isRead)
	}
	if limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}

	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []messageRow
	if response.WriteErr(c, query.Order("id DESC").Offset((current-1)*size).Limit(size).Find(&rows).Error) {
		return
	}
	list, err := s.decorate(ctx, rows)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "ok", page.New(current, size, total, list))
}

func (s *Service) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	ctx := c.Request.Context()

	var row messageRow
	err := s.db.WithContext(ctx).Where("id = ?", body.ID.Int()).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "记录不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	// user_messages has no SoftDeletes
	if response.WriteErr(c, s.db.WithContext(ctx).Where("id = ?", row.ID).Delete(&messageRow{}).Error) {
		return
	}
	response.OK(c, "删除成功", nil)
}

func (s *Service) decorate(ctx context.Context, rows []messageRow) ([]messageView, error) {
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
	list := make([]messageView, 0, len(rows))
	for _, r := range rows {
		var user *userBrief
		if u, ok := users[r.UserID]; ok {
			cp := u
			user = &cp
		}
		var admin *adminBrief
		if a, ok := admins[r.AdminID]; ok {
			cp := a
			admin = &cp
		}
		list = append(list, messageView{
			ID: r.ID, UserID: r.UserID, AdminID: r.AdminID, Title: r.Title, Content: r.Content,
			Type: r.Type, RelatedType: r.RelatedType, RelatedID: r.RelatedID, IsRead: r.IsRead,
			ReadAt: model.FormatTime(r.ReadAt), CreatedAt: model.FormatTime(r.CreatedAt),
			UpdatedAt: model.FormatTime(r.UpdatedAt), User: user, Admin: admin,
		})
	}
	return list, nil
}

func (s *Service) loadUsers(ctx context.Context, ids []int) (map[int]userBrief, error) {
	out := map[int]userBrief{}
	rows, err := s.messages.ListUsersBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = userBrief{ID: r.ID, NickName: r.NickName, Avatar: r.Avatar}
	}
	return out, nil
}

func (s *Service) loadAdmins(ctx context.Context, ids []int) (map[int]adminBrief, error) {
	out := map[int]adminBrief{}
	rows, err := s.messages.ListAdminsBrief(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = adminBrief{ID: r.ID, NickName: r.NickName}
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

func asInt(v any) int { return bind.AsInt(v) }

func asOptionalInt(v any) (int, bool) {
	if !bind.Present(v) {
		return 0, false
	}
	return bind.AsInt(v), true
}

func asOptionalString(v any) *string {
	if v == nil {
		return nil
	}
	switch n := v.(type) {
	case string:
		s := strings.TrimSpace(n)
		if s == "" {
			return nil
		}
		return &s
	default:
		return nil
	}
}

func asOptionalRelatedID(v any) *int {
	if v == nil {
		return nil
	}
	id := asInt(v)
	if id <= 0 {
		return nil
	}
	return &id
}

func asInts(v any) []int {
	switch n := v.(type) {
	case []any:
		out := make([]int, 0, len(n))
		for _, item := range n {
			if id := asInt(item); id > 0 {
				out = append(out, id)
			}
		}
		return out
	case []int:
		return n
	case []float64:
		out := make([]int, 0, len(n))
		for _, item := range n {
			if id := int(item); id > 0 {
				out = append(out, id)
			}
		}
		return out
	default:
		return nil
	}
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
