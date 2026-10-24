package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// CommentRepository covers course/chapter discussion reads and writes.
type CommentRepository interface {
	CountRoot(ctx context.Context, table string, withChapter bool, targetID int) (int64, error)
	ListRoot(ctx context.Context, table string, withChapter bool, targetID, page, limit int) ([]CommentRow, error)
	ListChildren(ctx context.Context, table string, pids []int) ([]CommentRow, error)
	CountUserCourse(ctx context.Context, uid, courseID int) (int64, error)
	CreateChapterComment(ctx context.Context, row *CommentRow) error
	CreateCourseComment(ctx context.Context, row *CommentRow) error
	ListUserBriefs(ctx context.Context, ids []int) ([]map[string]any, error)
}

type commentRepo struct{ db *gorm.DB }

func NewComment(db *gorm.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) CountRoot(ctx context.Context, table string, withChapter bool, targetID int) (int64, error) {
	q := r.db.WithContext(ctx).Table(table).Where("pid = 0 AND deleted_at IS NULL")
	if withChapter {
		q = q.Where("chapter_id = ?", targetID)
	} else {
		q = q.Where("course_id = ?", targetID)
	}
	var n int64
	err := q.Count(&n).Error
	return n, err
}

func (r *commentRepo) ListRoot(ctx context.Context, table string, withChapter bool, targetID, page, limit int) ([]CommentRow, error) {
	q := r.db.WithContext(ctx).Table(table).Where("pid = 0 AND deleted_at IS NULL")
	if withChapter {
		q = q.Where("chapter_id = ?", targetID)
	} else {
		q = q.Where("course_id = ?", targetID)
	}
	var rows []CommentRow
	err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *commentRepo) ListChildren(ctx context.Context, table string, pids []int) ([]CommentRow, error) {
	if len(pids) == 0 {
		return nil, nil
	}
	var rows []CommentRow
	err := r.db.WithContext(ctx).Table(table).
		Where("pid IN ? AND deleted_at IS NULL", pids).
		Order("id DESC").Find(&rows).Error
	return rows, err
}

func (r *commentRepo) CountUserCourse(ctx context.Context, uid, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", uid, courseID).Count(&n).Error
	return n, err
}

func (r *commentRepo) CreateChapterComment(ctx context.Context, row *CommentRow) error {
	if row.CreatedAt == nil {
		now := time.Now()
		row.CreatedAt = &now
		row.UpdatedAt = &now
	}
	return r.db.WithContext(ctx).Table("course_chapter_comments").Create(row).Error
}

func (r *commentRepo) CreateCourseComment(ctx context.Context, row *CommentRow) error {
	if row.CreatedAt == nil {
		now := time.Now()
		row.CreatedAt = &now
		row.UpdatedAt = &now
	}
	return r.db.WithContext(ctx).Table("course_comments").Create(row).Error
}

func (r *commentRepo) ListUserBriefs(ctx context.Context, ids []int) ([]map[string]any, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []map[string]any
	err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, avatar, coins").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}
