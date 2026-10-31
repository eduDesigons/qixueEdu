package admin

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CommentTitleBrief struct {
	ID       int    `gorm:"column:id" json:"id"`
	Title    string `gorm:"column:title" json:"title"`
	CourseID int    `gorm:"column:course_id" json:"course_id,omitempty"`
}

type CommentUserBrief struct {
	ID       int     `gorm:"column:id" json:"id"`
	NickName *string `gorm:"column:nick_name" json:"nick_name"`
	Avatar   string  `gorm:"column:avatar" json:"avatar"`
}

type AdminCommentRepository interface {
	ChildCounts(ctx context.Context, table string, ids []int) (map[int]int, error)
	ListUsersBrief(ctx context.Context, ids []int) (map[int]CommentUserBrief, error)
	ListTitles(ctx context.Context, table string, ids []int) (map[int]CommentTitleBrief, error)
	ListChapters(ctx context.Context, ids []int) (map[int]CommentTitleBrief, error)
	IDsByTitleLike(ctx context.Context, table, title string) ([]int, error)
	IDsByNickLike(ctx context.Context, nick string) ([]int, error)
	TitleOf(ctx context.Context, table string, id int, fallback string) string
	SendNotice(ctx context.Context, userID int, title, content, relatedType string, relatedID, adminID int) error
}

type adminCommentRepo struct{ db *gorm.DB }

func NewAdminComment(db *gorm.DB) AdminCommentRepository {
	return &adminCommentRepo{db: db}
}

func (r *adminCommentRepo) ChildCounts(ctx context.Context, table string, ids []int) (map[int]int, error) {
	out := map[int]int{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []struct {
		PID   int `gorm:"column:pid"`
		Total int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table(table).
		Select("pid, COUNT(*) AS total").
		Where("pid IN ? AND deleted_at IS NULL", ids).
		Group("pid").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.PID] = row.Total
	}
	return out, nil
}

func (r *adminCommentRepo) ListUsersBrief(ctx context.Context, ids []int) (map[int]CommentUserBrief, error) {
	out := map[int]CommentUserBrief{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []CommentUserBrief
	err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, avatar").Where("id IN ?", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = row
	}
	return out, nil
}

func (r *adminCommentRepo) ListTitles(ctx context.Context, table string, ids []int) (map[int]CommentTitleBrief, error) {
	out := map[int]CommentTitleBrief{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []CommentTitleBrief
	err := r.db.WithContext(ctx).Table(table).Select("id, title").Where("id IN ?", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = row
	}
	return out, nil
}

func (r *adminCommentRepo) ListChapters(ctx context.Context, ids []int) (map[int]CommentTitleBrief, error) {
	out := map[int]CommentTitleBrief{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []CommentTitleBrief
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, title, course_id").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = row
	}
	return out, nil
}

func (r *adminCommentRepo) IDsByTitleLike(ctx context.Context, table, title string) ([]int, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, nil
	}
	var ids []int
	q := r.db.WithContext(ctx).Table(table).Where("title LIKE ?", "%"+escapeLike(title)+"%")
	switch table {
	case "course_chapters", "courses", "paths", "papers":
		q = q.Where("deleted_at IS NULL")
	}
	err := q.Limit(1000).Pluck("id", &ids).Error
	return ids, err
}

func (r *adminCommentRepo) IDsByNickLike(ctx context.Context, nick string) ([]int, error) {
	nick = strings.TrimSpace(nick)
	if nick == "" {
		return nil, nil
	}
	var ids []int
	err := r.db.WithContext(ctx).Table("users").
		Where("nick_name LIKE ?", "%"+escapeLike(nick)+"%").Limit(1000).Pluck("id", &ids).Error
	return ids, err
}

func (r *adminCommentRepo) TitleOf(ctx context.Context, table string, id int, fallback string) string {
	var title string
	_ = r.db.WithContext(ctx).Table(table).Select("title").Where("id = ?", id).Scan(&title).Error
	if title == "" {
		return fallback
	}
	return title
}

func (r *adminCommentRepo) SendNotice(ctx context.Context, userID int, title, content, relatedType string, relatedID, adminID int) error {
	if userID <= 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Table("user_messages").Create(map[string]any{
		"user_id": userID, "admin_id": adminID, "title": title, "content": content,
		"type": "notice", "related_type": relatedType, "related_id": relatedID,
		"is_read": 0, "created_at": now, "updated_at": now,
	}).Error
}

func escapeLike(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(s)
}
