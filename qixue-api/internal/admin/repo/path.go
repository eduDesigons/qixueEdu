package admin

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type PathCourseSearchFilter struct {
	Title    string
	Category int
	Type     int
	PathID   int
	Offset   int
	Limit    int
}

type PathCourseSearchRow struct {
	ID            int    `gorm:"column:id"`
	Title         string `gorm:"column:title"`
	Thumb         string `gorm:"column:thumb"`
	MinNeedCredit int    `gorm:"column:min_need_credit"`
	Credit        int    `gorm:"column:credit"`
	Type          int    `gorm:"column:type"`
	CategoryID    int    `gorm:"column:category_id"`
}

type PathCourseDetailRow struct {
	ID            int    `gorm:"column:id"`
	Title         string `gorm:"column:title"`
	Thumb         string `gorm:"column:thumb"`
	Type          int    `gorm:"column:type"`
	Desc          string `gorm:"column:desc"`
	IsShow        int    `gorm:"column:is_show"`
	CategoryID    int    `gorm:"column:category_id"`
	IsRec         int    `gorm:"column:is_rec"`
	MinNeedCredit int    `gorm:"column:min_need_credit"`
	Credit        int    `gorm:"column:credit"`
}

type AdminUserPathRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	PathID    int        `gorm:"column:path_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	IsEnd     int        `gorm:"column:is_end"`
}

type PathUserCertRow struct {
	ID        int64      `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CertID    int64      `gorm:"column:cert_id"`
	CertNo    string     `gorm:"column:cert_no"`
	URL       *string    `gorm:"column:url"`
	PathID    int        `gorm:"column:path_id"`
	CourseID  *int       `gorm:"column:course_id"`
	Title     string     `gorm:"column:title"`
	Image     string     `gorm:"column:image"`
	IsValid   int        `gorm:"column:is_valid"`
	IssuedAt  *time.Time `gorm:"column:issued_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type PathUserListFilter struct {
	PathID   int
	CityID   int
	NickName string
	Mobile   string
	Offset   int
	Limit    int
}

type PathApplyRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	PathID    int        `gorm:"column:path_id"`
	Reason    *string    `gorm:"column:reason"`
	Status    int        `gorm:"column:status"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type PathApplyListFilter struct {
	PathID    *int
	UserID    *int
	PathTitle string
	NickName  string
	Status    *int
	StartTime string
	EndTime   string
	CityID    int
	Offset    int
	Limit     int
}

type PathThumbBrief struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
	Thumb string `gorm:"column:thumb"`
}

type AdminPathRepository interface {
	ListPathCategories(ctx context.Context) ([]IDName, error)
	ListCategoriesByIDs(ctx context.Context, ids []int) ([]IDName, error)
	ListAdminsByIDs(ctx context.Context, ids []int) ([]AdminNick, error)
	ListCertNamesByIDs(ctx context.Context, ids []int) ([]IDName, error)
	CourseCountsByPathIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	UserCountsByPathIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	FindCategory(ctx context.Context, id int) (*IDName, error)
	FindAdmin(ctx context.Context, id int) (*AdminNick, error)
	FindCertDetail(ctx context.Context, id int) (*CertDetail, error)
	CountPathCourses(ctx context.Context, pathID int) (int64, error)
	CountPathUsers(ctx context.Context, pathID int) (int64, error)
	CountPathEndedUsers(ctx context.Context, pathID int) (int64, error)
	CountPathApplies(ctx context.Context, pathID int) (int64, error)
	CountPendingCertUsers(ctx context.Context, pathID int) (int, error)
	CountIssuedCertUsers(ctx context.Context, pathID int) (int, error)
	DeletePathCascade(ctx context.Context, pathID int) error
	SearchCourses(ctx context.Context, f PathCourseSearchFilter) ([]PathCourseSearchRow, int64, error)
	ListCourseCategories(ctx context.Context, ids []int) ([]IDName, error)
	ListCoursesDetail(ctx context.Context, ids []int) ([]PathCourseDetailRow, error)
	ChapterCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	UserCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	CountPathCoursesIn(ctx context.Context, pathID int, courseIDs []int) (int64, error)
	PluckExistingCourseIDs(ctx context.Context, ids []int) ([]int, error)
	CountPathCoursesActive(ctx context.Context, pathID int) (int64, error)
	CreatePathCourses(ctx context.Context, rows []map[string]any) error
	CountAnyPathUsers(ctx context.Context, pathID int) (int64, error)
	ListPathUsers(ctx context.Context, f PathUserListFilter) ([]AdminUserPathRow, int64, error)
	ListUsersFull(ctx context.Context, ids []int) ([]UserBriefFull, error)
	ListPathUserCerts(ctx context.Context, pathID int, userIDs []int) ([]PathUserCertRow, error)
	SearchUsersByNick(ctx context.Context, nick string, cityID, limit int) ([]UserSearchBrief, error)
	CountEnrolled(ctx context.Context, pathID int, userIDs []int) (int64, error)
	FindUserPath(ctx context.Context, id int) (userID, pathID int, err error)
	PluckPathCourseIDs(ctx context.Context, pathID int) ([]int, error)
	DeleteUserPathCascade(ctx context.Context, rowID, userID, pathID int, courseIDs []int) error
	ListPathApplies(ctx context.Context, f PathApplyListFilter) ([]PathApplyRow, int64, error)
	ListPathThumbs(ctx context.Context, ids []int) ([]PathThumbBrief, error)
	FindPathApply(ctx context.Context, id int) (*PathApplyRow, error)
	PathTitle(ctx context.Context, id int) (string, error)
	UpdatePathApplyStatus(ctx context.Context, id, status int) error
	CountPathApply(ctx context.Context, id int) (int64, error)
	DeletePathApply(ctx context.Context, id int) error
	SendUserMessage(ctx context.Context, userID int, title, content, relatedType string, relatedID int) error
}

