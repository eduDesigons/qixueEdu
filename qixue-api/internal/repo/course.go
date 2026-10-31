package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CourseRepository wraps course/path/enrollment access for enroll, learn, and webapp catalog.
type CourseRepository interface {
	CountChapters(ctx context.Context, courseID int) (int64, error)
	ListChapterIDs(ctx context.Context, courseID int) ([]int, error)
	ListUserIDsInCourse(ctx context.Context, courseID int) ([]int, error)
	ListUserIDsInCourseAmong(ctx context.Context, courseID int, userIDs []int) ([]int, error)
	CreateUserCoursesIgnore(ctx context.Context, rows []map[string]any) error
	ListUserChapterIDs(ctx context.Context, uid, courseID int) ([]int, error)
	CreateUserChaptersIgnore(ctx context.Context, rows []map[string]any) error
	ListUserChapterProgress(ctx context.Context, uid, courseID int) ([]ChapterProgressRow, error)
	SetChapterUnlock(ctx context.Context, ids []int, unlock int) error
	FindUserCourseStatus(ctx context.Context, uid, courseID int) (*UserCourseStatus, error)
	FindUserCourse(ctx context.Context, uid, courseID int) (*UserCourseStatus, error)
	UpdateUserCourseStatus(ctx context.Context, id, hasEnd, isEnd int) error
	CreditGranted(ctx context.Context, uid int, field, remark string) (bool, error)

	ListPathIDsByCourse(ctx context.Context, courseID int) ([]int, error)
	FindUserPathStatus(ctx context.Context, uid, pathID int) (*UserPathStatus, error)
	ListCourseIDsByPath(ctx context.Context, pathID int) ([]int, error)
	CountEndedUserCourses(ctx context.Context, uid int, courseIDs []int) (int, error)
	UpdateUserPathStatus(ctx context.Context, id, hasEnd, isEnd int) error
	ListUserIDsOnPathAmong(ctx context.Context, pathID int, userIDs []int) ([]int, error)
	CreateUserPathsIgnore(ctx context.Context, rows []map[string]any) error

	ListShownCourses(ctx context.Context, f CourseListFilter) (total int64, rows []CourseListRow, err error)
	FindCourseDetail(ctx context.Context, id int) (*CourseDetailRow, error)
	CountPendingCourseApply(ctx context.Context, uid, courseID int) (int64, error)
	CountCourseLike(ctx context.Context, uid, courseID int) (int64, error)
	ListCourseAttaches(ctx context.Context, courseID int) ([]CourseAttachRow, error)
	FindValidCourseCert(ctx context.Context, uid, courseID int) (*UserCertBrief, error)
	IncrCourseView(ctx context.Context, id int) error
	ListChaptersWithMedia(ctx context.Context, courseID int) ([]ChapterListRow, error)
	FindChapterDetail(ctx context.Context, id int) (*ChapterDetailRow, error)
	FindUserChapter(ctx context.Context, uid, chapterID int) (*CatalogUserChapterRow, error)
	CountChapterLike(ctx context.Context, uid, chapterID int) (int64, error)
	IncrChapterView(ctx context.Context, id int) error
	FindLatestUserPaper(ctx context.Context, uid, paperID int) (*LatestUserPaper, error)
	ListPendingUserPapers(ctx context.Context, uid int, paperIDs []int) ([]PendingUserPaper, error)
	FindPrevChapter(ctx context.Context, courseID, chapterID, sort int) (*ChapterNeighborRow, error)
	FindNextChapter(ctx context.Context, courseID, chapterID, sort int) (*ChapterNeighborRow, error)
	ListNeighborProgress(ctx context.Context, uid int, ids []int) ([]NeighborProgressRow, error)
	UnlockUserChapter(ctx context.Context, uid, chapterID int) error

	ListShownPaths(ctx context.Context, categoryID, page, limit int) (total int64, rows []PathListRow, err error)
	FindPathDetail(ctx context.Context, id int) (*PathDetailRow, error)
	ListPathCourses(ctx context.Context, pathID int) ([]PathCourseBrief, error)
	CountPathUsers(ctx context.Context, pathID int) (int64, error)
	FindUserPath(ctx context.Context, uid, pathID int) (*UserPathBrief, error)
	CountPendingPathApply(ctx context.Context, uid, pathID int) (int64, error)
	FindValidPathCert(ctx context.Context, uid, pathID int) (*UserCertBrief, error)
	FindPathMinCredit(ctx context.Context, id int) (pathID, minCredit int, err error)
	CountUserOnPath(ctx context.Context, uid, pathID int) (int64, error)
	PathExists(ctx context.Context, id int) (bool, error)
	CreatePathApply(ctx context.Context, uid, pathID int, reason string) error
}

