package repo

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrAlreadyJoined = errors.New("already joined")

// LearnRepository covers learner study-progress reads/writes.
type LearnRepository interface {
	ListMyCourses(ctx context.Context, uid, typ, page, limit int) ([]MyCourseRow, error)
	ListContinueChapters(ctx context.Context, uid int, courseIDs []int) ([]ContinueChapterRow, error)

	ListMyTasks(ctx context.Context, uid, limit int) ([]MyTaskRow, error)

	CountHistory(ctx context.Context, uid int) (int64, error)
	ListHistory(ctx context.Context, uid, page, limit int) ([]HistoryRow, error)

	CountSoonExpire(ctx context.Context, uid int, deadline string) (int64, error)
	ListSoonExpire(ctx context.Context, uid, page, limit int, deadline string) ([]SoonExpireRow, error)

	CountEndedCourses(ctx context.Context, uid int) (int64, error)
	ListChapterProgress(ctx context.Context, uid int) ([]ChapterProgRow, error)
	SumMediaDuration(ctx context.Context, chapterIDs []int) (int64, error)

	FindLikeCourseID(ctx context.Context, uid, courseID int) (int, error)
	DeleteLikeCourse(ctx context.Context, id int) error
	CreateLikeCourse(ctx context.Context, uid, courseID int) error

	FindChapterCourseID(ctx context.Context, chapterID int) (courseID int, err error)
	FindLikeChapterID(ctx context.Context, uid, chapterID int) (int, error)
	DeleteLikeChapter(ctx context.Context, id int) error
	CreateLikeChapter(ctx context.Context, uid, chapterID, courseID int) error

	CountUserCourse(ctx context.Context, uid, courseID int) (int64, error)
	FindCourseMinCredit(ctx context.Context, courseID int) (minNeed int, err error)
	CountCourseChapters(ctx context.Context, courseID int) (int64, error)
	JoinCourse(ctx context.Context, uid, courseID int) error

	CourseExists(ctx context.Context, courseID int) (bool, error)
	CountPendingCourseApply(ctx context.Context, uid, courseID int) (int64, error)
	CreateCourseApply(ctx context.Context, uid, courseID int, reason string) error

	FindChapterUnlock(ctx context.Context, chapterID int) (*ChapterForUnlock, error)
	FindUserChapter(ctx context.Context, uid, chapterID int) (*LearnUserChapter, error)
	FindPrevChapterID(ctx context.Context, courseID, sort, id int) (int, error)
	FindUserChapterIsEnd(ctx context.Context, uid, chapterID int) (int, error)
	UnlockUserChapter(ctx context.Context, id int, expire time.Time) error

	FindChapterLearn(ctx context.Context, chapterID int) (*ChapterLearnRow, error)
	FindUserChapterFull(ctx context.Context, uid, chapterID int) (*UserChapterFullRow, error)
	LatestPaperPass(ctx context.Context, uid, paperID int) (int, error)
	UpdateUserChapterEnd(ctx context.Context, id, hasEnd, isEnd int) error
	CountEndedChapters(ctx context.Context, uid, courseID int) (int64, error)
	UpdateUserCourseHasEnd(ctx context.Context, uid, courseID, hasEnd int) error

	FindUserCourse(ctx context.Context, uid, courseID int) (*UserCourseEndRow, error)
	FindCourseEnd(ctx context.Context, courseID int) (*CourseEndRow, error)
	MarkUserCourseEnded(ctx context.Context, id int) error
}

type MyCourseRow struct {
	ID            int     `gorm:"column:id"`
	CourseID      int     `gorm:"column:course_id"`
	HasEndChapter int     `gorm:"column:has_end_chapter"`
	IsEnd         int     `gorm:"column:is_end"`
	CID           *int    `gorm:"column:cid"`
	Title         *string `gorm:"column:title"`
	Thumb         *string `gorm:"column:thumb"`
	ChapterCount  *int    `gorm:"column:chapter_count"`
}

type ContinueChapterRow struct {
	CourseID  int    `gorm:"column:course_id"`
	ChapterID int    `gorm:"column:chapter_id"`
	UnLock    int    `gorm:"column:un_lock"`
	IsEnd     int    `gorm:"column:is_end"`
	Title     string `gorm:"column:title"`
}

