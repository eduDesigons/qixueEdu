package admin

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IDName struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

type AdminNick struct {
	ID       int     `gorm:"column:id"`
	NickName *string `gorm:"column:nickName"`
}

type PaperScoreBrief struct {
	ID        int `gorm:"column:id"`
	Score     int `gorm:"column:score"`
	PassScore int `gorm:"column:pass_score"`
}

type CertDetail struct {
	ID            int     `gorm:"column:id"`
	Name          string  `gorm:"column:name"`
	PreviewImage  *string `gorm:"column:preview_image"`
	TemplateImage string  `gorm:"column:template_image"`
}

type CountByInt struct {
	Key   int `gorm:"column:key"`
	Total int `gorm:"column:total"`
}

type UserBriefFull struct {
	ID       int    `gorm:"column:id"`
	NickName string `gorm:"column:nick_name"`
	Credit   int    `gorm:"column:credit"`
	Avatar   string `gorm:"column:avatar"`
	Mobile   string `gorm:"column:mobile"`
}

type UserSearchBrief struct {
	ID       int    `gorm:"column:id"`
	NickName string `gorm:"column:nick_name"`
	Credit   int    `gorm:"column:credit"`
	Avatar   string `gorm:"column:avatar"`
}

type CourseThumbBrief struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
	Thumb string `gorm:"column:thumb"`
}

type UserCourseRow struct {
	ID            int        `gorm:"column:id"`
	UserID        int        `gorm:"column:user_id"`
	CourseID      int        `gorm:"column:course_id"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	HasEndChapter int        `gorm:"column:has_end_chapter"`
	IsEnd         int        `gorm:"column:is_end"`
}

type AdminCourseUserCertRow struct {
	ID        int64      `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CertID    int64      `gorm:"column:cert_id"`
	CertNo    string     `gorm:"column:cert_no"`
	URL       *string    `gorm:"column:url"`
	PathID    *int       `gorm:"column:path_id"`
	CourseID  int        `gorm:"column:course_id"`
	Title     string     `gorm:"column:title"`
	Image     string     `gorm:"column:image"`
	IsValid   int        `gorm:"column:is_valid"`
	IssuedAt  *time.Time `gorm:"column:issued_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type CourseApplyRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	Reason    *string    `gorm:"column:reason"`
	Status    int        `gorm:"column:status"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type CourseUserListFilter struct {
	CourseID int
	CityID   int
	NickName string
	Mobile   string
	Offset   int
	Limit    int
}

type CourseApplyListFilter struct {
	CourseID    *int
	CourseTitle string
	NickName    string
	UserID      *int
	Status      *int
	StartTime   string
	EndTime     string
	CityID      int
	Offset      int
	Limit       int
}

type AdminCourseRepository interface {
	ListCategoriesByIDs(ctx context.Context, ids []int) ([]IDName, error)
	ListCertNamesByIDs(ctx context.Context, ids []int) ([]IDName, error)
	ChapterCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	UserCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	FindCategory(ctx context.Context, id int) (*IDName, error)
	FindAdmin(ctx context.Context, id int) (*AdminNick, error)
	FindPaperScore(ctx context.Context, id int) (*PaperScoreBrief, error)
	FindCertDetail(ctx context.Context, id int) (*CertDetail, error)
	CountChapters(ctx context.Context, courseID int) (int64, error)
	CountComments(ctx context.Context, courseID int) (int64, error)
	CountUsers(ctx context.Context, courseID int) (int64, error)
	CountEndedUsers(ctx context.Context, courseID int) (int64, error)
	CountLikes(ctx context.Context, courseID int) (int64, error)
	PluckChapterIDs(ctx context.Context, courseID int) ([]int, error)
	PluckChapterPaperIDs(ctx context.Context, courseID int) ([]int, error)
	DeleteCourseCascade(ctx context.Context, courseID int, paperID *int, chapterIDs, chapterPaperIDs []int) error
	ListCourseUsers(ctx context.Context, f CourseUserListFilter) ([]UserCourseRow, int64, error)
	ListUsersFull(ctx context.Context, ids []int) ([]UserBriefFull, error)
	ListCourseUserCerts(ctx context.Context, courseID int, userIDs []int) ([]AdminCourseUserCertRow, error)
	SearchUsersByNick(ctx context.Context, nick string, cityID, limit int) ([]UserSearchBrief, error)
	CountEnrolled(ctx context.Context, courseID int, userIDs []int) (int64, error)
	FindUserCourse(ctx context.Context, id int) (userID, courseID int, err error)
	DeleteCourseUserCascade(ctx context.Context, rowID, userID, courseID int) error
	ListCourseApplies(ctx context.Context, f CourseApplyListFilter) ([]CourseApplyRow, int64, error)
	ListCourseThumbs(ctx context.Context, ids []int) ([]CourseThumbBrief, error)
	FindCourseApply(ctx context.Context, id int) (*CourseApplyRow, error)
	CourseTitle(ctx context.Context, id int) (string, error)
	UpdateCourseApplyStatus(ctx context.Context, id, status int) error
	PluckPendingApplyCourseIDs(ctx context.Context, ids []int) ([]int, error)
	CountCourseApply(ctx context.Context, id int) (int64, error)
	DeleteCourseApply(ctx context.Context, id int) error
	SendUserMessage(ctx context.Context, userID int, title, content, relatedType string, relatedID int) error
	ChaptersExist(ctx context.Context, courseID int) (bool, error)
	CountPendingCertUsers(ctx context.Context, courseID int) (int, error)
	CountIssuedCertUsers(ctx context.Context, courseID int) (int, error)
}

type adminCourseRepo struct{ db *gorm.DB }

func NewAdminCourse(db *gorm.DB) AdminCourseRepository { return &adminCourseRepo{db: db} }

func (r *adminCourseRepo) ListCategoriesByIDs(ctx context.Context, ids []int) ([]IDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []IDName
	err := r.db.WithContext(ctx).Table("course_categories").Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminCourseRepo) ListCertNamesByIDs(ctx context.Context, ids []int) ([]IDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []IDName
	err := r.db.WithContext(ctx).Table("certs").Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminCourseRepo) ChapterCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		CourseID int `gorm:"column:course_id"`
		Total    int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("course_id, COUNT(*) AS total").Where("course_id IN ? AND deleted_at IS NULL", ids).
		Group("course_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.CourseID, Total: row.Total})
	}
	return out, err
}

