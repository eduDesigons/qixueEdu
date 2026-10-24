package admincomment

import (
	"strings"
	"time"

	"qixue-api/internal/http/bind"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/page"
	"qixue-api/internal/http/response"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/scope"
)

type Service struct {
	comments repo.AdminCommentRepository
	db       *gorm.DB
	rdb      *redis.Client
}

func New(comments repo.AdminCommentRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{comments: comments, db: db, rdb: rdb}
}

type commentRow struct {
	ID        int            `gorm:"column:id;primaryKey"`
	PID       int            `gorm:"column:pid"`
	UserID    int            `gorm:"column:user_id"`
	CourseID  int            `gorm:"column:course_id"`
	ChapterID int            `gorm:"column:chapter_id"`
	Content   string         `gorm:"column:content"`
	CreatedAt *time.Time     `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (commentRow) TableName() string { return "course_comments" }

type chapterCommentRow struct {
	ID        int            `gorm:"column:id;primaryKey"`
	PID       int            `gorm:"column:pid"`
	UserID    int            `gorm:"column:user_id"`
	CourseID  int            `gorm:"column:course_id"`
	ChapterID int            `gorm:"column:chapter_id"`
	Content   string         `gorm:"column:content"`
	CreatedAt *time.Time     `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (chapterCommentRow) TableName() string { return "course_chapter_comments" }

type userBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
}

type titleBrief struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	CourseID int    `json:"course_id,omitempty"`
}

type courseCommentView struct {
	ID            int         `json:"id"`
	PID           int         `json:"pid"`
	CourseID      int         `json:"course_id"`
	UserID        int         `json:"user_id"`
	Content       string      `json:"content"`
	CreatedAt     *string     `json:"created_at"`
	ChildrenCount int         `json:"children_count"`
	User          *userBrief  `json:"user"`
	Course        *titleBrief `json:"course"`
}

type chapterCommentView struct {
	ID            int         `json:"id"`
	PID           int         `json:"pid"`
	CourseID      int         `json:"course_id"`
	ChapterID     int         `json:"chapter_id"`
	UserID        int         `json:"user_id"`
	Content       string      `json:"content"`
	CreatedAt     *string     `json:"created_at"`
	ChildrenCount int         `json:"children_count"`
	User          *userBrief  `json:"user"`
	Course        *titleBrief `json:"course"`
	Chapter       *titleBrief `json:"chapter"`
}

func (s *Service) GetCourseCommentList(c *gin.Context) {
	s.listCourse(c)
}

func (s *Service) GetCourseChapterCommentList(c *gin.Context) {
	s.listChapter(c)
}