type ChapterProgressRow struct {
	ID            int `gorm:"column:id"`
	ChapterID     int `gorm:"column:chapter_id"`
	IsEnd         int `gorm:"column:is_end"`
	HasEndChapter int `gorm:"column:has_end_chapter"`
	UnLock        int `gorm:"column:un_lock"`
}

// UserChapterProgress is kept as an alias for enroll callers.
type UserChapterProgress = ChapterProgressRow

type LatestUserPaper struct {
	ID            int `gorm:"column:id"`
	IsPass        int `gorm:"column:is_pass"`
	IsPeopleCheck int `gorm:"column:is_people_check"`
}

type PendingUserPaper struct {
	ID      int `gorm:"column:id"`
	PaperID int `gorm:"column:paper_id"`
}

type UserCourseStatus struct {
	ID            int `gorm:"column:id"`
	HasEndChapter int `gorm:"column:has_end_chapter"`
	IsEnd         int `gorm:"column:is_end"`
}

type UserPathStatus struct {
	ID           int `gorm:"column:id"`
	HasEndCourse int `gorm:"column:has_end_course"`
	IsEnd        int `gorm:"column:is_end"`
}

type courseRepo struct{ db *gorm.DB }

func NewCourse(db *gorm.DB) CourseRepository {
	return &courseRepo{db: db}
}

func (r *courseRepo) CountChapters(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("course_chapters").
		Where("course_id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n, err
}

func (r *courseRepo) ListChapterIDs(ctx context.Context, courseID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("course_chapters").
		Where("course_id = ? AND deleted_at IS NULL", courseID).
		Order("sort ASC, id ASC").
		Pluck("id", &ids).Error
	return ids, err
}

func (r *courseRepo) ListUserIDsInCourse(ctx context.Context, courseID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("user_course").
		Where("course_id = ?", courseID).
		Pluck("user_id", &ids).Error
	return ids, err
}

func (r *courseRepo) ListUserIDsInCourseAmong(ctx context.Context, courseID int, userIDs []int) ([]int, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var ids []int
	err := r.db.WithContext(ctx).Table("user_course").
		Where("course_id = ? AND user_id IN ?", courseID, userIDs).
		Pluck("user_id", &ids).Error
	return ids, err
}

func (r *courseRepo) CreateUserCoursesIgnore(ctx context.Context, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).
		Table("user_course").Create(&rows).Error
}

func (r *courseRepo) ListUserChapterIDs(ctx context.Context, uid, courseID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND chapter_id IN (SELECT id FROM course_chapters WHERE course_id = ? AND deleted_at IS NULL)", uid, courseID).
		Pluck("chapter_id", &ids).Error
	return ids, err
}

func (r *courseRepo) CreateUserChaptersIgnore(ctx context.Context, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).
		Table("user_course_chapter").Create(&rows).Error
}

func (r *courseRepo) ListUserChapterProgress(ctx context.Context, uid, courseID int) ([]ChapterProgressRow, error) {
	var rows []ChapterProgressRow
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND course_id = ?", uid, courseID).
		Find(&rows).Error
	return rows, err
}