type adminPathRepo struct{ db *gorm.DB }

func NewAdminPath(db *gorm.DB) AdminPathRepository { return &adminPathRepo{db: db} }

func (r *adminPathRepo) ListPathCategories(ctx context.Context) ([]IDName, error) {
	var rows []IDName
	err := r.db.WithContext(ctx).Table("path_categories").
		Select("id, name").Where("parent_id = 0").Order("sort ASC").Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) ListCategoriesByIDs(ctx context.Context, ids []int) ([]IDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []IDName
	err := r.db.WithContext(ctx).Table("path_categories").Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) ListAdminsByIDs(ctx context.Context, ids []int) ([]AdminNick, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminNick
	err := r.db.WithContext(ctx).Table("administrator").Select("id, nickName").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) ListCertNamesByIDs(ctx context.Context, ids []int) ([]IDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []IDName
	err := r.db.WithContext(ctx).Table("certs").Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) CourseCountsByPathIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		PathID int `gorm:"column:path_id"`
		Total  int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("path_courses").
		Select("path_id, COUNT(*) AS total").Where("path_id IN ?", ids).
		Group("path_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.PathID, Total: row.Total})
	}
	return out, err
}

func (r *adminPathRepo) UserCountsByPathIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		PathID int `gorm:"column:path_id"`
		Total  int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("user_path").
		Select("path_id, COUNT(*) AS total").Where("path_id IN ? AND deleted_at IS NULL", ids).
		Group("path_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.PathID, Total: row.Total})
	}
	return out, err
}