func (s *Service) listCourse(c *gin.Context) {
	var body struct {
		CourseID    any      `json:"course_id"`
		CourseTitle string   `json:"course_title"`
		NickName    string   `json:"nick_name"`
		Content     string   `json:"content"`
		StartTime   string   `json:"start_time"`
		EndTime     string   `json:"end_time"`
		Page        bind.Int `json:"page"`
		Size        bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	query := s.db.WithContext(ctx).Model(&commentRow{}).Where("pid = 0")
	if id := asInt(body.CourseID); id > 0 {
		query = query.Where("course_id = ?", id)
	}
	query, ok, err := s.applyTitleIDs(ctx, query, "course_id", "courses", body.CourseTitle)
	if response.WriteErr(c, err) {
		return
	}
	if !ok {
		response.OK(c, "", page.New(1, size, 0, []any{}))
		return
	}
	query, ok, err = s.applyNickIDs(ctx, query, body.NickName)
	if response.WriteErr(c, err) {
		return
	}
	if !ok {
		response.OK(c, "", page.New(1, size, 0, []any{}))
		return
	}
	query = applyContent(query, "course_comments", body.Content)
	query = applyTime(query, body.StartTime, body.EndTime)
	if limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []commentRow
	if response.WriteErr(c, query.Order("id DESC").Offset((current-1)*size).Limit(size).Find(&rows).Error) {
		return
	}
	list, err := s.decorateCourse(c, rows)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) listChapter(c *gin.Context) {
	var body struct {
		CourseID     any      `json:"course_id"`
		ChapterID    any      `json:"chapter_id"`
		CourseTitle  string   `json:"course_title"`
		ChapterTitle string   `json:"chapter_title"`
		NickName     string   `json:"nick_name"`
		Content      string   `json:"content"`
		StartTime    string   `json:"start_time"`
		EndTime      string   `json:"end_time"`
		Page         bind.Int `json:"page"`
		Size         bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	current, size := page.Clamp(body.Page.Int(), body.Size.Int())
	ctx := c.Request.Context()
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, middleware.UID(c))
	if response.WriteErr(c, err) {
		return
	}
	query := s.db.WithContext(ctx).Model(&chapterCommentRow{}).Where("pid = 0")
	if id := asInt(body.CourseID); id > 0 {
		query = query.Where("course_id = ?", id)
	}
	if id := asInt(body.ChapterID); id > 0 {
		query = query.Where("chapter_id = ?", id)
	}
	query, ok, err := s.applyTitleIDs(ctx, query, "course_id", "courses", body.CourseTitle)
	if response.WriteErr(c, err) {
		return
	}
	if !ok {
		response.OK(c, "", page.New(1, size, 0, []any{}))
		return
	}
	query, ok, err = s.applyTitleIDs(ctx, query, "chapter_id", "course_chapters", body.ChapterTitle)
	if response.WriteErr(c, err) {
		return
	}
	if !ok {
		response.OK(c, "", page.New(1, size, 0, []any{}))
		return
	}
	query, ok, err = s.applyNickIDs(ctx, query, body.NickName)
	if response.WriteErr(c, err) {
		return
	}
	if !ok {
		response.OK(c, "", page.New(1, size, 0, []any{}))
		return
	}
	query = applyContent(query, "course_chapter_comments", body.Content)
	query = applyTime(query, body.StartTime, body.EndTime)
	if limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var total int64
	if response.WriteErr(c, query.Count(&total).Error) {
		return
	}
	var rows []chapterCommentRow
	if response.WriteErr(c, query.Order("id DESC").Offset((current-1)*size).Limit(size).Find(&rows).Error) {
		return
	}
	list, err := s.decorateChapter(c, rows)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", page.New(current, size, total, list))
}

func (s *Service) AddCourseCommentReply(c *gin.Context) {
	s.reply(c, false)
}

func (s *Service) AddChapterCommentReply(c *gin.Context) {
	s.reply(c, true)
}

func (s *Service) reply(c *gin.Context, chapter bool) {
	var body struct {
		CommentID bind.Int `json:"comment_id"`
		Content   string   `json:"content"`
	}
	_ = c.ShouldBindJSON(&body)
	content := strings.TrimSpace(body.Content)
	if body.CommentID.Int() <= 0 {
		response.Fail(c, 1, "请指定要回复的评论")
		return
	}
	if content == "" {
		response.Fail(c, 1, "回复内容不能为空")
		return
	}
	ctx := c.Request.Context()
	now := time.Now()
	if chapter {
		var parent chapterCommentRow
		err := s.db.WithContext(ctx).Where("id = ?", body.CommentID.Int()).Take(&parent).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "评论不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		row := chapterCommentRow{
			PID: parent.ID, CourseID: parent.CourseID, ChapterID: parent.ChapterID,
			Content: content, UserID: 0, CreatedAt: &now,
		}
		if response.WriteErr(c, s.db.WithContext(ctx).Create(&row).Error) {
			return
		}
		courseTitle := s.comments.TitleOf(ctx, "courses", parent.CourseID, "课程")
		chapterTitle := s.comments.TitleOf(ctx, "course_chapters", parent.ChapterID, "课节")
		_ = s.comments.SendNotice(ctx, parent.UserID, "评论收到回复",
			replyHTML("课节「"+chapterTitle+"」（课程「"+courseTitle+"」）", parent.Content, content),
			"chapter", parent.ChapterID, middleware.UID(c))
		response.OK(c, "回复成功", []any{})
		return
	}
	var parent commentRow
	err := s.db.WithContext(ctx).Where("id = ?", body.CommentID.Int()).Take(&parent).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "评论不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	row := commentRow{
		PID: parent.ID, CourseID: parent.CourseID, Content: content, UserID: 0, CreatedAt: &now,
	}
	if response.WriteErr(c, s.db.WithContext(ctx).Create(&row).Error) {
		return
	}
	courseTitle := s.comments.TitleOf(ctx, "courses", parent.CourseID, "课程")
	_ = s.comments.SendNotice(ctx, parent.UserID, "评论收到回复",
		replyHTML("课程「"+courseTitle+"」", parent.Content, content),
		"course", parent.CourseID, middleware.UID(c))
	response.OK(c, "回复成功", []any{})
}