type MyTaskRow struct {
	CourseID     int     `gorm:"column:course_id"`
	ChapterID    int     `gorm:"column:chapter_id"`
	PaperID      *int    `gorm:"column:paper_id"`
	Title        *string `gorm:"column:title"`
	ChapterTitle *string `gorm:"column:chapter_title"`
}

type HistoryRow struct {
	ID           int       `gorm:"column:id"`
	ChapterID    int       `gorm:"column:chapter_id"`
	CourseID     int       `gorm:"column:course_id"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	ChapterTitle *string   `gorm:"column:chapter_title"`
	CourseTitle  *string   `gorm:"column:course_title"`
	CourseThumb  *string   `gorm:"column:course_thumb"`
}

type SoonExpireRow struct {
	ID           int        `gorm:"column:id"`
	UserID       int        `gorm:"column:user_id"`
	CourseID     int        `gorm:"column:course_id"`
	ChapterID    int        `gorm:"column:chapter_id"`
	IsEnd        int        `gorm:"column:is_end"`
	ExpireTime   *time.Time `gorm:"column:expire_time"`
	CourseTitle  *string    `gorm:"column:course_title"`
	CourseThumb  *string    `gorm:"column:course_thumb"`
	ChapterTitle *string    `gorm:"column:chapter_title"`
	ChapterSort  *int       `gorm:"column:chapter_sort"`
}

type ChapterProgRow struct {
	ChapterID     int        `gorm:"column:chapter_id"`
	IsEnd         int        `gorm:"column:is_end"`
	HasEndChapter int        `gorm:"column:has_end_chapter"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

type ChapterForUnlock struct {
	ID       int  `gorm:"column:id"`
	CourseID int  `gorm:"column:course_id"`
	Sort     int  `gorm:"column:sort"`
	LimitDay *int `gorm:"column:limit_day"`
}

type LearnUserChapter struct {
	ID     int `gorm:"column:id"`
	IsEnd  int `gorm:"column:is_end"`
	UnLock int `gorm:"column:un_lock"`
}

type ChapterLearnRow struct {
	ID       int    `gorm:"column:id"`
	CourseID int    `gorm:"column:course_id"`
	Type     int    `gorm:"column:type"`
	Title    string `gorm:"column:title"`
	PaperID  *int   `gorm:"column:paper_id"`
}

type UserChapterFullRow struct {
	ID            int `gorm:"column:id"`
	IsEnd         int `gorm:"column:is_end"`
	HasEndChapter int `gorm:"column:has_end_chapter"`
	UnLock        int `gorm:"column:un_lock"`
	UserID        int `gorm:"column:user_id"`
	CourseID      int `gorm:"column:course_id"`
	ChapterID     int `gorm:"column:chapter_id"`
	LearnTime     int `gorm:"column:learn_time"`
	ExpireTime    any `gorm:"column:expire_time"`
	CreatedAt     any `gorm:"column:created_at"`
	UpdatedAt     any `gorm:"column:updated_at"`
}

type UserCourseEndRow struct {
	ID    int `gorm:"column:id"`
	IsEnd int `gorm:"column:is_end"`
}

type CourseEndRow struct {
	ID      int    `gorm:"column:id"`
	Title   string `gorm:"column:title"`
	Credit  int    `gorm:"column:credit"`
	PaperID *int   `gorm:"column:paper_id"`
}

type learnRepo struct{ db *gorm.DB }

func NewLearn(db *gorm.DB) LearnRepository {
	return &learnRepo{db: db}
}

func (r *learnRepo) ListMyCourses(ctx context.Context, uid, typ, page, limit int) ([]MyCourseRow, error) {
	query := r.db.WithContext(ctx).Table("user_course AS uc").
		Select("uc.id, uc.course_id, uc.has_end_chapter, uc.is_end, uc.updated_at, c.id AS cid, c.title, c.thumb, c.chapter_count").
		Joins("LEFT JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL").
		Where("uc.user_id = ?", uid)
	switch typ {
	case 2:
		query = query.Order("uc.created_at DESC")
	case 3:
		query = query.Where("uc.is_end = 1").Order("uc.updated_at DESC")
	default:
		query = query.Order("uc.updated_at DESC")
	}
	var rows []MyCourseRow
	err := query.Offset((page - 1) * limit).Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *learnRepo) ListContinueChapters(ctx context.Context, uid int, courseIDs []int) ([]ContinueChapterRow, error) {
	if len(courseIDs) == 0 {
		return nil, nil
	}
	var rows []ContinueChapterRow
	err := r.db.WithContext(ctx).Table("user_course_chapter AS ucc").
		Select("ucc.course_id, ucc.chapter_id, ucc.un_lock, ucc.is_end, cc.title").
		Joins("JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL").
		Where("ucc.user_id = ? AND ucc.course_id IN ?", uid, courseIDs).
		Order("cc.sort ASC, cc.id ASC").Find(&rows).Error
	return rows, err
}

func (r *learnRepo) ListMyTasks(ctx context.Context, uid, limit int) ([]MyTaskRow, error) {
	var rows []MyTaskRow
	err := r.db.WithContext(ctx).Table("user_course_chapter AS ucc").
		Select("ucc.course_id, ucc.chapter_id, cc.paper_id, c.title, cc.title AS chapter_title").
		Joins("LEFT JOIN courses c ON c.id = ucc.course_id").
		Joins("LEFT JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL").
		Where("ucc.user_id = ? AND ucc.is_end = 0 AND ucc.has_end_chapter = 1", uid).
		Order("ucc.updated_at DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *learnRepo) CountHistory(ctx context.Context, uid int) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Table("user_course_chapter AS ucc").
		Where("ucc.user_id = ? AND ucc.has_end_chapter = 1", uid).Count(&total).Error
	return total, err
}

func (r *learnRepo) ListHistory(ctx context.Context, uid, page, limit int) ([]HistoryRow, error) {
	var rows []HistoryRow
	err := r.db.WithContext(ctx).Table("user_course_chapter AS ucc").
		Select("ucc.id, ucc.chapter_id, ucc.course_id, ucc.updated_at, cc.title AS chapter_title, c.title AS course_title, c.thumb AS course_thumb").
		Joins("LEFT JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL").
		Joins("LEFT JOIN courses c ON c.id = ucc.course_id AND c.deleted_at IS NULL").
		Where("ucc.user_id = ? AND ucc.has_end_chapter = 1", uid).
		Order("ucc.updated_at DESC, ucc.id DESC").
		Offset((page - 1) * limit).Limit(limit + 1).
		Find(&rows).Error
	return rows, err
}

func (r *learnRepo) CountSoonExpire(ctx context.Context, uid int, deadline string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Table("user_course_chapter AS ucc").
		Where("ucc.user_id = ? AND ucc.is_end = 0 AND ucc.un_lock = 1 AND ucc.expire_time IS NOT NULL AND ucc.expire_time < ?", uid, deadline).
		Count(&total).Error
	return total, err
}

func (r *learnRepo) ListSoonExpire(ctx context.Context, uid, page, limit int, deadline string) ([]SoonExpireRow, error) {
	var rows []SoonExpireRow
	err := r.db.WithContext(ctx).Table("user_course_chapter AS ucc").
		Select("ucc.id, ucc.user_id, ucc.course_id, ucc.chapter_id, ucc.is_end, ucc.expire_time, c.title AS course_title, c.thumb AS course_thumb, cc.title AS chapter_title, cc.sort AS chapter_sort").
		Joins("LEFT JOIN courses c ON c.id = ucc.course_id").
		Joins("LEFT JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL").
		Where("ucc.user_id = ? AND ucc.is_end = 0 AND ucc.un_lock = 1 AND ucc.expire_time IS NOT NULL AND ucc.expire_time < ?", uid, deadline).
		Order("ucc.expire_time ASC, ucc.id ASC").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *learnRepo) CountEndedCourses(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND is_end = 1", uid).Count(&n).Error
	return n, err
}

func (r *learnRepo) ListChapterProgress(ctx context.Context, uid int) ([]ChapterProgRow, error) {
	var rows []ChapterProgRow
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Select("chapter_id, is_end, has_end_chapter, updated_at").
		Where("user_id = ?", uid).Find(&rows).Error
	return rows, err
}

func (r *learnRepo) SumMediaDuration(ctx context.Context, chapterIDs []int) (int64, error) {
	if len(chapterIDs) == 0 {
		return 0, nil
	}
	var sum *int64
	err := r.db.WithContext(ctx).Table("course_chapters AS cc").
		Select("SUM(mf.duration)").
		Joins("JOIN media_file mf ON mf.id = cc.media_id").
		Where("cc.id IN ?", chapterIDs).Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	if sum == nil {
		return 0, nil
	}
	return *sum, nil
}

func (r *learnRepo) FindLikeCourseID(ctx context.Context, uid, courseID int) (int, error) {
	var existing struct{ ID int }
	err := r.db.WithContext(ctx).Table("user_like_courses").
		Where("user_id = ? AND course_id = ?", uid, courseID).Take(&existing).Error
	if err != nil {
		return 0, err
	}
	return existing.ID, nil
}

func (r *learnRepo) DeleteLikeCourse(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM user_like_courses WHERE id = ?", id).Error
}

func (r *learnRepo) CreateLikeCourse(ctx context.Context, uid, courseID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).Table("user_like_courses").Create(map[string]any{
		"user_id": uid, "course_id": courseID, "created_at": now, "updated_at": now,
	}).Error
}

func (r *learnRepo) FindChapterCourseID(ctx context.Context, chapterID int) (int, error) {
	var chapter struct {
		ID       int `gorm:"column:id"`
		CourseID int `gorm:"column:course_id"`
	}
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, course_id").Where("id = ? AND deleted_at IS NULL", chapterID).Take(&chapter).Error
	if err != nil {
		return 0, err
	}
	return chapter.CourseID, nil
}

func (r *learnRepo) FindLikeChapterID(ctx context.Context, uid, chapterID int) (int, error) {
	var existing struct{ ID int }
	err := r.db.WithContext(ctx).Table("user_like_chapters").
		Where("user_id = ? AND chapter_id = ?", uid, chapterID).Take(&existing).Error
	if err != nil {
		return 0, err
	}
	return existing.ID, nil
}

func (r *learnRepo) DeleteLikeChapter(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM user_like_chapters WHERE id = ?", id).Error
}

func (r *learnRepo) CreateLikeChapter(ctx context.Context, uid, chapterID, courseID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).Table("user_like_chapters").Create(map[string]any{
		"user_id": uid, "chapter_id": chapterID, "course_id": courseID,
		"created_at": now, "updated_at": now,
	}).Error
}

