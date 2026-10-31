package webapp

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// PaperRepository covers learner exam/assessment reads and writes.
type PaperRepository interface {
	FindPassedID(ctx context.Context, uid, paperID int) (id int, err error)
	FindLatestAttempt(ctx context.Context, uid, paperID int) (*UserPaperAttempt, error)
	FindUserPaper(ctx context.Context, id int) (*UserPaperLog, error)
	CountPassed(ctx context.Context, uid, paperID int) (int64, error)
	FindMeta(ctx context.Context, id int) (*PaperMeta, error)
	FindCourseBrief(ctx context.Context, id int) (courseID int, title string, err error)
	FindChapterBrief(ctx context.Context, id int) (chapterID int, title string, err error)
	ListQuestions(ctx context.Context, paperID int) ([]PaperQuestionRow, error)
	ListQuestionLogs(ctx context.Context, userPaperID int) ([]map[string]any, error)
	CreateUserPaper(ctx context.Context, uid, paperID int, now time.Time) (id int, err error)
	CreateQuestionLogs(ctx context.Context, rows []map[string]any) error
	UpdateUserPaper(ctx context.Context, id int, fields map[string]any) error
	FindUserChapter(ctx context.Context, uid, chapterID int) (*UserChapterRef, error)
	MarkChapterEnded(ctx context.Context, id int) error
	IncrCourseHasEnd(ctx context.Context, uid, courseID int) error
	ChapterTitle(ctx context.Context, id int) (string, error)
	MarkCourseEnded(ctx context.Context, uid, courseID int) error
	FindCourseCredit(ctx context.Context, id int) (title string, credit int, err error)
}

type PaperMeta struct {
	ID        int
	CourseID  *int
	ChapterID *int
	PassScore int
	Score     int
	Type      string
}

type UserPaperAttempt struct {
	ID            int `gorm:"column:id"`
	IsPass        int `gorm:"column:is_pass"`
	IsPeopleCheck int `gorm:"column:is_people_check"`
}

type UserPaperLog struct {
	ID            int `gorm:"column:id"`
	UserID        int `gorm:"column:user_id"`
	PaperID       int `gorm:"column:paper_id"`
	IsPass        int `gorm:"column:is_pass"`
	IsPeopleCheck int `gorm:"column:is_people_check"`
}

type PaperQuestionRow struct {
	ID        int64           `gorm:"column:id"`
	PaperID   int             `gorm:"column:paper_id"`
	Sort      int             `gorm:"column:sort"`
	Type      string          `gorm:"column:type"`
	Title     string          `gorm:"column:title"`
	Answer    json.RawMessage `gorm:"column:answer"`
	Score     int             `gorm:"column:score"`
	Options   json.RawMessage `gorm:"column:options"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at"`
}

type UserChapterRef struct {
	ID       int  `gorm:"column:id"`
	CourseID *int `gorm:"column:course_id"`
}

type paperRepo struct{ db *gorm.DB }

func NewPaper(db *gorm.DB) PaperRepository {
	return &paperRepo{db: db}
}

func (r *paperRepo) FindPassedID(ctx context.Context, uid, paperID int) (int, error) {
	var row struct {
		ID int `gorm:"column:id"`
	}
	err := r.db.WithContext(ctx).Table("user_paper").
		Select("id").Where("user_id = ? AND paper_id = ? AND is_pass = 1 AND deleted_at IS NULL", uid, paperID).
		Order("id DESC").Take(&row).Error
	if err != nil {
		return 0, err
	}
	return row.ID, nil
}

func (r *paperRepo) FindLatestAttempt(ctx context.Context, uid, paperID int) (*UserPaperAttempt, error) {
	var row UserPaperAttempt
	err := r.db.WithContext(ctx).Table("user_paper").
		Where("user_id = ? AND paper_id = ? AND deleted_at IS NULL", uid, paperID).
		Order("id DESC").Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *paperRepo) FindUserPaper(ctx context.Context, id int) (*UserPaperLog, error) {
	var row UserPaperLog
	err := r.db.WithContext(ctx).Table("user_paper").Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *paperRepo) CountPassed(ctx context.Context, uid, paperID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_paper").
		Where("user_id = ? AND paper_id = ? AND is_pass = 1 AND deleted_at IS NULL", uid, paperID).
		Count(&n).Error
	return n, err
}

