package adminchapter

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
)

func (s *Service) GetChapterUserList(c *gin.Context) {
	var raw map[string]any
	_ = c.ShouldBindJSON(&raw)
	if raw == nil {
		raw = map[string]any{}
	}
	chapterID := asInt(raw["chapter_id"])
	if chapterID <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	current, size := page.Clamp(asInt(raw["page"]), asInt(raw["size"]))
	ctx := c.Request.Context()
	cityID, err := s.resolveCity(c, asInt(raw["city_id"]))
	if response.WriteErr(c, err) {
		return
	}

	rows, total, err := s.chapters.ListChapterUsers(ctx, repo.ChapterUserListFilter{
		ChapterID: chapterID, CityID: cityID,
		NickName: strings.TrimSpace(asString(raw["nick_name"])),
		Mobile:   strings.TrimSpace(asString(raw["mobile"])),
		Offset:   (current - 1) * size, Limit: size,
	})
	if response.WriteErr(c, err) {
		return
	}

	userIDs := make([]int, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}
	users := map[int]gin.H{}
	if len(userIDs) > 0 {
		userRows, err := s.chapters.ListUsersFull(ctx, userIDs)
		if response.WriteErr(c, err) {
			return
		}
		for _, u := range userRows {
			users[u.ID] = gin.H{
				"id": u.ID, "nick_name": u.NickName, "credit": u.Credit,
				"avatar": u.Avatar, "mobile": u.Mobile,
			}
		}
	}

	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var user any
		if found, ok := users[row.UserID]; ok {
			user = found
		}
		list = append(list, gin.H{
			"id": row.ID, "user_id": row.UserID, "chapter_id": row.ChapterID,
			"course_id": row.CourseID, "created_at": model.FormatTime(row.CreatedAt),
			"updated_at": model.FormatTime(row.UpdatedAt), "un_lock": row.UnLock, "is_end": row.IsEnd,
			"user": user,
		})
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) UnlockChapterUser(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	id := asInt(raw["id"])
	if id <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	row, err := s.chapters.FindUserChapter(ctx, id)
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "记录不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if row.UnLock == 1 {
		response.Fail(c, 1, "已经解锁过了")
		return
	}
	chapter, err := s.chapters.FindChapterLimitDay(ctx, row.ChapterID)
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "课节不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	limitDay := 30
	if chapter.LimitDay != nil && *chapter.LimitDay > 0 {
		limitDay = *chapter.LimitDay
	}
	expire := time.Now().AddDate(0, 0, limitDay)
	if response.WriteErr(c, s.chapters.UnlockUserChapter(ctx, id, expire)) {
		return
	}
	response.OK(c, "解锁成功", []any{})
}