func (r *learnRepo) CountUserCourse(ctx context.Context, uid, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", uid, courseID).Count(&n).Error
	return n, err
}

func (r *learnRepo) FindCourseMinCredit(ctx context.Context, courseID int) (int, error) {
	var course struct {
		ID            int `gorm:"column:id"`
		MinNeedCredit int `gorm:"column:min_need_credit"`
	}
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, min_need_credit").
		Where("id = ? AND deleted_at IS NULL", courseID).Take(&course).Error
	if err != nil {
		return 0, err
	}
	return course.MinNeedCredit, nil
}

func (r *learnRepo) CountCourseChapters(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("course_chapters").
		Where("course_id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n, err
}

func (r *learnRepo) JoinCourse(ctx context.Context, uid, courseID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var locked int64
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("user_course").
			Where("user_id = ? AND course_id = ?", uid, courseID).Count(&locked).Error; err != nil {
			return err
		}
		if locked > 0 {
			return ErrAlreadyJoined
		}
		now := time.Now()
		return tx.Table("user_course").Create(map[string]any{
			"user_id": uid, "course_id": courseID, "created_at": now, "updated_at": now,
		}).Error
	})
}

func (r *learnRepo) CourseExists(ctx context.Context, courseID int) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("courses").
		Where("id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n > 0, err
}

