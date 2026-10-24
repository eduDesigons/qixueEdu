package comment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/coin"
)

var (
	ErrEmptyContent   = errors.New("评论内容不能为空")
	ErrNeedChapter    = errors.New("请选择章节")
	ErrNeedCourse     = errors.New("请选择课程")
	ErrNeedJoinCourse = errors.New("请先加入课程")
)

type Service struct {
	comments repo.CommentRepository
	db       *gorm.DB
	rdb      *redis.Client
}

func New(comments repo.CommentRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{comments: comments, db: db, rdb: rdb}
}

func (s *Service) ListChapter(ctx context.Context, chapterID, pageNum, limit int) (page.MoreBody, error) {
	return s.list(ctx, "course_chapter_comments", true, chapterID, pageNum, limit)
}

func (s *Service) ListCourse(ctx context.Context, courseID, pageNum, limit int) (page.MoreBody, error) {
	return s.list(ctx, "course_comments", false, courseID, pageNum, limit)
}

func (s *Service) list(ctx context.Context, table string, withChapter bool, targetID, pageNum, limit int) (page.MoreBody, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 10
	}
	total, err := s.comments.CountRoot(ctx, table, withChapter, targetID)
	if err != nil {
		return page.MoreBody{}, err
	}
	rows, err := s.comments.ListRoot(ctx, table, withChapter, targetID, pageNum, limit)
	if err != nil {
		return page.MoreBody{}, err
	}
	list, err := s.decorate(ctx, table, rows, withChapter)
	if err != nil {
		return page.MoreBody{}, err
	}
	return page.NewMore(pageNum, limit, total, list), nil
}

func (s *Service) ChapterComment(ctx context.Context, uid, chapterID, courseID, pid int, content string) (map[string]any, error) {
	if strings.TrimSpace(content) == "" {
		return nil, ErrEmptyContent
	}
	if chapterID == 0 {
		return nil, ErrNeedChapter
	}
	if courseID == 0 {
		return nil, ErrNeedCourse
	}
	n, err := s.comments.CountUserCourse(ctx, uid, courseID)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrNeedJoinCourse
	}
	now := time.Now()
	row := repo.CommentRow{
		PID: pid, CourseID: courseID, ChapterID: chapterID,
		UserID: uid, Content: content, CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.comments.CreateChapterComment(ctx, &row); err != nil {
		return nil, err
	}
	coin.OnComment(ctx, s.db, s.rdb, uid, pid, content)
	return oneView(row, true), nil
}

func (s *Service) CourseComment(ctx context.Context, uid, courseID, pid int, content string) (map[string]any, error) {
	if strings.TrimSpace(content) == "" {
		return nil, ErrEmptyContent
	}
	if courseID == 0 {
		return nil, ErrNeedCourse
	}
	n, err := s.comments.CountUserCourse(ctx, uid, courseID)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrNeedJoinCourse
	}
	now := time.Now()
	row := repo.CommentRow{
		PID: pid, CourseID: courseID, UserID: uid,
		Content: content, CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.comments.CreateCourseComment(ctx, &row); err != nil {
		return nil, err
	}
	coin.OnComment(ctx, s.db, s.rdb, uid, pid, content)
	return oneView(row, false), nil
}

func oneView(row repo.CommentRow, withChapter bool) map[string]any {
	item := map[string]any{
		"id": row.ID, "pid": row.PID, "course_id": row.CourseID, "user_id": row.UserID,
		"content":    row.Content,
		"created_at": model.FormatTime(row.CreatedAt), "updated_at": model.FormatTime(row.UpdatedAt),
		"children": []any{},
	}
	if withChapter {
		item["chapter_id"] = row.ChapterID
	}
	return item
}

func (s *Service) decorate(ctx context.Context, table string, rows []repo.CommentRow, withChapter bool) ([]map[string]any, error) {
	if len(rows) == 0 {
		return []map[string]any{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	for _, r := range rows {
		ids = append(ids, r.ID)
		userIDs[r.UserID] = struct{}{}
	}
	childRows, err := s.comments.ListChildren(ctx, table, ids)
	if err != nil {
		return nil, err
	}
	children := map[int][]repo.CommentRow{}
	for _, kid := range childRows {
		children[kid.PID] = append(children[kid.PID], kid)
		userIDs[kid.UserID] = struct{}{}
	}
	users, err := s.loadUsers(ctx, keys(userIDs))
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		kids := children[r.ID]
		childViews := make([]map[string]any, 0, len(kids))
		for _, kid := range kids {
			childViews = append(childViews, rowView(kid, users, withChapter, nil))
		}
		list = append(list, rowView(r, users, withChapter, childViews))
	}
	return list, nil
}

func rowView(r repo.CommentRow, users map[int]map[string]any, withChapter bool, children []map[string]any) map[string]any {
	if children == nil {
		children = []map[string]any{}
	}
	item := map[string]any{
		"id": r.ID, "pid": r.PID, "course_id": r.CourseID, "user_id": r.UserID,
		"content":    r.Content,
		"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
		"user":     userBrief(users, r.UserID),
		"children": children,
	}
	if withChapter {
		item["chapter_id"] = r.ChapterID
	}
	return item
}

func (s *Service) loadUsers(ctx context.Context, ids []int) (map[int]map[string]any, error) {
	out := map[int]map[string]any{}
	rows, err := s.comments.ListUserBriefs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		id := bind.AsInt(r["id"])
		if id <= 0 {
			continue
		}
		out[id] = normalizeUserBrief(r, id)
	}
	return out, nil
}

func normalizeUserBrief(r map[string]any, id int) map[string]any {
	nick := ""
	switch v := r["nick_name"].(type) {
	case string:
		nick = v
	case []byte:
		nick = string(v)
	case nil:
		nick = ""
	default:
		if v != nil {
			nick = fmt.Sprint(v)
		}
	}
	avatar := ""
	switch v := r["avatar"].(type) {
	case string:
		avatar = v
	case []byte:
		avatar = string(v)
	}
	return map[string]any{
		"id":        id,
		"nick_name": nick,
		"avatar":    avatar,
		"coins":     bind.AsInt(r["coins"]),
	}
}

func userBrief(users map[int]map[string]any, uid int) map[string]any {
	if uid == 0 {
		return map[string]any{"id": 0, "nick_name": "管理员", "avatar": "", "coins": 0}
	}
	if u, ok := users[uid]; ok {
		return u
	}
	return map[string]any{"id": uid, "nick_name": "", "avatar": "", "coins": 0}
}

func keys(m map[int]struct{}) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