func (s *Service) GetCourseCommentReplies(c *gin.Context) {
	s.replies(c, false)
}

func (s *Service) GetChapterCommentReplies(c *gin.Context) {
	s.replies(c, true)
}

func (s *Service) replies(c *gin.Context, chapter bool) {
	id := atoi(c.Query("comment_id"))
	if id <= 0 {
		response.Fail(c, 1, "参数错误")
		return
	}
	ctx := c.Request.Context()
	if chapter {
		var parent chapterCommentRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&parent).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "评论不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		var rows []chapterCommentRow
		if response.WriteErr(c, s.db.WithContext(ctx).Where("pid = ?", id).Order("id DESC").Find(&rows).Error) {
			return
		}
		list, err := s.decorateChapter(c, rows)
		if response.WriteErr(c, err) {
			return
		}
		response.OK(c, "", list)
		return
	}
	var parent commentRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&parent).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "评论不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	var rows []commentRow
	if response.WriteErr(c, s.db.WithContext(ctx).Where("pid = ?", id).Order("id DESC").Find(&rows).Error) {
		return
	}
	list, err := s.decorateCourse(c, rows)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "", list)
}

func (s *Service) CourseCommentDelete(c *gin.Context) {
	s.deleteOne(c, false)
}

func (s *Service) CourseChapterCommentDelete(c *gin.Context) {
	s.deleteOne(c, true)
}

func (s *Service) deleteOne(c *gin.Context, chapter bool) {
	id := atoi(c.Query("id"))
	if id <= 0 {
		response.Fail(c, 1, "评论不存在")
		return
	}
	ctx := c.Request.Context()
	if chapter {
		var row chapterCommentRow
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "评论不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		_ = s.db.WithContext(ctx).Where("pid = ?", id).Delete(&chapterCommentRow{}).Error
		if response.WriteErr(c, s.db.WithContext(ctx).Delete(&row).Error) {
			return
		}
		response.OK(c, "删除成功", []any{})
		return
	}
	var row commentRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "评论不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	_ = s.db.WithContext(ctx).Where("pid = ?", id).Delete(&commentRow{}).Error
	if response.WriteErr(c, s.db.WithContext(ctx).Delete(&row).Error) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (s *Service) CourseCommentBatchDelete(c *gin.Context) {
	s.batchDelete(c, false)
}

func (s *Service) CourseChapterCommentBatchDelete(c *gin.Context) {
	s.batchDelete(c, true)
}

func (s *Service) batchDelete(c *gin.Context, chapter bool) {
	var body struct {
		IDs []bind.Int `json:"ids"`
	}
	_ = c.ShouldBindJSON(&body)
	ids := uniquePositive(bind.Ints(body.IDs))
	if len(ids) == 0 {
		response.Fail(c, 1, "请选择要删除的评论")
		return
	}
	ctx := c.Request.Context()
	if chapter {
		_ = s.db.WithContext(ctx).Where("pid IN ?", ids).Delete(&chapterCommentRow{}).Error
		result := s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&chapterCommentRow{})
		if response.WriteErr(c, result.Error) {
			return
		}
		if result.RowsAffected == 0 {
			response.Fail(c, 1, "评论不存在")
			return
		}
		response.OK(c, "删除成功", []any{})
		return
	}
	_ = s.db.WithContext(ctx).Where("pid IN ?", ids).Delete(&commentRow{}).Error
	result := s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&commentRow{})
	if response.WriteErr(c, result.Error) {
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, 1, "评论不存在")
		return
	}
	response.OK(c, "删除成功", []any{})
}
