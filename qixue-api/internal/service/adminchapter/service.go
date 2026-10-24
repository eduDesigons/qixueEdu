package adminchapter

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/model"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/enroll"
)

type Service struct {
	chapters repo.AdminChapterRepository
	db       *gorm.DB
	rdb      *redis.Client
}

func New(chapters repo.AdminChapterRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{chapters: chapters, db: db, rdb: rdb}
}

type chapterRow struct {
	ID        int        `gorm:"column:id;primaryKey"`
	Title     string     `gorm:"column:title"`
	CourseID  int        `gorm:"column:course_id"`
	Type      int        `gorm:"column:type"`
	Content   *string    `gorm:"column:content"`
	LimitDay  *int       `gorm:"column:limit_day"`
	MediaID   *int       `gorm:"column:media_id"`
	MediaURL  *string    `gorm:"column:media_url"`
	PaperID   *int       `gorm:"column:paper_id"`
	Sort      int        `gorm:"column:sort"`
	CreatedID int        `gorm:"column:created_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (chapterRow) TableName() string { return "course_chapters" }

func (s *Service) GetChapterList(c *gin.Context) {
	courseID := bind.Query(c, "courseId", "course_id")
	title := strings.TrimSpace(c.Query("title"))
	if courseID <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	query := s.db.WithContext(ctx).Model(&chapterRow{}).
		Select("id, title, paper_id, limit_day, sort").
		Where("course_id = ? AND deleted_at IS NULL", courseID)
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	var rows []chapterRow
	if response.WriteErr(c, query.Order("sort ASC").Find(&rows).Error) {
		return
	}
	ids := make([]int, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	comments, users, ends := map[int]int{}, map[int]int{}, map[int]int{}
	if len(ids) > 0 {
		commentRows, err := s.chapters.CommentCountsByChapterIDs(ctx, ids)
		if response.WriteErr(c, err) {
			return
		}
		for _, row := range commentRows {
			comments[row.Key] = row.Total
		}
		userRows, err := s.chapters.UserCountsByChapterIDs(ctx, ids)
		if response.WriteErr(c, err) {
			return
		}
		for _, row := range userRows {
			users[row.Key] = row.Total
		}
		endRows, err := s.chapters.EndedUserCountsByChapterIDs(ctx, ids)
		if response.WriteErr(c, err) {
			return
		}
		for _, row := range endRows {
			ends[row.Key] = row.Total
		}
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, gin.H{
			"id": row.ID, "title": row.Title, "paper_id": row.PaperID,
			"limit_day": row.LimitDay, "sort": row.Sort,
			"comments_count": comments[row.ID], "users_count": users[row.ID],
			"user_chapter_end_count": ends[row.ID],
		})
	}
	response.OK(c, "", list)
}

func (s *Service) GetChapterDetail(c *gin.Context) {
	id := bind.Query(c, "id")
	if id <= 0 {
		response.Fail(c, 1, "课节不存在")
		return
	}
	ctx := c.Request.Context()
	var row chapterRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "课节不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}

	var creator any
	if row.CreatedID > 0 {
		if admin, err := s.chapters.FindAdmin(ctx, row.CreatedID); err == nil {
			creator = gin.H{"id": admin.ID, "nickName": admin.NickName}
		}
	}
	var paper any
	if row.PaperID != nil && *row.PaperID > 0 {
		if p, err := s.chapters.FindPaperScore(ctx, *row.PaperID); err == nil {
			paper = gin.H{"id": p.ID, "score": p.Score, "pass_score": p.PassScore}
		}
	}
	mediaThumb := ""
	var media any
	if row.MediaID != nil && *row.MediaID > 0 {
		if m, err := s.chapters.FindMedia(ctx, *row.MediaID); err == nil {
			media = gin.H{
				"id": m.ID, "title": m.Title, "name": m.Title,
				"thumb": m.Thumb, "url": m.URL, "type": m.Type, "size": m.Size,
			}
			mediaThumb = m.Thumb
		}
	}

	commentsCount, _ := s.chapters.CountComments(ctx, id)
	usersCount, _ := s.chapters.CountUsers(ctx, id)
	endCount, _ := s.chapters.CountEndedUsers(ctx, id)

	response.OK(c, "", gin.H{
		"id": row.ID, "title": row.Title, "course_id": row.CourseID, "type": row.Type,
		"content": row.Content, "limit_day": row.LimitDay, "media_id": row.MediaID,
		"media_url": row.MediaURL, "created_at": model.FormatTime(row.CreatedAt),
		"updated_at": model.FormatTime(row.UpdatedAt), "created_id": row.CreatedID,
		"paper_id": row.PaperID, "creator": creator, "paper": paper, "media": media,
		"media_thumb": mediaThumb, "comments_count": commentsCount,
		"users_count": usersCount, "user_chapter_end_count": endCount,
	})
}

func (s *Service) ChapterSave(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	title := strings.TrimSpace(asString(raw["title"]))
	courseID := asInt(raw["course_id"])
	if title == "" || courseID <= 0 {
		response.Fail(c, 1, "请填写课节标题")
		return
	}
	ctx := c.Request.Context()
	courseCount, err := s.chapters.CountCourse(ctx, courseID)
	if response.WriteErr(c, err) {
		return
	}
	if courseCount == 0 {
		response.Fail(c, 1, "课程不存在")
		return
	}

	now := time.Now()
	id := asInt(raw["id"])
	data := map[string]any{
		"title": title, "course_id": courseID, "type": asInt(raw["type"]),
		"content": asString(raw["content"]), "limit_day": asInt(raw["limit_day"]),
		"updated_at": now,
	}
	if v, ok := raw["media_id"]; ok {
		if asInt(v) > 0 {
			data["media_id"] = asInt(v)
		} else {
			data["media_id"] = nil
		}
	}
	if v, ok := raw["media_url"]; ok {
		data["media_url"] = asString(v)
	}

	if id > 0 {
		var row chapterRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "课节不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		cols := []string{"title", "course_id", "type", "content", "limit_day", "updated_at"}
		if _, ok := raw["media_id"]; ok {
			cols = append(cols, "media_id")
		}
		if _, ok := raw["media_url"]; ok {
			cols = append(cols, "media_url")
		}
		err = s.db.WithContext(ctx).Model(&chapterRow{}).Where("id = ?", id).
			Select(cols).Updates(data).Error
		if response.WriteErr(c, err) {
			return
		}
		response.OK(c, "更新成功", []any{})
		return
	}

	var sort int64
	if response.WriteErr(c, s.db.WithContext(ctx).Model(&chapterRow{}).
		Where("course_id = ?", courseID).Count(&sort).Error) {
		return
	}
	limitDay := asInt(raw["limit_day"])
	row := chapterRow{
		Title: title, CourseID: courseID, Type: asInt(raw["type"]),
		LimitDay: &limitDay, Sort: int(sort) + 1,
		CreatedID: middleware.UID(c), CreatedAt: &now, UpdatedAt: &now,
	}
	if content := asString(raw["content"]); content != "" {
		row.Content = &content
	}
	if mid := asInt(raw["media_id"]); mid > 0 {
		row.MediaID = &mid
	}
	if url := asString(raw["media_url"]); url != "" {
		row.MediaURL = &url
	}
	if response.WriteErr(c, s.chapters.CreateChapterAndInc(ctx, &row, courseID)) {
		return
	}
	if response.WriteErr(c, enroll.SyncEnrolledChapters(ctx, s.db, s.rdb, courseID)) {
		return
	}
	response.OK(c, "课节创建成功", []any{})
}

func (s *Service) ChapterDelete(c *gin.Context) {
	id := bind.Query(c, "id")
	if id <= 0 {
		response.Fail(c, 1, "课节不存在")
		return
	}
	ctx := c.Request.Context()
	var row chapterRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "课节不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	userCount, err := s.chapters.CountUsers(ctx, id)
	if response.WriteErr(c, err) {
		return
	}
	if userCount > 0 {
		response.Fail(c, 1, "学员已报名，无法删除")
		return
	}

	if response.WriteErr(c, s.chapters.DeleteChapterCascade(ctx, id, row.CourseID, row.PaperID)) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (s *Service) ChapterSort(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	direction := asString(raw["direction"])
	if direction != "up" && direction != "down" {
		response.Fail(c, 1, "参数错误")
		return
	}
	id := asInt(raw["id"])
	ctx := c.Request.Context()
	var current chapterRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&current).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "课节不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}

	var other chapterRow
	q := s.db.WithContext(ctx).Where("course_id = ?", current.CourseID)
	if direction == "up" {
		err = q.Where("sort < ?", current.Sort).Order("sort DESC").Take(&other).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "已经是第一节")
			return
		}
	} else {
		err = q.Where("sort > ?", current.Sort).Order("sort ASC").Take(&other).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "已经是最后一节")
			return
		}
	}
	if response.WriteErr(c, err) {
		return
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		temp := current.Sort
		if err := tx.Model(&chapterRow{}).Where("id = ?", current.ID).
			Update("sort", other.Sort).Error; err != nil {
			return err
		}
		return tx.Model(&chapterRow{}).Where("id = ?", other.ID).
			Update("sort", temp).Error
	})
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "成功", []any{})
}
