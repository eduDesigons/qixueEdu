package adminchapter

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

func (s *Service) GetChapterUserList(ctx context.Context, adminID int, raw map[string]any) (page.Body, error) {
	if raw == nil {
		raw = map[string]any{}
	}
	chapterID := asInt(raw["chapter_id"])
	if chapterID <= 0 {
		return page.Body{}, biz("参数错误")
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	cityID, err := s.resolveCity(ctx, adminID, asInt(raw["city_id"]))
	if err != nil {
		return page.Body{}, err
	}

	rows, total, err := s.chapters.ListChapterUsers(ctx, repo.ChapterUserListFilter{
		ChapterID: chapterID, CityID: cityID,
		NickName: strings.TrimSpace(asString(raw["nick_name"])),
		Mobile:   strings.TrimSpace(asString(raw["mobile"])),
		Offset:   (current - 1) * size, Limit: size,
	})
	if err != nil {
		return page.Body{}, err
	}

	userIDs := make([]int, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}
	users := map[int]map[string]any{}
	if len(userIDs) > 0 {
		userRows, err := s.chapters.ListUsersFull(ctx, userIDs)
		if err != nil {
			return page.Body{}, err
		}
		for _, u := range userRows {
			users[u.ID] = map[string]any{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		var user any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		list = append(list, map[string]any{
			"id": row.ID, "user_id": row.UserID, "chapter_id": row.ChapterID,
			"course_id": row.CourseID, "created_at": model.FormatTime(row.CreatedAt),
			"updated_at": model.FormatTime(row.UpdatedAt), "un_lock": row.UnLock, "is_end": row.IsEnd,
			"user": user,
		})
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) UnlockChapterUser(ctx context.Context, raw map[string]any) error {
	if raw == nil {
		return biz("参数错误")
	}
	id := asInt(raw["id"])
	if id <= 0 {
		return biz("参数错误")
	}
	row, err := s.chapters.FindUserChapter(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return biz("记录不存在")
	}
	if err != nil {
		return err
	}
	if row.UnLock == 1 {
		return biz("已经解锁过了")
	}
	chapter, err := s.chapters.FindChapterLimitDay(ctx, row.ChapterID)
	if err == gorm.ErrRecordNotFound {
		return biz("课节不存在")
	}
	if err != nil {
		return err
	}
	limitDay := 30
	if chapter.LimitDay != nil && *chapter.LimitDay > 0 {
		limitDay = *chapter.LimitDay
	}
	expire := time.Now().AddDate(0, 0, limitDay)
	return s.chapters.UnlockUserChapter(ctx, id, expire)
}
