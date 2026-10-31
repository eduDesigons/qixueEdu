package admincomment

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/page"
	"qixue-api/internal/model"
	"qixue-api/internal/platform/bg"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/scope"
)

var (
	ErrBadParams     = errors.New("参数错误")
	ErrNeedCommentID = errors.New("请指定要回复的评论")
	ErrEmptyReply    = errors.New("回复内容不能为空")
	ErrNotFound      = errors.New("评论不存在")
	ErrNeedIDs       = errors.New("请选择要删除的评论")
)

type Service struct {
	comments adminrepo.AdminCommentRepository
	db       *gorm.DB
	rdb      *redis.Client
}

func New(comments adminrepo.AdminCommentRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{comments: comments, db: db, rdb: rdb}
}

type UserBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
}

type TitleBrief struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	CourseID int    `json:"course_id,omitempty"`
}

type CourseCommentView struct {
	ID            int         `json:"id"`
	PID           int         `json:"pid"`
	CourseID      int         `json:"course_id"`
	UserID        int         `json:"user_id"`
	Content       string      `json:"content"`
	CreatedAt     *string     `json:"created_at"`
	ChildrenCount int         `json:"children_count"`
	User          *UserBrief  `json:"user"`
	Course        *TitleBrief `json:"course"`
}

type ChapterCommentView struct {
	ID            int         `json:"id"`
	PID           int         `json:"pid"`
	CourseID      int         `json:"course_id"`
	ChapterID     int         `json:"chapter_id"`
	UserID        int         `json:"user_id"`
	Content       string      `json:"content"`
	CreatedAt     *string     `json:"created_at"`
	ChildrenCount int         `json:"children_count"`
	User          *UserBrief  `json:"user"`
	Course        *TitleBrief `json:"course"`
	Chapter       *TitleBrief `json:"chapter"`
}

type CourseListQuery struct {
	CourseID    int
	CourseTitle string
	NickName    string
	Content     string
	StartTime   string
	EndTime     string
	Page        int
	Size        int
}

type ChapterListQuery struct {
	CourseID     int
	ChapterID    int
	CourseTitle  string
	ChapterTitle string
	NickName     string
	Content      string
	StartTime    string
	EndTime      string
	Page         int
	Size         int
}

type ReplyInput struct {
	CommentID int
	Content   string
}

func (s *Service) ListCourse(ctx context.Context, adminID int, q CourseListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	query := s.db.WithContext(ctx).Model(&model.CourseComment{}).Where("pid = 0")
	if q.CourseID > 0 {
		query = query.Where("course_id = ?", q.CourseID)
	}
	query, ok, err := s.applyTitleIDs(ctx, query, "course_id", "courses", q.CourseTitle)
	if err != nil {
		return page.Body{}, err
	}
	if !ok {
		return page.New(1, size, 0, []any{}), nil
	}
	query, ok, err = s.applyNickIDs(ctx, query, q.NickName)
	if err != nil {
		return page.Body{}, err
	}
	if !ok {
		return page.New(1, size, 0, []any{}), nil
	}
	query = applyContent(query, "course_comments", q.Content)
	query = applyTime(query, q.StartTime, q.EndTime)
	if limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []model.CourseComment
	if err := query.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	list, err := s.decorateCourse(ctx, rows)
	if err != nil {
		return page.Body{}, err
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) ListChapter(ctx context.Context, adminID int, q ChapterListQuery) (page.Body, error) {
	current, size := page.Clamp(q.Page, q.Size)
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return page.Body{}, err
	}
	query := s.db.WithContext(ctx).Model(&model.CourseChapterComment{}).Where("pid = 0")
	if q.CourseID > 0 {
		query = query.Where("course_id = ?", q.CourseID)
	}
	if q.ChapterID > 0 {
		query = query.Where("chapter_id = ?", q.ChapterID)
	}
	query, ok, err := s.applyTitleIDs(ctx, query, "course_id", "courses", q.CourseTitle)
	if err != nil {
		return page.Body{}, err
	}
	if !ok {
		return page.New(1, size, 0, []any{}), nil
	}
	query, ok, err = s.applyTitleIDs(ctx, query, "chapter_id", "course_chapters", q.ChapterTitle)
	if err != nil {
		return page.Body{}, err
	}
	if !ok {
		return page.New(1, size, 0, []any{}), nil
	}
	query, ok, err = s.applyNickIDs(ctx, query, q.NickName)
	if err != nil {
		return page.Body{}, err
	}
	if !ok {
		return page.New(1, size, 0, []any{}), nil
	}
	query = applyContent(query, "course_chapter_comments", q.Content)
	query = applyTime(query, q.StartTime, q.EndTime)
	if limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return page.Body{}, err
	}
	var rows []model.CourseChapterComment
	if err := query.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return page.Body{}, err
	}
	list, err := s.decorateChapter(ctx, rows)
	if err != nil {
		return page.Body{}, err
	}
	return page.New(current, size, total, list), nil
}