func (r *adminCourseRepo) UserCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		CourseID int `gorm:"column:course_id"`
		Total    int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("user_course").
		Select("course_id, COUNT(*) AS total").Where("course_id IN ?", ids).
		Group("course_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.CourseID, Total: row.Total})
	}
	return out, err
}

func (r *adminCourseRepo) FindCategory(ctx context.Context, id int) (*IDName, error) {
	var row IDName
	err := r.db.WithContext(ctx).Table("course_categories").Select("id, name").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminCourseRepo) FindAdmin(ctx context.Context, id int) (*AdminNick, error) {
	var row AdminNick
	err := r.db.WithContext(ctx).Table("administrator").Select("id, nickName").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminCourseRepo) FindPaperScore(ctx context.Context, id int) (*PaperScoreBrief, error) {
	var row PaperScoreBrief
	err := r.db.WithContext(ctx).Table("papers").
		Select("id, score, pass_score").Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminCourseRepo) FindCertDetail(ctx context.Context, id int) (*CertDetail, error) {
	var row CertDetail
	err := r.db.WithContext(ctx).Table("certs").
		Select("id, name, preview_image, template_image").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminCourseRepo) CountChapters(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("course_chapters").Where("course_id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n, err
}

func (r *adminCourseRepo) CountComments(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("course_comments").Where("course_id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n, err
}

func (r *adminCourseRepo) CountUsers(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").Where("course_id = ?", courseID).Count(&n).Error
	return n, err
}

func (r *adminCourseRepo) CountEndedUsers(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").Where("course_id = ? AND is_end = 1", courseID).Count(&n).Error
	return n, err
}

func (r *adminCourseRepo) CountLikes(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_like_courses").Where("course_id = ?", courseID).Count(&n).Error
	return n, err
}

func (r *adminCourseRepo) PluckChapterIDs(ctx context.Context, courseID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("course_chapters").Where("course_id = ?", courseID).Pluck("id", &ids).Error
	return ids, err
}

func (r *adminCourseRepo) PluckChapterPaperIDs(ctx context.Context, courseID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("course_chapters").Where("course_id = ? AND paper_id > 0", courseID).Pluck("paper_id", &ids).Error
	return ids, err
}

func (r *adminCourseRepo) DeleteCourseCascade(ctx context.Context, courseID int, paperID *int, chapterIDs, chapterPaperIDs []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(chapterIDs) > 0 {
			now := time.Now()
			if err := tx.Table("course_chapter_comments").
				Where("chapter_id IN ? AND deleted_at IS NULL", chapterIDs).
				Update("deleted_at", now).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_like_chapters WHERE chapter_id IN ?", chapterIDs).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_course_chapter WHERE chapter_id IN ?", chapterIDs).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM course_chapters WHERE id IN ?", chapterIDs).Error; err != nil {
				return err
			}
		}
		paperIDs := uniquePositiveInts(chapterPaperIDs)
		if paperID != nil && *paperID > 0 {
			paperIDs = uniquePositiveInts(append(paperIDs, *paperID))
		}
		if len(paperIDs) > 0 {
			var userPaperIDs []int64
			if err := tx.Table("user_paper").Where("paper_id IN ?", paperIDs).Pluck("id", &userPaperIDs).Error; err != nil {
				return err
			}
			if len(userPaperIDs) > 0 {
				if err := tx.Exec("DELETE FROM user_paper_question_log WHERE user_paper_id IN ?", userPaperIDs).Error; err != nil {
					return err
				}
			}
			now := time.Now()
			if err := tx.Table("user_paper").Where("paper_id IN ? AND deleted_at IS NULL", paperIDs).
				Update("deleted_at", now).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec("DELETE FROM user_like_courses WHERE course_id = ?", courseID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_course WHERE course_id = ?", courseID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_course_apply WHERE course_id = ?", courseID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM path_courses WHERE course_id = ?", courseID).Error; err != nil {
			return err
		}
		now := time.Now()
		if err := tx.Table("course_comments").Where("course_id = ? AND deleted_at IS NULL", courseID).
			Update("deleted_at", now).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM course_attach WHERE course_id = ?", courseID).Error; err != nil {
			return err
		}
		return tx.Exec("UPDATE courses SET deleted_at = ? WHERE id = ?", now, courseID).Error
	})
}