func (r *learnRepo) CountPendingCourseApply(ctx context.Context, uid, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course_apply").
		Where("user_id = ? AND course_id = ? AND status != 2", uid, courseID).
		Count(&n).Error
	return n, err
}

func (r *learnRepo) CreateCourseApply(ctx context.Context, uid, courseID int, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Table("user_course_apply").Create(map[string]any{
		"user_id": uid, "course_id": courseID, "status": 0, "reason": reason,
		"created_at": now, "updated_at": now,
	}).Error
}

func (r *learnRepo) FindChapterUnlock(ctx context.Context, chapterID int) (*ChapterForUnlock, error) {
	var chapter ChapterForUnlock
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, course_id, sort, limit_day").Where("id = ? AND deleted_at IS NULL", chapterID).Take(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (r *learnRepo) FindUserChapter(ctx context.Context, uid, chapterID int) (*LearnUserChapter, error) {
	var row LearnUserChapter
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND chapter_id = ?", uid, chapterID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *learnRepo) FindPrevChapterID(ctx context.Context, courseID, sort, id int) (int, error) {
	var pre struct{ ID int }
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id").
		Where("course_id = ? AND deleted_at IS NULL AND (sort < ? OR (sort = ? AND id < ?))", courseID, sort, sort, id).
		Order("sort DESC, id DESC").Take(&pre).Error
	if err != nil {
		return 0, err
	}
	return pre.ID, nil
}