func (r *courseRepo) SetChapterUnlock(ctx context.Context, ids []int, unlock int) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Table("user_course_chapter").
		Where("id IN ?", ids).
		UpdateColumn("un_lock", unlock).Error
}

func (r *courseRepo) FindUserCourseStatus(ctx context.Context, uid, courseID int) (*UserCourseStatus, error) {
	return r.FindUserCourse(ctx, uid, courseID)
}

func (r *courseRepo) FindUserCourse(ctx context.Context, uid, courseID int) (*UserCourseStatus, error) {
	var row UserCourseStatus
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", uid, courseID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) UpdateUserCourseStatus(ctx context.Context, id, hasEnd, isEnd int) error {
	return r.db.WithContext(ctx).Table("user_course").Where("id = ?", id).
		UpdateColumns(map[string]any{"has_end_chapter": hasEnd, "is_end": isEnd}).Error
}

func (r *courseRepo) CreditGranted(ctx context.Context, uid int, field, remark string) (bool, error) {
	if uid <= 0 || remark == "" {
		return false, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Table("user_credit_records").
		Where("user_id = ? AND field = ? AND remark = ?", uid, field, remark).
		Count(&n).Error
	return n > 0, err
}

func (r *courseRepo) ListPathIDsByCourse(ctx context.Context, courseID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("path_courses").
		Where("course_id = ?", courseID).
		Pluck("path_id", &ids).Error
	return ids, err
}

func (r *courseRepo) FindUserPathStatus(ctx context.Context, uid, pathID int) (*UserPathStatus, error) {
	var row UserPathStatus
	err := r.db.WithContext(ctx).Table("user_path").
		Where("user_id = ? AND path_id = ?", uid, pathID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *courseRepo) ListCourseIDsByPath(ctx context.Context, pathID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("path_courses").
		Where("path_id = ? AND deleted_at IS NULL", pathID).
		Order("sort ASC").
		Pluck("course_id", &ids).Error
	return ids, err
}

func (r *courseRepo) CountEndedUserCourses(ctx context.Context, uid int, courseIDs []int) (int, error) {
	if len(courseIDs) == 0 {
		return 0, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id IN ? AND is_end = 1", uid, courseIDs).
		Count(&n).Error
	return int(n), err
}

func (r *courseRepo) UpdateUserPathStatus(ctx context.Context, id, hasEnd, isEnd int) error {
	return r.db.WithContext(ctx).Table("user_path").Where("id = ?", id).
		UpdateColumns(map[string]any{"has_end_course": hasEnd, "is_end": isEnd}).Error
}

func (r *courseRepo) ListUserIDsOnPathAmong(ctx context.Context, pathID int, userIDs []int) ([]int, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var ids []int
	err := r.db.WithContext(ctx).Table("user_path").
		Where("path_id = ? AND user_id IN ?", pathID, userIDs).
		Pluck("user_id", &ids).Error
	return ids, err
}

func (r *courseRepo) CreateUserPathsIgnore(ctx context.Context, rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).
		Table("user_path").Create(&rows).Error
}

// NewUserPathRow builds a user_path insert map.
func NewUserPathRow(uid, pathID int, now time.Time) map[string]any {
	return map[string]any{
		"user_id": uid, "path_id": pathID,
		"created_at": now, "updated_at": now,
	}
}

// NewUserCourseRow builds a user_course insert map.
func NewUserCourseRow(uid, courseID int, now time.Time) map[string]any {
	return map[string]any{
		"user_id": uid, "course_id": courseID,
		"created_at": now, "updated_at": now,
	}
}

// NewUserChapterRow builds a user_course_chapter insert map.
func NewUserChapterRow(uid, courseID, chapterID int, now time.Time) map[string]any {
	return map[string]any{
		"user_id": uid, "course_id": courseID, "chapter_id": chapterID,
		"un_lock": 0, "created_at": now, "updated_at": now,
	}
}
