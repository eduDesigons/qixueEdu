package comment

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/async"
	"qixue-api/internal/platform/bg"
	webapprepo "qixue-api/internal/webapp/repo"
	"qixue-api/internal/service/coin"
)

var (
	ErrEmptyContent   = errors.New("评论内容不能为空")
	ErrNeedChapter    = errors.New("请选择章节")
	ErrNeedCourse     = errors.New("请选择课程")
	ErrNeedJoinCourse = errors.New("请先加入课程")
)

type Service struct {
	comments webapprepo.CommentRepository
	db       *gorm.DB
	rdb      *redis.Client
}

func New(comments webapprepo.CommentRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{comments: comments, db: db, rdb: rdb}
}

func clampPage(pageNum, limit int) (int, int) {
	if pageNum < 1 {
		pageNum = 1
	}
	if limit < 1 {
		limit = 10
	}
	return pageNum, limit
}

func (s *Service) ListCourse(ctx context.Context, courseID, pageNum, limit int) (page.MoreBody, error) {
	pageNum, limit = clampPage(pageNum, limit)
	var (
		total int64
		rows  []model.CourseComment
	)
	if err := async.Run(
		func() error {
			var callErr error
			total, callErr = s.comments.CountCourseRoot(ctx, courseID)
			return callErr
		},
		func() error {
			var callErr error
			rows, callErr = s.comments.ListCourseRoot(ctx, courseID, pageNum, limit)
			return callErr
		},
	); err != nil {
		return page.MoreBody{}, err
	}
	list, err := s.decorateCourse(ctx, rows)
	if err != nil {
		return page.MoreBody{}, err
	}
	return page.NewMore(pageNum, limit, total, list), nil
}

func (s *Service) ListChapter(ctx context.Context, chapterID, pageNum, limit int) (page.MoreBody, error) {
	pageNum, limit = clampPage(pageNum, limit)
	var (
		total int64
		rows  []model.CourseChapterComment
	)
	if err := async.Run(
		func() error {
			var callErr error
			total, callErr = s.comments.CountChapterRoot(ctx, chapterID)
			return callErr
		},
		func() error {
			var callErr error
			rows, callErr = s.comments.ListChapterRoot(ctx, chapterID, pageNum, limit)
			return callErr
		},
	); err != nil {
		return page.MoreBody{}, err
	}
	list, err := s.decorateChapter(ctx, rows)
	if err != nil {
		return page.MoreBody{}, err
	}
	return page.NewMore(pageNum, limit, total, list), nil
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
	row := &model.CourseComment{
		PID: pid, CourseID: courseID, UserID: uid, Content: content,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.comments.CreateCourseComment(ctx, row); err != nil {
		return nil, err
	}
	bg.Go(ctx, func(ctx context.Context) {
		coin.OnComment(ctx, s.db, s.rdb, uid, pid, content)
	})
	return courseView(*row, nil, map[int]map[string]any{}), nil
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
	row := &model.CourseChapterComment{
		PID: pid, CourseID: courseID, ChapterID: chapterID, UserID: uid, Content: content,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.comments.CreateChapterComment(ctx, row); err != nil {
		return nil, err
	}
	bg.Go(ctx, func(ctx context.Context) {
		coin.OnComment(ctx, s.db, s.rdb, uid, pid, content)
	})
	return chapterView(*row, nil, map[int]map[string]any{}), nil
}

func (s *Service) decorateCourse(ctx context.Context, rows []model.CourseComment) ([]map[string]any, error) {
	if len(rows) == 0 {
		return []map[string]any{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	for _, r := range rows {
		ids = append(ids, int(r.ID))
		userIDs[r.UserID] = struct{}{}
	}
	kids, err := s.comments.ListCourseChildren(ctx, ids)
	if err != nil {
		return nil, err
	}
	children := map[int][]model.CourseComment{}
	for _, k := range kids {
		children[k.PID] = append(children[k.PID], k)
		userIDs[k.UserID] = struct{}{}
	}
	users, err := s.loadUsers(ctx, keys(userIDs))
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		childViews := make([]map[string]any, 0, len(children[int(r.ID)]))
		for _, kid := range children[int(r.ID)] {
			childViews = append(childViews, courseView(kid, nil, users))
		}
		list = append(list, courseView(r, childViews, users))
	}
	return list, nil
}

func (s *Service) decorateChapter(ctx context.Context, rows []model.CourseChapterComment) ([]map[string]any, error) {
	if len(rows) == 0 {
		return []map[string]any{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	for _, r := range rows {
		ids = append(ids, int(r.ID))
		userIDs[r.UserID] = struct{}{}
	}
	kids, err := s.comments.ListChapterChildren(ctx, ids)
	if err != nil {
		return nil, err
	}
	children := map[int][]model.CourseChapterComment{}
	for _, k := range kids {
		children[k.PID] = append(children[k.PID], k)
		userIDs[k.UserID] = struct{}{}
	}
	users, err := s.loadUsers(ctx, keys(userIDs))
	if err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		childViews := make([]map[string]any, 0, len(children[int(r.ID)]))
		for _, kid := range children[int(r.ID)] {
			childViews = append(childViews, chapterView(kid, nil, users))
		}
		list = append(list, chapterView(r, childViews, users))
	}
	return list, nil
}

func courseView(r model.CourseComment, children []map[string]any, users map[int]map[string]any) map[string]any {
	if children == nil {
		children = []map[string]any{}
	}
	return map[string]any{
		"id": int(r.ID), "pid": r.PID, "course_id": r.CourseID, "user_id": r.UserID,
		"content":    r.Content,
		"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
		"user": userBrief(users, r.UserID), "children": children,
	}
}

func chapterView(r model.CourseChapterComment, children []map[string]any, users map[int]map[string]any) map[string]any {
	if children == nil {
		children = []map[string]any{}
	}
	return map[string]any{
		"id": int(r.ID), "pid": r.PID, "course_id": r.CourseID, "chapter_id": r.ChapterID, "user_id": r.UserID,
		"content":    r.Content,
		"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
		"user": userBrief(users, r.UserID), "children": children,
	}
}

func (s *Service) loadUsers(ctx context.Context, ids []int) (map[int]map[string]any, error) {
	out := map[int]map[string]any{}
	rows, err := s.comments.ListUserBriefs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r.ID <= 0 {
			continue
		}
		nick := strings.TrimSpace(r.NickName)
		if nick == "" {
			nick = "学员"
		}
		out[r.ID] = map[string]any{
			"id": r.ID, "nick_name": nick, "avatar": r.Avatar, "coins": r.Coins,
		}
	}
	return out, nil
}

func userBrief(users map[int]map[string]any, uid int) map[string]any {
	if uid == 0 {
		return map[string]any{"id": 0, "nick_name": "管理员", "avatar": "", "coins": 0}
	}
	if u, ok := users[uid]; ok {
		return u
	}
	return map[string]any{"id": uid, "nick_name": "学员", "avatar": "", "coins": 0}
}

func keys(m map[int]struct{}) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		if k > 0 {
			out = append(out, k)
		}
	}
	return out
}