func (r *paperRepo) FindMeta(ctx context.Context, id int) (*PaperMeta, error) {
	var row PaperMeta
	err := r.db.WithContext(ctx).Table("papers").
		Select("id, course_id, chapter_id, score, pass_score, type").
		Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *paperRepo) FindCourseBrief(ctx context.Context, id int) (int, string, error) {
	var row struct {
		ID    int    `gorm:"column:id"`
		Title string `gorm:"column:title"`
	}
	err := r.db.WithContext(ctx).Table("courses").Select("id, title").
		Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return 0, "", err
	}
	return row.ID, row.Title, nil
}

func (r *paperRepo) FindChapterBrief(ctx context.Context, id int) (int, string, error) {
	var row struct {
		ID    int    `gorm:"column:id"`
		Title string `gorm:"column:title"`
	}
	err := r.db.WithContext(ctx).Table("course_chapters").Select("id, title").
		Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return 0, "", err
	}
	return row.ID, row.Title, nil
}

func (r *paperRepo) ListQuestions(ctx context.Context, paperID int) ([]PaperQuestionRow, error) {
	var rows []PaperQuestionRow
	err := r.db.WithContext(ctx).Table("paper_questions").
		Where("paper_id = ?", paperID).Order("sort ASC").Find(&rows).Error
	return rows, err
}

func (r *paperRepo) ListQuestionLogs(ctx context.Context, userPaperID int) ([]map[string]any, error) {
	var rows []map[string]any
	err := r.db.WithContext(ctx).Table("user_paper_question_log").
		Where("user_paper_id = ?", userPaperID).Find(&rows).Error
	return rows, err
}

func (r *paperRepo) CreateUserPaper(ctx context.Context, uid, paperID int, now time.Time) (int, error) {
	row := struct {
		ID            int       `gorm:"column:id;primaryKey"`
		UserID        int       `gorm:"column:user_id"`
		PaperID       int       `gorm:"column:paper_id"`
		Score         int       `gorm:"column:score"`
		IsPass        int       `gorm:"column:is_pass"`
		IsPeopleCheck int       `gorm:"column:is_people_check"`
		CreatedAt     time.Time `gorm:"column:created_at"`
		UpdatedAt     time.Time `gorm:"column:updated_at"`
	}{
		UserID: uid, PaperID: paperID, CreatedAt: now, UpdatedAt: now,
	}
	err := r.db.WithContext(ctx).Table("user_paper").Create(&row).Error
	return row.ID, err
}

func (r *paperRepo) CreateQuestionLogs(ctx context.Context, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Table("user_paper_question_log").Create(&rows).Error
}

func (r *paperRepo) UpdateUserPaper(ctx context.Context, id int, fields map[string]any) error {
	return r.db.WithContext(ctx).Table("user_paper").Where("id = ?", id).Updates(fields).Error
}

func (r *paperRepo) FindUserChapter(ctx context.Context, uid, chapterID int) (*UserChapterRef, error) {
	var row UserChapterRef
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Select("id, course_id").
		Where("user_id = ? AND chapter_id = ?", uid, chapterID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *paperRepo) MarkChapterEnded(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Table("user_course_chapter").Where("id = ?", id).
		Updates(map[string]any{"is_end": 1, "has_end_chapter": 1}).Error
}

func (r *paperRepo) IncrCourseHasEnd(ctx context.Context, uid, courseID int) error {
	return r.db.WithContext(ctx).Exec(
		"UPDATE user_course SET has_end_chapter = has_end_chapter + 1 WHERE user_id = ? AND course_id = ?",
		uid, courseID,
	).Error
}

func (r *paperRepo) ChapterTitle(ctx context.Context, id int) (string, error) {
	var title string
	err := r.db.WithContext(ctx).Table("course_chapters").Select("title").
		Where("id = ? AND deleted_at IS NULL", id).Scan(&title).Error
	return title, err
}

func (r *paperRepo) MarkCourseEnded(ctx context.Context, uid, courseID int) error {
	return r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", uid, courseID).
		Update("is_end", 1).Error
}

func (r *paperRepo) FindCourseCredit(ctx context.Context, id int) (string, int, error) {
	var row struct {
		Title  string `gorm:"column:title"`
		Credit int    `gorm:"column:credit"`
	}
	err := r.db.WithContext(ctx).Table("courses").Select("title, credit").
		Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return "", 0, err
	}
	return row.Title, row.Credit, nil
}
