package webapp

import (
	"context"
	"time"

	"gorm.io/gorm"

	"qixue-api/internal/model"
)

// CommentRepository 学员端课程评论 / 课节评论（两张表、两套模型，禁止混用）。
type CommentRepository interface {
	CountUserCourse(ctx context.Context, uid, courseID int) (int64, error)

	CountCourseRoot(ctx context.Context, courseID int) (int64, error)
	ListCourseRoot(ctx context.Context, courseID, page, limit int) ([]model.CourseComment, error)
	ListCourseChildren(ctx context.Context, pids []int) ([]model.CourseComment, error)
	CreateCourseComment(ctx context.Context, row *model.CourseComment) error
	CountUserCourseComments(ctx context.Context, uid int) (int64, error)
	ListUserCourseComments(ctx context.Context, uid, page, limit int) ([]model.CourseComment, error)

	CountChapterRoot(ctx context.Context, chapterID int) (int64, error)
	ListChapterRoot(ctx context.Context, chapterID, page, limit int) ([]model.CourseChapterComment, error)
	ListChapterChildren(ctx context.Context, pids []int) ([]model.CourseChapterComment, error)
	CreateChapterComment(ctx context.Context, row *model.CourseChapterComment) error
	CountUserChapterComments(ctx context.Context, uid int) (int64, error)
	ListUserChapterComments(ctx context.Context, uid, page, limit int) ([]model.CourseChapterComment, error)

	ListUserBriefs(ctx context.Context, ids []int) ([]UserBriefRow, error)
	ListCourseBriefs(ctx context.Context, ids []int) ([]CourseBrief, error)
	ListChapterBriefs(ctx context.Context, ids []int) ([]ChapterBrief, error)
}

type commentRepo struct{ db *gorm.DB }

func NewComment(db *gorm.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) CountUserCourse(ctx context.Context, uid, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", uid, courseID).Count(&n).Error
	return n, err
}

func (r *commentRepo) CountCourseRoot(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.CourseComment{}).
		Where("course_id = ? AND pid = 0", courseID).
		Count(&n).Error
	return n, err
}

func (r *commentRepo) ListCourseRoot(ctx context.Context, courseID, page, limit int) ([]model.CourseComment, error) {
	var rows []model.CourseComment
	err := r.db.WithContext(ctx).Model(&model.CourseComment{}).
		Where("course_id = ? AND pid = 0", courseID).
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *commentRepo) ListCourseChildren(ctx context.Context, pids []int) ([]model.CourseComment, error) {
	if len(pids) == 0 {
		return nil, nil
	}
	var rows []model.CourseComment
	err := r.db.WithContext(ctx).Model(&model.CourseComment{}).
		Where("pid IN ?", pids).
		Order("id DESC").Find(&rows).Error
	return rows, err
}

func (r *commentRepo) CreateCourseComment(ctx context.Context, row *model.CourseComment) error {
	now := time.Now()
	if row.CreatedAt == nil {
		row.CreatedAt = &now
	}
	if row.UpdatedAt == nil {
		row.UpdatedAt = &now
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *commentRepo) CountUserCourseComments(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.CourseComment{}).
		Where("user_id = ?", uid).Count(&n).Error
	return n, err
}

func (r *commentRepo) ListUserCourseComments(ctx context.Context, uid, page, limit int) ([]model.CourseComment, error) {
	var rows []model.CourseComment
	err := r.db.WithContext(ctx).Model(&model.CourseComment{}).
		Where("user_id = ?", uid).
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *commentRepo) CountChapterRoot(ctx context.Context, chapterID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.CourseChapterComment{}).
		Where("chapter_id = ? AND pid = 0", chapterID).
		Count(&n).Error
	return n, err
}

func (r *commentRepo) ListChapterRoot(ctx context.Context, chapterID, page, limit int) ([]model.CourseChapterComment, error) {
	var rows []model.CourseChapterComment
	err := r.db.WithContext(ctx).Model(&model.CourseChapterComment{}).
		Where("chapter_id = ? AND pid = 0", chapterID).
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *commentRepo) ListChapterChildren(ctx context.Context, pids []int) ([]model.CourseChapterComment, error) {
	if len(pids) == 0 {
		return nil, nil
	}
	var rows []model.CourseChapterComment
	err := r.db.WithContext(ctx).Model(&model.CourseChapterComment{}).
		Where("pid IN ?", pids).
		Order("id DESC").Find(&rows).Error
	return rows, err
}

func (r *commentRepo) CreateChapterComment(ctx context.Context, row *model.CourseChapterComment) error {
	now := time.Now()
	if row.CreatedAt == nil {
		row.CreatedAt = &now
	}
	if row.UpdatedAt == nil {
		row.UpdatedAt = &now
	}
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *commentRepo) CountUserChapterComments(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.CourseChapterComment{}).
		Where("user_id = ?", uid).Count(&n).Error
	return n, err
}

func (r *commentRepo) ListUserChapterComments(ctx context.Context, uid, page, limit int) ([]model.CourseChapterComment, error) {
	var rows []model.CourseChapterComment
	err := r.db.WithContext(ctx).Model(&model.CourseChapterComment{}).
		Where("user_id = ?", uid).
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).
		Find(&rows).Error
	return rows, err
}

type UserBriefRow struct {
	ID       int    `gorm:"column:id"`
	NickName string `gorm:"column:nick_name"`
	Avatar   string `gorm:"column:avatar"`
	Coins    int    `gorm:"column:coins"`
}

func (r *commentRepo) ListUserBriefs(ctx context.Context, ids []int) ([]UserBriefRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []UserBriefRow
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("id, nick_name, avatar, coins").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *commentRepo) ListCourseBriefs(ctx context.Context, ids []int) ([]CourseBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []CourseBrief
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb, `desc`").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *commentRepo) ListChapterBriefs(ctx context.Context, ids []int) ([]ChapterBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []ChapterBrief
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, title").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}
