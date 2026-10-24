package user

import (
	"context"
	"time"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

func (s *Service) decorateComments(ctx context.Context, table string, rows []repo.CommentRow, withChapter bool) ([]map[string]any, error) {
	if len(rows) == 0 {
		return []map[string]any{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	courseIDs := map[int]struct{}{}
	chapterIDs := map[int]struct{}{}
	for _, r := range rows {
		ids = append(ids, r.ID)
		userIDs[r.UserID] = struct{}{}
		courseIDs[r.CourseID] = struct{}{}
		if withChapter && r.ChapterID > 0 {
			chapterIDs[r.ChapterID] = struct{}{}
		}
	}
	childRows, err := s.learner.ListCommentChildren(ctx, table, ids)
	if err != nil {
		return nil, err
	}
	children := map[int][]map[string]any{}
	formatTimes(childRows)
	for _, r := range childRows {
		pid := asInt(r["pid"])
		children[pid] = append(children[pid], r)
		userIDs[asInt(r["user_id"])] = struct{}{}
	}
	userRows, err := s.learner.ListUserBriefs(ctx, keys(userIDs))
	if err != nil {
		return nil, err
	}
	users := map[int]map[string]any{}
	for _, r := range userRows {
		id := asInt(r["id"])
		if id <= 0 {
			continue
		}
		nick := ""
		switch v := r["nick_name"].(type) {
		case string:
			nick = v
		case []byte:
			nick = string(v)
		}
		avatar := ""
		switch v := r["avatar"].(type) {
		case string:
			avatar = v
		case []byte:
			avatar = string(v)
		}
		users[id] = map[string]any{
			"id": id, "nick_name": nick, "avatar": avatar, "coins": asInt(r["coins"]),
		}
	}
	courseRows, err := s.learner.ListCourseBriefs(ctx, keys(courseIDs))
	if err != nil {
		return nil, err
	}
	courses := map[int]map[string]any{}
	for _, r := range courseRows {
		courses[r.ID] = map[string]any{"id": r.ID, "title": r.Title, "thumb": r.Thumb, "desc": deref(r.Desc)}
	}
	var chapters map[int]map[string]any
	if withChapter {
		chapterRows, err := s.learner.ListChapterBriefs(ctx, keys(chapterIDs))
		if err != nil {
			return nil, err
		}
		chapters = map[int]map[string]any{}
		for _, r := range chapterRows {
			chapters[r.ID] = map[string]any{"id": r.ID, "title": r.Title}
		}
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		item := map[string]any{
			"id": r.ID, "pid": r.PID, "course_id": r.CourseID, "user_id": r.UserID,
			"content":    r.Content,
			"created_at": model.FormatTime(r.CreatedAt), "updated_at": model.FormatTime(r.UpdatedAt),
			"user":     userBrief(users, r.UserID),
			"course":   courses[r.CourseID],
			"children": attachUsers(children[r.ID], users),
		}
		if withChapter {
			item["chapter_id"] = r.ChapterID
			item["chapter"] = chapters[r.ChapterID]
		}
		list = append(list, item)
	}
	return list, nil
}

func (s *Service) loadPathCerts(ctx context.Context, uid int, pathIDs []int) (map[int]any, error) {
	out := map[int]any{}
	rows, err := s.learner.ListValidPathCerts(ctx, uid, pathIDs)
	if err != nil {
		return nil, err
	}
	formatTimes(rows)
	for _, r := range rows {
		out[asInt(r["path_id"])] = r
	}
	return out, nil
}

func attachUsers(rows []map[string]any, users map[int]map[string]any) []map[string]any {
	if rows == nil {
		return []map[string]any{}
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		cp := map[string]any{}
		for k, v := range r {
			cp[k] = v
		}
		cp["user"] = userBrief(users, asInt(r["user_id"]))
		out = append(out, cp)
	}
	return out
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

func asInt(v any) int { return bind.AsInt(v) }

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
			case string:
				// already formatted
			}
		}
	}
}