func (r *adminPathRepo) FindCategory(ctx context.Context, id int) (*IDName, error) {
	var row IDName
	err := r.db.WithContext(ctx).Table("path_categories").Select("id, name").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPathRepo) FindAdmin(ctx context.Context, id int) (*AdminNick, error) {
	var row AdminNick
	err := r.db.WithContext(ctx).Table("administrator").Select("id, nickName").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPathRepo) FindCertDetail(ctx context.Context, id int) (*CertDetail, error) {
	var row CertDetail
	err := r.db.WithContext(ctx).Table("certs").
		Select("id, name, preview_image, template_image").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPathRepo) CountPathCourses(ctx context.Context, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("path_courses").Where("path_id = ?", pathID).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) CountPathUsers(ctx context.Context, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path").Where("path_id = ? AND deleted_at IS NULL", pathID).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) CountPathEndedUsers(ctx context.Context, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path").Where("path_id = ? AND is_end = 1 AND deleted_at IS NULL", pathID).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) CountPathApplies(ctx context.Context, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path_apply").Where("path_id = ?", pathID).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) CountPendingCertUsers(ctx context.Context, pathID int) (int, error) {
	if pathID <= 0 {
		return 0, nil
	}
	var certID *int
	if err := r.db.WithContext(ctx).Table("paths").Select("cert_id").
		Where("id = ? AND deleted_at IS NULL", pathID).Scan(&certID).Error; err != nil {
		return 0, err
	}
	if certID == nil || *certID <= 0 {
		return 0, nil
	}
	var ended []int
	if err := r.db.WithContext(ctx).Table("user_path").
		Where("path_id = ? AND is_end = 1", pathID).
		Pluck("user_id", &ended).Error; err != nil {
		return 0, err
	}
	if len(ended) == 0 {
		return 0, nil
	}
	var hasCert []int
	if err := r.db.WithContext(ctx).Table("user_cert").
		Where("path_id = ? AND course_id = 0 AND is_valid = 1 AND user_id IN ?", pathID, ended).
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

func (r *adminPathRepo) CountIssuedCertUsers(ctx context.Context, pathID int) (int, error) {
	if pathID <= 0 {
		return 0, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Table("user_cert").
		Where("path_id = ? AND course_id = 0", pathID).Count(&n).Error
	return int(n), err
}

func (r *adminPathRepo) DeletePathCascade(ctx context.Context, pathID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM user_path_apply WHERE path_id = ?", pathID).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Exec("UPDATE paths SET deleted_at = ? WHERE id = ?", now, pathID).Error
	})
}