func (r *learnRepo) FindUserChapterIsEnd(ctx context.Context, uid, chapterID int) (int, error) {
	var preUser struct{ IsEnd int }
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Select("is_end").Where("user_id = ? AND chapter_id = ?", uid, chapterID).Take(&preUser).Error
	if err != nil {
		return 0, err
	}
	return preUser.IsEnd, nil
}

func (r *learnRepo) UnlockUserChapter(ctx context.Context, id int, expire time.Time) error {
	return r.db.WithContext(ctx).Table("user_course_chapter").
		Where("id = ?", id).
		Updates(map[string]any{"expire_time": expire, "un_lock": 1}).Error
}

func (r *learnRepo) FindChapterLearn(ctx context.Context, chapterID int) (*ChapterLearnRow, error) {
	var chapter ChapterLearnRow
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, course_id, type, title, paper_id").Where("id = ? AND deleted_at IS NULL", chapterID).Take(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (r *learnRepo) FindUserChapterFull(ctx context.Context, uid, chapterID int) (*UserChapterFullRow, error) {
	var row UserChapterFullRow
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND chapter_id = ?", uid, chapterID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *learnRepo) LatestPaperPass(ctx context.Context, uid, paperID int) (int, error) {
	var pass int
	err := r.db.WithContext(ctx).Table("user_paper").
		Select("is_pass").Where("user_id = ? AND paper_id = ? AND deleted_at IS NULL", uid, paperID).
		Order("id DESC").Limit(1).Scan(&pass).Error
	return pass, err
}

func (r *learnRepo) UpdateUserChapterEnd(ctx context.Context, id, hasEnd, isEnd int) error {
	return r.db.WithContext(ctx).Table("user_course_chapter").
		Where("id = ?", id).
		Updates(map[string]any{
			"has_end_chapter": hasEnd,
			"is_end":          isEnd,
		}).Error
}

func (r *learnRepo) CountEndedChapters(ctx context.Context, uid, courseID int) (int64, error) {
	var ended int64
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Where("user_id = ? AND course_id = ? AND is_end = 1", uid, courseID).
		Count(&ended).Error
	return ended, err
}

func (r *learnRepo) UpdateUserCourseHasEnd(ctx context.Context, uid, courseID, hasEnd int) error {
	return r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", uid, courseID).
		Update("has_end_chapter", hasEnd).Error
}

func (r *learnRepo) FindUserCourse(ctx context.Context, uid, courseID int) (*UserCourseEndRow, error) {
	var row UserCourseEndRow
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", uid, courseID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *learnRepo) FindCourseEnd(ctx context.Context, courseID int) (*CourseEndRow, error) {
	var course CourseEndRow
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, credit, paper_id").
		Where("id = ? AND deleted_at IS NULL", courseID).Take(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *learnRepo) MarkUserCourseEnded(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Table("user_course").
		Where("id = ?", id).Update("is_end", 1).Error
}
