package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// AdminCertRepository covers admin certificate Table-based access.
type AdminCertRepository interface {
	ListOptions(ctx context.Context) ([]AdminCertOptionRow, error)
	CountPathsByCert(ctx context.Context, certID int) (int64, error)
	CountCoursesByCert(ctx context.Context, certID int) (int64, error)
	CountUserCertsByCert(ctx context.Context, certID int) (int64, error)
	ListIssued(ctx context.Context, filter AdminCertIssuedFilter, offset, limit int) ([]AdminCertIssuedRow, int64, error)
	ListUsersBrief(ctx context.Context, ids []int) ([]AdminCertUserBrief, error)
	ListTitles(ctx context.Context, table string, ids []int) ([]AdminCertTitleBrief, error)
	ListCertNames(ctx context.Context, ids []int) ([]AdminCertNameBrief, error)
}

type AdminCertOptionRow struct {
	ID            int     `gorm:"column:id" json:"id"`
	Name          string  `gorm:"column:name" json:"name"`
	TemplateImage string  `gorm:"column:template_image" json:"template_image"`
	PreviewImage  *string `gorm:"column:preview_image" json:"preview_image"`
	IsShow        int     `gorm:"column:is_show" json:"is_show"`
}

type AdminCertIssuedFilter struct {
	CourseID int
	PathID   int
	Name     string
	CityID   int
	Limited  bool
}

type AdminCertIssuedRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CertID    int        `gorm:"column:cert_id"`
	PathID    *int       `gorm:"column:path_id"`
	CourseID  int        `gorm:"column:course_id"`
	CertNo    string     `gorm:"column:cert_no"`
	Title     string     `gorm:"column:title"`
	Image     string     `gorm:"column:image"`
	IsValid   int        `gorm:"column:is_valid"`
	IssuedAt  *time.Time `gorm:"column:issued_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
}

type AdminCertUserBrief struct {
	ID       int     `gorm:"column:id" json:"id"`
	NickName *string `gorm:"column:nick_name" json:"nick_name"`
	Avatar   string  `gorm:"column:avatar" json:"avatar"`
}

type AdminCertTitleBrief struct {
	ID    int    `gorm:"column:id" json:"id"`
	Title string `gorm:"column:title" json:"title"`
}

type AdminCertNameBrief struct {
	ID   int    `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type adminCertRepo struct{ db *gorm.DB }

func NewAdminCert(db *gorm.DB) AdminCertRepository {
	return &adminCertRepo{db: db}
}

func (r *adminCertRepo) ListOptions(ctx context.Context) ([]AdminCertOptionRow, error) {
	var rows []AdminCertOptionRow
	err := r.db.WithContext(ctx).Table("certs").
		Select("id, name, template_image, preview_image, is_show").
		Order("id DESC").Find(&rows).Error
	return rows, err
}

func (r *adminCertRepo) CountPathsByCert(ctx context.Context, certID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("paths").
		Where("cert_id = ? AND deleted_at IS NULL", certID).Count(&n).Error
	return n, err
}

func (r *adminCertRepo) CountCoursesByCert(ctx context.Context, certID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("courses").
		Where("cert_id = ? AND deleted_at IS NULL", certID).Count(&n).Error
	return n, err
}

func (r *adminCertRepo) CountUserCertsByCert(ctx context.Context, certID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_cert").Where("cert_id = ?", certID).Count(&n).Error
	return n, err
}

func (r *adminCertRepo) ListIssued(ctx context.Context, filter AdminCertIssuedFilter, offset, limit int) ([]AdminCertIssuedRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_cert").
		Select("id, user_id, cert_id, path_id, course_id, cert_no, title, image, is_valid, issued_at, created_at")
	if filter.CourseID > 0 {
		query = query.Where("course_id = ? AND path_id = 0", filter.CourseID)
	}
	if filter.PathID > 0 {
		query = query.Where("path_id = ? AND course_id = 0", filter.PathID)
	}
	if filter.Name != "" {
		like := "%" + filter.Name + "%"
		query = query.Where(`(
			title LIKE ? OR cert_no LIKE ?
			OR user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)
			OR path_id IN (SELECT id FROM paths WHERE title LIKE ? AND deleted_at IS NULL)
			OR course_id IN (SELECT id FROM courses WHERE title LIKE ? AND deleted_at IS NULL)
		)`, like, like, like, like, like)
	}
	if filter.Limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", filter.CityID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []AdminCertIssuedRow
	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminCertRepo) ListUsersBrief(ctx context.Context, ids []int) ([]AdminCertUserBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminCertUserBrief
	err := r.db.WithContext(ctx).Table("users").Select("id, nick_name, avatar").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminCertRepo) ListTitles(ctx context.Context, table string, ids []int) ([]AdminCertTitleBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminCertTitleBrief
	err := r.db.WithContext(ctx).Table(table).Select("id, title").
		Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminCertRepo) ListCertNames(ctx context.Context, ids []int) ([]AdminCertNameBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []AdminCertNameBrief
	err := r.db.WithContext(ctx).Table("certs").Select("id, name").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}
