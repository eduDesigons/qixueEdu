package user

import (
	"context"
	"time"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/async"
)

func (s *Service) decorateMyCourseComments(ctx context.Context, rows []model.CourseComment) ([]map[string]any, error) {
	if len(rows) == 0 {
		return []map[string]any{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	courseIDs := map[int]struct{}{}
	for _, r := range rows {
		ids = append(ids, int(r.ID))
		userIDs[r.UserID] = struct{}{}
		courseIDs[r.CourseID] = struct{}{}
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
	var (
		users   map[int]map[string]any
		courses map[int]map[string]any
	)
	if err := async.Run(
		func() error {
			var callErr error
			users, callErr = s.loadCommentUsers(ctx, keys(userIDs))
			return callErr
		},
		func() error {
			var callErr error
			courses, callErr = s.loadCommentCourses(ctx, keys(courseIDs))
			return callErr
		},
	); err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		childViews := make([]map[string]any, 0, len(children[int(r.ID)]))
		for _, kid := range children[int(r.ID)] {
			childViews = append(childViews, map[string]any{
				"id": int(kid.ID), "pid": kid.PID, "course_id": kid.CourseID, "user_id": kid.UserID,
				"content": kid.Content,
				"created_at": model.FormatTime(kid.CreatedAt), "updated_at": model.FormatTime(kid.UpdatedAt),
				"user": userBrief(users, kid.UserID),
			})
		}
		list = append(list, map[string]any{
			"id": int(r.ID), "pid": r.PID, "course_id": r.CourseID, "user_id": r.UserID,
			"content": r.Content,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"user": userBrief(users, r.UserID), "course": courses[r.CourseID], "children": childViews,
		})
	}
	return list, nil
}

func (s *Service) decorateMyChapterComments(ctx context.Context, rows []model.CourseChapterComment) ([]map[string]any, error) {
	if len(rows) == 0 {
		return []map[string]any{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	courseIDs := map[int]struct{}{}
	chapterIDs := map[int]struct{}{}
	for _, r := range rows {
		ids = append(ids, int(r.ID))
		userIDs[r.UserID] = struct{}{}
		courseIDs[r.CourseID] = struct{}{}
		if r.ChapterID > 0 {
			chapterIDs[r.ChapterID] = struct{}{}
		}
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
	var (
		users    map[int]map[string]any
		courses  map[int]map[string]any
		chapters map[int]map[string]any
	)
	if err := async.Run(
		func() error {
			var callErr error
			users, callErr = s.loadCommentUsers(ctx, keys(userIDs))
			return callErr
		},
		func() error {
			var callErr error
			courses, callErr = s.loadCommentCourses(ctx, keys(courseIDs))
			return callErr
		},
		func() error {
			var callErr error
			chapters, callErr = s.loadCommentChapters(ctx, keys(chapterIDs))
			return callErr
		},
	); err != nil {
		return nil, err
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		childViews := make([]map[string]any, 0, len(children[int(r.ID)]))
		for _, kid := range children[int(r.ID)] {
			childViews = append(childViews, map[string]any{
				"id": int(kid.ID), "pid": kid.PID, "course_id": kid.CourseID, "chapter_id": kid.ChapterID, "user_id": kid.UserID,
				"content": kid.Content,
				"created_at": model.FormatTime(kid.CreatedAt), "updated_at": model.FormatTime(kid.UpdatedAt),
				"user": userBrief(users, kid.UserID),
			})
		}
		list = append(list, map[string]any{
			"id": int(r.ID), "pid": r.PID, "course_id": r.CourseID, "chapter_id": r.ChapterID, "user_id": r.UserID,
			"content": r.Content,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"user": userBrief(users, r.UserID), "course": courses[r.CourseID], "chapter": chapters[r.ChapterID],
			"children": childViews,
		})
	}
	return list, nil
}

func (s *Service) loadCommentUsers(ctx context.Context, ids []int) (map[int]map[string]any, error) {
	out := map[int]map[string]any{}
	rows, err := s.comments.ListUserBriefs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r.ID <= 0 {
			continue
		}
		nick := r.NickName
		if nick == "" {
			nick = "学员"
		}
		out[r.ID] = map[string]any{
			"id": r.ID, "nick_name": nick, "avatar": r.Avatar, "coins": r.Coins,
		}
	}
	return out, nil
}

func (s *Service) loadCommentCourses(ctx context.Context, ids []int) (map[int]map[string]any, error) {
	out := map[int]map[string]any{}
	rows, err := s.comments.ListCourseBriefs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = map[string]any{"id": r.ID, "title": r.Title, "thumb": r.Thumb, "desc": deref(r.Desc)}
	}
	return out, nil
}

func (s *Service) loadCommentChapters(ctx context.Context, ids []int) (map[int]map[string]any, error) {
	out := map[int]map[string]any{}
	rows, err := s.comments.ListChapterBriefs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.ID] = map[string]any{"id": r.ID, "title": r.Title}
	}
	return out, nil
}

func (s *Service) loadPathCerts(ctx context.Context, uid int, pathIDs []int) (map[int]any, error) {
	out := map[int]any{}
	rows, err := s.learner.ListValidPathCerts(ctx, uid, pathIDs)
	if err != nil {
		return nil, err
	}
	formatTimes(rows)
	for _, r := range rows {
		out[bind.AsInt(r["path_id"])] = r
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

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(p *int, fallback int) int {
	if p != nil {
		return *p
	}
	return fallback
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

func formatTimes(rows []map[string]any) {
	for _, row := range rows {
		for _, key := range []string{"created_at", "updated_at", "read_at", "issued_at", "replied_at"} {
			switch v := row[key].(type) {
			case time.Time:
				row[key] = v.Format("2006-01-02 15:04:05")
			case *time.Time:
				row[key] = model.FormatTime(v)
			case []byte:
				if len(v) > 0 {
					row[key] = string(v)
				}
			}
		}
	}
}