func (r *adminPathRepo) SearchCourses(ctx context.Context, f PathCourseSearchFilter) ([]PathCourseSearchRow, int64, error) {
	query := r.db.WithContext(ctx).Table("courses").Where("deleted_at IS NULL")
	if f.Title != "" {
		query = query.Where("title LIKE ?", "%"+f.Title+"%")
	}
	if f.Category > 0 {
		query = query.Where("category_id = ?", f.Category)
	}
	if f.Type > 0 {
		query = query.Where("type = ?", f.Type)
	}
	if f.PathID > 0 {
		query = query.Where("id NOT IN (SELECT course_id FROM path_courses WHERE path_id = ? AND deleted_at IS NULL)", f.PathID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []PathCourseSearchRow
	err := query.Select("id, title, thumb, min_need_credit, credit, type, category_id").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminPathRepo) ListCourseCategories(ctx context.Context, ids []int) ([]IDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []IDName
	err := r.db.WithContext(ctx).Table("course_categories").Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) ListCoursesDetail(ctx context.Context, ids []int) ([]PathCourseDetailRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []PathCourseDetailRow
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb, type, courses.desc, is_show, category_id, is_rec, min_need_credit, credit").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) ChapterCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		CourseID int `gorm:"column:course_id"`
		Total    int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("course_id, COUNT(*) AS total").Where("course_id IN ?", ids).
		Group("course_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.CourseID, Total: row.Total})
	}
	return out, err
}

func (r *adminPathRepo) UserCountsByCourseIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
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

func (r *adminPathRepo) CountPathCoursesIn(ctx context.Context, pathID int, courseIDs []int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("path_courses").
		Where("path_id = ? AND course_id IN ? AND deleted_at IS NULL", pathID, courseIDs).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) PluckExistingCourseIDs(ctx context.Context, ids []int) ([]int, error) {
	var exist []int
	err := r.db.WithContext(ctx).Table("courses").Where("id IN ? AND deleted_at IS NULL", ids).Pluck("id", &exist).Error
	return exist, err
}

func (r *adminPathRepo) CountPathCoursesActive(ctx context.Context, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("path_courses").
		Where("path_id = ? AND deleted_at IS NULL", pathID).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) CreatePathCourses(ctx context.Context, rows []map[string]any) error {
	return r.db.WithContext(ctx).Table("path_courses").Create(&rows).Error
}

func (r *adminPathRepo) CountAnyPathUsers(ctx context.Context, pathID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path").Where("path_id = ?", pathID).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) ListPathUsers(ctx context.Context, f PathUserListFilter) ([]AdminUserPathRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_path").Where("path_id = ? AND deleted_at IS NULL", f.PathID)
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
	var rows []AdminUserPathRow
	err := query.Select("id, user_id, path_id, created_at, updated_at, is_end").
		Order("id ASC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminPathRepo) ListUsersFull(ctx context.Context, ids []int) ([]UserBriefFull, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []UserBriefFull
	err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, credit, avatar, mobile").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) ListPathUserCerts(ctx context.Context, pathID int, userIDs []int) ([]PathUserCertRow, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var rows []PathUserCertRow
	err := r.db.WithContext(ctx).Table("user_cert").
		Where("path_id = ? AND user_id IN ?", pathID, userIDs).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) SearchUsersByNick(ctx context.Context, nick string, cityID, limit int) ([]UserSearchBrief, error) {
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

func (r *adminPathRepo) CountEnrolled(ctx context.Context, pathID int, userIDs []int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path").
		Where("path_id = ? AND user_id IN ? AND deleted_at IS NULL", pathID, userIDs).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) FindUserPath(ctx context.Context, id int) (int, int, error) {
	var row struct {
		UserID int `gorm:"column:user_id"`
		PathID int `gorm:"column:path_id"`
	}
	err := r.db.WithContext(ctx).Table("user_path").Where("id = ?", id).Take(&row).Error
	return row.UserID, row.PathID, err
}

func (r *adminPathRepo) PluckPathCourseIDs(ctx context.Context, pathID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Table("path_courses").Where("path_id = ?", pathID).Pluck("course_id", &ids).Error
	return ids, err
}

func (r *adminPathRepo) DeleteUserPathCascade(ctx context.Context, rowID, userID, pathID int, courseIDs []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(courseIDs) > 0 {
			if err := tx.Exec("DELETE FROM user_like_courses WHERE user_id = ? AND course_id IN ?", userID, courseIDs).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_course WHERE user_id = ? AND course_id IN ?", userID, courseIDs).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_course_chapter WHERE user_id = ? AND course_id IN ?", userID, courseIDs).Error; err != nil {
				return err
			}
			now := time.Now()
			if err := tx.Table("course_comments").
				Where("user_id = ? AND course_id IN ? AND deleted_at IS NULL", userID, courseIDs).
				Update("deleted_at", now).Error; err != nil {
				return err
			}
			if err := tx.Table("course_chapter_comments").
				Where("user_id = ? AND course_id IN ? AND deleted_at IS NULL", userID, courseIDs).
				Update("deleted_at", now).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec("DELETE FROM user_path_apply WHERE user_id = ? AND path_id = ?", userID, pathID).Error; err != nil {
			return err
		}
		return tx.Exec("DELETE FROM user_path WHERE id = ?", rowID).Error
	})
}

func (r *adminPathRepo) ListPathApplies(ctx context.Context, f PathApplyListFilter) ([]PathApplyRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_path_apply")
	if f.PathID != nil {
		query = query.Where("path_id = ?", *f.PathID)
	}
	if f.UserID != nil {
		query = query.Where("user_id = ?", *f.UserID)
	}
	if f.PathTitle != "" {
		query = query.Where("path_id IN (SELECT id FROM paths WHERE title LIKE ? AND deleted_at IS NULL)", "%"+f.PathTitle+"%")
	}
	if f.NickName != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
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
	var rows []PathApplyRow
	err := query.Select("id, user_id, path_id, reason, status, created_at, updated_at").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminPathRepo) ListPathThumbs(ctx context.Context, ids []int) ([]PathThumbBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []PathThumbBrief
	err := r.db.WithContext(ctx).Table("paths").
		Select("id, title, thumb").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPathRepo) FindPathApply(ctx context.Context, id int) (*PathApplyRow, error) {
	var row PathApplyRow
	err := r.db.WithContext(ctx).Table("user_path_apply").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPathRepo) PathTitle(ctx context.Context, id int) (string, error) {
	var title string
	err := r.db.WithContext(ctx).Table("paths").Select("title").Where("id = ? AND deleted_at IS NULL", id).Scan(&title).Error
	return title, err
}

func (r *adminPathRepo) UpdatePathApplyStatus(ctx context.Context, id, status int) error {
	return r.db.WithContext(ctx).Table("user_path_apply").Where("id = ?", id).Update("status", status).Error
}

func (r *adminPathRepo) CountPathApply(ctx context.Context, id int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_path_apply").Where("id = ?", id).Count(&n).Error
	return n, err
}

func (r *adminPathRepo) DeletePathApply(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM user_path_apply WHERE id = ?", id).Error
}

func (r *adminPathRepo) SendUserMessage(ctx context.Context, userID int, title, content, relatedType string, relatedID int) error {
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