func (s *Service) ReplyCourse(ctx context.Context, adminID int, in ReplyInput) error {
	return s.reply(ctx, adminID, in, false)
}

func (s *Service) ReplyChapter(ctx context.Context, adminID int, in ReplyInput) error {
	return s.reply(ctx, adminID, in, true)
}

func (s *Service) reply(ctx context.Context, adminID int, in ReplyInput, chapter bool) error {
	content := strings.TrimSpace(in.Content)
	if in.CommentID <= 0 {
		return ErrNeedCommentID
	}
	if content == "" {
		return ErrEmptyReply
	}
	now := time.Now()
	if chapter {
		var parent model.CourseChapterComment
		err := s.db.WithContext(ctx).Where("id = ?", in.CommentID).Take(&parent).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		row := model.CourseChapterComment{
			PID: int(parent.ID), CourseID: parent.CourseID, ChapterID: parent.ChapterID,
			Content: content, UserID: 0, CreatedAt: &now, UpdatedAt: &now,
		}
		if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
			return err
		}
		noticeParent := parent
		bg.Go(ctx, func(ctx context.Context) {
			courseTitle := s.comments.TitleOf(ctx, "courses", noticeParent.CourseID, "课程")
			chapterTitle := s.comments.TitleOf(ctx, "course_chapters", noticeParent.ChapterID, "课节")
			_ = s.comments.SendNotice(ctx, noticeParent.UserID, "评论收到回复",
				replyHTML("课节「"+chapterTitle+"」（课程「"+courseTitle+"」）", noticeParent.Content, content),
				"chapter", noticeParent.ChapterID, adminID)
		})
		return nil
	}
	var parent model.CourseComment
	err := s.db.WithContext(ctx).Where("id = ?", in.CommentID).Take(&parent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	row := model.CourseComment{
		PID: int(parent.ID), CourseID: parent.CourseID, Content: content, UserID: 0,
		CreatedAt: &now, UpdatedAt: &now,
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return err
	}
	noticeParent := parent
	bg.Go(ctx, func(ctx context.Context) {
		courseTitle := s.comments.TitleOf(ctx, "courses", noticeParent.CourseID, "课程")
		_ = s.comments.SendNotice(ctx, noticeParent.UserID, "评论收到回复",
			replyHTML("课程「"+courseTitle+"」", noticeParent.Content, content),
			"course", noticeParent.CourseID, adminID)
	})
	return nil
}

func (s *Service) CourseReplies(ctx context.Context, commentID int) ([]CourseCommentView, error) {
	if commentID <= 0 {
		return nil, ErrBadParams
	}
	var parent model.CourseComment
	err := s.db.WithContext(ctx).Where("id = ?", commentID).Take(&parent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var rows []model.CourseComment
	if err := s.db.WithContext(ctx).Where("pid = ?", commentID).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return s.decorateCourse(ctx, rows)
}

func (s *Service) ChapterReplies(ctx context.Context, commentID int) ([]ChapterCommentView, error) {
	if commentID <= 0 {
		return nil, ErrBadParams
	}
	var parent model.CourseChapterComment
	err := s.db.WithContext(ctx).Where("id = ?", commentID).Take(&parent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var rows []model.CourseChapterComment
	if err := s.db.WithContext(ctx).Where("pid = ?", commentID).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return s.decorateChapter(ctx, rows)
}

func (s *Service) DeleteCourse(ctx context.Context, id int) error {
	return s.deleteOne(ctx, id, false)
}

func (s *Service) DeleteChapter(ctx context.Context, id int) error {
	return s.deleteOne(ctx, id, true)
}

func (s *Service) deleteOne(ctx context.Context, id int, chapter bool) error {
	if id <= 0 {
		return ErrNotFound
	}
	if chapter {
		var row model.CourseChapterComment
		err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		_ = s.db.WithContext(ctx).Where("pid = ?", id).Delete(&model.CourseChapterComment{}).Error
		return s.db.WithContext(ctx).Delete(&row).Error
	}
	var row model.CourseComment
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	_ = s.db.WithContext(ctx).Where("pid = ?", id).Delete(&model.CourseComment{}).Error
	return s.db.WithContext(ctx).Delete(&row).Error
}

func (s *Service) BatchDeleteCourse(ctx context.Context, ids []int) error {
	return s.batchDelete(ctx, ids, false)
}

func (s *Service) BatchDeleteChapter(ctx context.Context, ids []int) error {
	return s.batchDelete(ctx, ids, true)
}

func (s *Service) batchDelete(ctx context.Context, ids []int, chapter bool) error {
	ids = uniquePositive(ids)
	if len(ids) == 0 {
		return ErrNeedIDs
	}
	if chapter {
		_ = s.db.WithContext(ctx).Where("pid IN ?", ids).Delete(&model.CourseChapterComment{}).Error
		result := s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.CourseChapterComment{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	}
	_ = s.db.WithContext(ctx).Where("pid IN ?", ids).Delete(&model.CourseComment{}).Error
	result := s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.CourseComment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