func (r *adminCourseRepo) ListCourseUsers(ctx context.Context, f CourseUserListFilter) ([]UserCourseRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_course").Where("course_id = ?", f.CourseID)
	query = applyUserCityFilter(query, f.CityID)
	if f.NickName != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Mobile != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE mobile = ?)", f.Mobile)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []UserCourseRow
	err := query.Select("id, user_id, course_id, created_at, updated_at, has_end_chapter, is_end").
		Order("id ASC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminCourseRepo) ListUsersFull(ctx context.Context, ids []int) ([]UserBriefFull, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []UserBriefFull
	err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, credit, avatar, mobile").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminCourseRepo) ListCourseUserCerts(ctx context.Context, courseID int, userIDs []int) ([]AdminCourseUserCertRow, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var rows []AdminCourseUserCertRow
	err := r.db.WithContext(ctx).Table("user_cert").
		Where("course_id = ? AND path_id = 0 AND user_id IN ?", courseID, userIDs).Find(&rows).Error
	return rows, err
}

func (r *adminCourseRepo) SearchUsersByNick(ctx context.Context, nick string, cityID, limit int) ([]UserSearchBrief, error) {
	query := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, credit, avatar").
		Where("nick_name LIKE ?", "%"+nick+"%")
	if cityID > 0 {
		query = query.Where("city_id = ?", cityID)
	}
	var rows []UserSearchBrief
	err := query.Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *adminCourseRepo) CountEnrolled(ctx context.Context, courseID int, userIDs []int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").
		Where("course_id = ? AND user_id IN ?", courseID, userIDs).Count(&n).Error
	return n, err
}

func (r *adminCourseRepo) FindUserCourse(ctx context.Context, id int) (int, int, error) {
	var row struct {
		UserID   int `gorm:"column:user_id"`
		CourseID int `gorm:"column:course_id"`
	}
	err := r.db.WithContext(ctx).Table("user_course").Where("id = ?", id).Take(&row).Error
	return row.UserID, row.CourseID, err
}

func (r *adminCourseRepo) DeleteCourseUserCascade(ctx context.Context, rowID, userID, courseID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM user_course WHERE id = ?", rowID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_like_courses WHERE user_id = ? AND course_id = ?", userID, courseID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_course_chapter WHERE user_id = ? AND course_id = ?", userID, courseID).Error; err != nil {
			return err
		}
		now := time.Now()
		if err := tx.Table("course_comments").
			Where("user_id = ? AND course_id = ? AND deleted_at IS NULL", userID, courseID).
			Update("deleted_at", now).Error; err != nil {
			return err
		}
		if err := tx.Table("course_chapter_comments").
			Where("user_id = ? AND course_id = ? AND deleted_at IS NULL", userID, courseID).
			Update("deleted_at", now).Error; err != nil {
			return err
		}
		return tx.Exec("DELETE FROM user_course_apply WHERE user_id = ? AND course_id = ?", userID, courseID).Error
	})
}

