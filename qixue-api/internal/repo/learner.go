package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// LearnerRepository covers learner-portal reads/writes across related tables.
type LearnerRepository interface {
	CountUserCourses(ctx context.Context, uid int) (int64, error)
	CountUnreadMessages(ctx context.Context, uid int) (int64, error)

	ListMessages(ctx context.Context, uid, page, limit int) (total int64, rows []map[string]any, err error)
	MarkMessagesRead(ctx context.Context, uid, id int) (int64, error)

	CountLikeCourses(ctx context.Context, uid int) (int64, error)
	ListLikeCourses(ctx context.Context, uid, page, limit int) ([]LikeCourseRow, error)
	CountLikeChapters(ctx context.Context, uid int) (int64, error)
	ListLikeChapters(ctx context.Context, uid, page, limit int) ([]LikeChapterRow, error)

	CountComments(ctx context.Context, table string, uid int) (int64, error)
	ListComments(ctx context.Context, table string, uid, page, limit int) ([]CommentRow, error)
	ListCommentChildren(ctx context.Context, table string, pids []int) ([]map[string]any, error)
	ListUserBriefs(ctx context.Context, ids []int) ([]map[string]any, error)
	ListCourseBriefs(ctx context.Context, ids []int) ([]CourseBrief, error)
	ListChapterBriefs(ctx context.Context, ids []int) ([]ChapterBrief, error)

	CountUserPapers(ctx context.Context, uid int, paperType string) (int64, error)
	ListUserPapers(ctx context.Context, uid int, paperType string, page, limit int) ([]UserPaperRow, error)

	CountRecords(ctx context.Context, table string, uid int) (int64, error)
	ListRecords(ctx context.Context, table string, uid, page, limit int) ([]map[string]any, error)

	CountUserPaths(ctx context.Context, uid int) (int64, error)
	ListUserPaths(ctx context.Context, uid, page, limit int) ([]UserPathRow, error)
	ListValidPathCerts(ctx context.Context, uid int, pathIDs []int) ([]map[string]any, error)

	CountValidCerts(ctx context.Context, uid int) (int64, error)
	ListValidCerts(ctx context.Context, uid, page, limit int) ([]UserCertRow, error)
}

type LikeCourseRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	Title     *string    `gorm:"column:title"`
	Thumb     *string    `gorm:"column:thumb"`
	Desc      *string    `gorm:"column:desc"`
	CID       *int       `gorm:"column:cid"`
}

type LikeChapterRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	ChapterID int        `gorm:"column:chapter_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	ChTitle   *string    `gorm:"column:ch_title"`
	CTitle    *string    `gorm:"column:c_title"`
	CThumb    *string    `gorm:"column:c_thumb"`
	CDesc     *string    `gorm:"column:c_desc"`
	CID       *int       `gorm:"column:cid"`
	ChID      *int       `gorm:"column:chid"`
}

type CommentRow struct {
	ID        int        `gorm:"column:id"`
	PID       int        `gorm:"column:pid"`
	CourseID  int        `gorm:"column:course_id"`
	ChapterID int        `gorm:"column:chapter_id"`
	UserID    int        `gorm:"column:user_id"`
	Content   string     `gorm:"column:content"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type CourseBrief struct {
	ID    int     `gorm:"column:id"`
	Title string  `gorm:"column:title"`
	Thumb string  `gorm:"column:thumb"`
	Desc  *string `gorm:"column:desc"`
}

type ChapterBrief struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
}

type UserPaperRow struct {
	ID            int        `gorm:"column:id"`
	UserID        int        `gorm:"column:user_id"`
	PaperID       int        `gorm:"column:paper_id"`
	Score         int        `gorm:"column:score"`
	IsPass        int        `gorm:"column:is_pass"`
	IsPeopleCheck int        `gorm:"column:is_people_check"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	PType         string     `gorm:"column:ptype"`
	CourseID      int        `gorm:"column:course_id"`
	ChapterID     *int       `gorm:"column:chapter_id"`
	CourseTitle   *string    `gorm:"column:course_title"`
	CourseThumb   *string    `gorm:"column:course_thumb"`
	ChapterTitle  *string    `gorm:"column:chapter_title"`
	CID           *int       `gorm:"column:cid"`
	ChID          *int       `gorm:"column:chid"`
}

type UserPathRow struct {
	ID           int        `gorm:"column:id"`
	UserID       int        `gorm:"column:user_id"`
	PathID       int        `gorm:"column:path_id"`
	HasEndCourse int        `gorm:"column:has_end_course"`
	IsEnd        int        `gorm:"column:is_end"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	PTitle       *string    `gorm:"column:ptitle"`
	PThumb       *string    `gorm:"column:pthumb"`
	PDesc        *string    `gorm:"column:pdesc"`
	PCertID      *int       `gorm:"column:pcert_id"`
	PID          *int       `gorm:"column:pid"`
}

type UserCertRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CertID    int        `gorm:"column:cert_id"`
	PathID    int        `gorm:"column:path_id"`
	CourseID  int        `gorm:"column:course_id"`
	CertNo    string     `gorm:"column:cert_no"`
	Title     string     `gorm:"column:title"`
	Image     string     `gorm:"column:image"`
	IsValid   int        `gorm:"column:is_valid"`
	IssuedAt  *time.Time `gorm:"column:issued_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	PTitle    *string    `gorm:"column:ptitle"`
	CTitle    *string    `gorm:"column:ctitle"`
	PID       *int       `gorm:"column:pid"`
	CID       *int       `gorm:"column:cid"`
}

type learnerRepo struct{ db *gorm.DB }

func NewLearner(db *gorm.DB) LearnerRepository {
	return &learnerRepo{db: db}
}

func (r *learnerRepo) CountUserCourses(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").Where("user_id = ?", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) CountUnreadMessages(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_messages").Where("user_id = ? AND is_read = 0", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListMessages(ctx context.Context, uid, page, limit int) (int64, []map[string]any, error) {
	q := r.db.WithContext(ctx).Table("user_messages").Where("user_id = ?", uid)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []map[string]any
	err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	if rows == nil {
		rows = []map[string]any{}
	}
	return total, rows, err
}

func (r *learnerRepo) MarkMessagesRead(ctx context.Context, uid, id int) (int64, error) {
	q := r.db.WithContext(ctx).Table("user_messages").Where("user_id = ? AND is_read = 0", uid)
	if id > 0 {
		q = q.Where("id = ?", id)
	}
	now := time.Now()
	res := q.Updates(map[string]any{"is_read": 1, "read_at": now, "updated_at": now})
	return res.RowsAffected, res.Error
}

func (r *learnerRepo) CountLikeCourses(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_like_courses").Where("user_id = ? AND deleted_at IS NULL", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListLikeCourses(ctx context.Context, uid, page, limit int) ([]LikeCourseRow, error) {
	var rows []LikeCourseRow
	err := r.db.WithContext(ctx).Table("user_like_courses AS ul").
		Select("ul.*, c.id AS cid, c.title, c.thumb, c.desc").
		Joins("LEFT JOIN courses c ON c.id = ul.course_id AND c.deleted_at IS NULL").
		Where("ul.user_id = ? AND ul.deleted_at IS NULL", uid).
		Order("ul.id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) CountLikeChapters(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_like_chapters").Where("user_id = ? AND deleted_at IS NULL", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListLikeChapters(ctx context.Context, uid, page, limit int) ([]LikeChapterRow, error) {
	var rows []LikeChapterRow
	err := r.db.WithContext(ctx).Table("user_like_chapters AS ul").
		Select("ul.*, cc.id AS chid, cc.title AS ch_title, c.id AS cid, c.title AS c_title, c.thumb AS c_thumb, c.desc AS c_desc").
		Joins("LEFT JOIN course_chapters cc ON cc.id = ul.chapter_id AND cc.deleted_at IS NULL").
		Joins("LEFT JOIN courses c ON c.id = ul.course_id AND c.deleted_at IS NULL").
		Where("ul.user_id = ? AND ul.deleted_at IS NULL", uid).
		Order("ul.id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) CountComments(ctx context.Context, table string, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table(table).Where("user_id = ? AND deleted_at IS NULL", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListComments(ctx context.Context, table string, uid, page, limit int) ([]CommentRow, error) {
	var rows []CommentRow
	err := r.db.WithContext(ctx).Table(table).
		Where("user_id = ? AND deleted_at IS NULL", uid).
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) ListCommentChildren(ctx context.Context, table string, pids []int) ([]map[string]any, error) {
	if len(pids) == 0 {
		return nil, nil
	}
	var rows []map[string]any
	err := r.db.WithContext(ctx).Table(table).
		Where("pid IN ? AND deleted_at IS NULL", pids).
		Order("id DESC").Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) ListUserBriefs(ctx context.Context, ids []int) ([]map[string]any, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []map[string]any
	err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, avatar, coins").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) ListCourseBriefs(ctx context.Context, ids []int) ([]CourseBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []CourseBrief
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb, `desc`").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) ListChapterBriefs(ctx context.Context, ids []int) ([]ChapterBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []ChapterBrief
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, title").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) CountUserPapers(ctx context.Context, uid int, paperType string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_paper AS up").
		Joins("JOIN papers p ON p.id = up.paper_id AND p.deleted_at IS NULL").
		Where("up.user_id = ? AND p.type = ?", uid, paperType).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListUserPapers(ctx context.Context, uid int, paperType string, page, limit int) ([]UserPaperRow, error) {
	var rows []UserPaperRow
	err := r.db.WithContext(ctx).Table("user_paper AS up").
		Select(`up.*, p.type AS ptype, p.course_id, p.chapter_id,
			c.id AS cid, c.title AS course_title, c.thumb AS course_thumb,
			cc.id AS chid, cc.title AS chapter_title`).
		Joins("JOIN papers p ON p.id = up.paper_id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN courses c ON c.id = p.course_id AND c.deleted_at IS NULL").
		Joins("LEFT JOIN course_chapters cc ON cc.id = p.chapter_id AND cc.deleted_at IS NULL").
		Where("up.user_id = ? AND p.type = ?", uid, paperType).
		Order("up.id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) CountRecords(ctx context.Context, table string, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table(table).Where("user_id = ? AND deleted_at IS NULL", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListRecords(ctx context.Context, table string, uid, page, limit int) ([]map[string]any, error) {
	var rows []map[string]any
	err := r.db.WithContext(ctx).Table(table).
		Where("user_id = ? AND deleted_at IS NULL", uid).
		Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	if rows == nil {
		rows = []map[string]any{}
	}
	return rows, err
}

func (r *learnerRepo) CountUserPaths(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path").Where("user_id = ? AND deleted_at IS NULL", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListUserPaths(ctx context.Context, uid, page, limit int) ([]UserPathRow, error) {
	var rows []UserPathRow
	err := r.db.WithContext(ctx).Table("user_path AS up").
		Select("up.*, p.id AS pid, p.title AS ptitle, p.thumb AS pthumb, p.desc AS pdesc, p.cert_id AS pcert_id").
		Joins("LEFT JOIN paths p ON p.id = up.path_id AND p.deleted_at IS NULL").
		Where("up.user_id = ? AND up.deleted_at IS NULL", uid).
		Order("up.id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) ListValidPathCerts(ctx context.Context, uid int, pathIDs []int) ([]map[string]any, error) {
	if len(pathIDs) == 0 {
		return nil, nil
	}
	var rows []map[string]any
	err := r.db.WithContext(ctx).Table("user_cert").
		Where("user_id = ? AND path_id IN ? AND is_valid = 1", uid, pathIDs).
		Find(&rows).Error
	return rows, err
}

func (r *learnerRepo) CountValidCerts(ctx context.Context, uid int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_cert").Where("user_id = ? AND is_valid = 1", uid).Count(&n).Error
	return n, err
}

func (r *learnerRepo) ListValidCerts(ctx context.Context, uid, page, limit int) ([]UserCertRow, error) {
	var rows []UserCertRow
	err := r.db.WithContext(ctx).Table("user_cert AS uc").
		Select("uc.id, uc.user_id, uc.cert_id, uc.path_id, uc.course_id, uc.cert_no, uc.title, uc.image, uc.is_valid, uc.issued_at, uc.created_at, p.id AS pid, p.title AS ptitle, c.id AS cid, c.title AS ctitle").
		Joins("LEFT JOIN paths p ON p.id = uc.path_id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL").
		Where("uc.user_id = ? AND uc.is_valid = 1", uid).
		Order("uc.id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, err
}
