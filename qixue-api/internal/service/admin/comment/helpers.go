package admincomment

import (
	"context"
	"html"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"

	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

func (s *Service) decorateCourse(ctx context.Context, rows []commentRow) ([]CourseCommentView, error) {
	if len(rows) == 0 {
		return []CourseCommentView{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	courseIDs := map[int]struct{}{}
	for _, row := range rows {
		ids = append(ids, row.ID)
		if row.UserID > 0 {
			userIDs[row.UserID] = struct{}{}
		}
		courseIDs[row.CourseID] = struct{}{}
	}
	counts, err := s.comments.ChildCounts(ctx, "course_comments", ids)
	if err != nil {
		return nil, err
	}
	users, err := s.comments.ListUsersBrief(ctx, keys(userIDs))
	if err != nil {
		return nil, err
	}
	courses, err := s.comments.ListTitles(ctx, "courses", keys(courseIDs))
	if err != nil {
		return nil, err
	}
	list := make([]CourseCommentView, 0, len(rows))
	for _, row := range rows {
		item := CourseCommentView{
			ID: row.ID, PID: row.PID, CourseID: row.CourseID, UserID: row.UserID,
			Content: row.Content, CreatedAt: model.FormatTime(row.CreatedAt), ChildrenCount: counts[row.ID],
		}
		if u, ok := users[row.UserID]; ok {
			copy := toUserBrief(u)
			item.User = &copy
		}
		if t, ok := courses[row.CourseID]; ok {
			copy := toTitleBrief(t)
			item.Course = &copy
		}
		list = append(list, item)
	}
	return list, nil
}

func (s *Service) decorateChapter(ctx context.Context, rows []chapterCommentRow) ([]ChapterCommentView, error) {
	if len(rows) == 0 {
		return []ChapterCommentView{}, nil
	}
	ids := make([]int, 0, len(rows))
	userIDs := map[int]struct{}{}
	courseIDs := map[int]struct{}{}
	chapterIDs := map[int]struct{}{}
	for _, row := range rows {
		ids = append(ids, row.ID)
		if row.UserID > 0 {
			userIDs[row.UserID] = struct{}{}
		}
		courseIDs[row.CourseID] = struct{}{}
		chapterIDs[row.ChapterID] = struct{}{}
	}
	counts, err := s.comments.ChildCounts(ctx, "course_chapter_comments", ids)
	if err != nil {
		return nil, err
	}
	users, err := s.comments.ListUsersBrief(ctx, keys(userIDs))
	if err != nil {
		return nil, err
	}
	courses, err := s.comments.ListTitles(ctx, "courses", keys(courseIDs))
	if err != nil {
		return nil, err
	}
	chapters, err := s.comments.ListChapters(ctx, keys(chapterIDs))
	if err != nil {
		return nil, err
	}
	list := make([]ChapterCommentView, 0, len(rows))
	for _, row := range rows {
		item := ChapterCommentView{
			ID: row.ID, PID: row.PID, CourseID: row.CourseID, ChapterID: row.ChapterID, UserID: row.UserID,
			Content: row.Content, CreatedAt: model.FormatTime(row.CreatedAt), ChildrenCount: counts[row.ID],
		}
		if u, ok := users[row.UserID]; ok {
			copy := toUserBrief(u)
			item.User = &copy
		}
		if t, ok := courses[row.CourseID]; ok {
			copy := toTitleBrief(t)
			item.Course = &copy
		}
		if t, ok := chapters[row.ChapterID]; ok {
			copy := toTitleBrief(t)
			item.Chapter = &copy
		}
		list = append(list, item)
	}
	return list, nil
}

func toUserBrief(u repo.CommentUserBrief) UserBrief {
	return UserBrief{ID: u.ID, NickName: u.NickName, Avatar: u.Avatar}
}

func toTitleBrief(t repo.CommentTitleBrief) TitleBrief {
	return TitleBrief{ID: t.ID, Title: t.Title, CourseID: t.CourseID}
}

func (s *Service) applyTitleIDs(ctx context.Context, query *gorm.DB, column, table, title string) (*gorm.DB, bool, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return query, true, nil
	}
	ids, err := s.comments.IDsByTitleLike(ctx, table, title)
	if err != nil {
		return query, false, err
	}
	if len(ids) == 0 {
		return query, false, nil
	}
	return query.Where(column+" IN ?", ids), true, nil
}

func (s *Service) applyNickIDs(ctx context.Context, query *gorm.DB, nick string) (*gorm.DB, bool, error) {
	nick = strings.TrimSpace(nick)
	if nick == "" {
		return query, true, nil
	}
	ids, err := s.comments.IDsByNickLike(ctx, nick)
	if err != nil {
		return query, false, err
	}
	if len(ids) == 0 {
		return query, false, nil
	}
	return query.Where("user_id IN ?", ids), true, nil
}

func applyContent(query *gorm.DB, table, content string) *gorm.DB {
	content = strings.TrimSpace(content)
	if content == "" {
		return query
	}
	like := "%" + escapeLike(content) + "%"
	return query.Where(
		"(content LIKE ? OR EXISTS (SELECT 1 FROM "+table+" AS reply WHERE reply.pid = "+table+".id AND reply.pid > 0 AND reply.deleted_at IS NULL AND reply.content LIKE ?))",
		like, like,
	)
}

func applyTime(query *gorm.DB, start, end string) *gorm.DB {
	if strings.TrimSpace(start) != "" {
		query = query.Where("created_at >= ?", start)
	}
	if strings.TrimSpace(end) != "" {
		query = query.Where("created_at <= ?", end)
	}
	return query
}

func escapeLike(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(s)
}

func replyHTML(scope, origin, reply string) string {
	return "<p>你在" + html.EscapeString(scope) + "的评论收到了回复。</p>" +
		"<p>原评论：" + html.EscapeString(clipText(origin)) + "</p>" +
		"<p>回复：" + strings.ReplaceAll(html.EscapeString(reply), "\n", "<br>") + "</p>"
}

func clipText(text string) string {
	text = strings.TrimSpace(strings.Join(strings.Fields(text), " "))
	if utf8.RuneCountInString(text) > 80 {
		return string([]rune(text)[:80]) + "…"
	}
	return text
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

func keys(set map[int]struct{}) []int {
	out := make([]int, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}