func (r *adminCourseRepo) ListCourseApplies(ctx context.Context, f CourseApplyListFilter) ([]CourseApplyRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_course_apply")
	if f.CourseID != nil {
		query = query.Where("course_id = ?", *f.CourseID)
	}
	if f.CourseTitle != "" {
		query = query.Where("course_id IN (SELECT id FROM courses WHERE title LIKE ? AND deleted_at IS NULL)", "%"+f.CourseTitle+"%")
	}
	if f.NickName != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.UserID != nil {
		query = query.Where("user_id = ?", *f.UserID)
	}
	if f.Status != nil {
		query = query.Where("status = ?", *f.Status)
	}
	if f.StartTime != "" {
		query = query.Where("created_at >= ?", f.StartTime)
	}
	if f.EndTime != "" {
		query = query.Where("created_at <= ?", f.EndTime)
	}
	query = applyUserCityFilter(query, f.CityID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []CourseApplyRow
	err := query.Select("id, user_id, course_id, reason, status, created_at, updated_at").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminCourseRepo) ListCourseThumbs(ctx context.Context, ids []int) ([]CourseThumbBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []CourseThumbBrief
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminCourseRepo) FindCourseApply(ctx context.Context, id int) (*CourseApplyRow, error) {
	var row CourseApplyRow
	err := r.db.WithContext(ctx).Table("user_course_apply").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminCourseRepo) CourseTitle(ctx context.Context, id int) (string, error) {
	var title string
	err := r.db.WithContext(ctx).Table("courses").Select("title").Where("id = ? AND deleted_at IS NULL", id).Scan(&title).Error
	return title, err
}

func (r *adminCourseRepo) UpdateCourseApplyStatus(ctx context.Context, id, status int) error {
	return r.db.WithContext(ctx).Table("user_course_apply").Where("id = ?", id).Update("status", status).Error
}

func (r *adminCourseRepo) PluckPendingApplyCourseIDs(ctx context.Context, ids []int) ([]int, error) {
	var courseIDs []int
	err := r.db.WithContext(ctx).Table("user_course_apply").
		Where("id IN ? AND status = 0", ids).Pluck("course_id", &courseIDs).Error
	return courseIDs, err
}

func (r *adminCourseRepo) CountCourseApply(ctx context.Context, id int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course_apply").Where("id = ?", id).Count(&n).Error
	return n, err
}

func (r *adminCourseRepo) DeleteCourseApply(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM user_course_apply WHERE id = ?", id).Error
}

func (r *adminCourseRepo) SendUserMessage(ctx context.Context, userID int, title, content, relatedType string, relatedID int) error {
	if userID <= 0 || title == "" || content == "" {
		return nil
	}
	now := r.db.NowFunc()
	return r.db.WithContext(ctx).Table("user_messages").Create(map[string]any{
		"user_id": userID, "admin_id": 0, "title": title, "content": content,
		"type": "notice", "related_type": relatedType, "related_id": relatedID,
		"is_read": 0, "created_at": now, "updated_at": now,
	}).Error
}

func (r *adminCourseRepo) ChaptersExist(ctx context.Context, courseID int) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("course_chapters").Where("course_id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n > 0, err
}

func (r *adminCourseRepo) CountPendingCertUsers(ctx context.Context, courseID int) (int, error) {
	if courseID <= 0 {
		return 0, nil
	}
	var certID *int
	if err := r.db.WithContext(ctx).Table("courses").Select("cert_id").
		Where("id = ? AND deleted_at IS NULL", courseID).Scan(&certID).Error; err != nil {
		return 0, err
	}
	if certID == nil || *certID <= 0 {
		return 0, nil
	}
	var ended []int
	if err := r.db.WithContext(ctx).Table("user_course").
		Where("course_id = ? AND is_end = 1", courseID).
		Pluck("user_id", &ended).Error; err != nil {
		return 0, err
	}
	if len(ended) == 0 {
		return 0, nil
	}
	var hasCert []int
	if err := r.db.WithContext(ctx).Table("user_cert").
		Where("course_id = ? AND path_id = 0 AND is_valid = 1 AND user_id IN ?", courseID, ended).
		Pluck("user_id", &hasCert).Error; err != nil {
		return 0, err
	}
	set := map[int]struct{}{}
	for _, id := range hasCert {
		set[id] = struct{}{}
	}
	pending := 0
	for _, id := range ended {
		if _, ok := set[id]; !ok {
			pending++
		}
	}
	return pending, nil
}

func (r *adminCourseRepo) CountIssuedCertUsers(ctx context.Context, courseID int) (int, error) {
	if courseID <= 0 {
		return 0, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Table("user_cert").
		Where("course_id = ? AND path_id = 0", courseID).Count(&n).Error
	return int(n), err
}

func applyUserCityFilter(query *gorm.DB, cityID int) *gorm.DB {
	if cityID <= 0 {
		return query
	}
	return query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
}

func uniquePositiveInts(ids []int) []int {
	seen := map[int]struct{}{}
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
