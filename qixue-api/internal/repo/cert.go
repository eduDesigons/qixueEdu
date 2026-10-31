package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CertRepository covers certificate issue / render data access.
type CertRepository interface {
	GetCourseBrief(ctx context.Context, courseID int) (*CertSubjectBrief, error)
	GetPathBrief(ctx context.Context, pathID int) (*CertSubjectBrief, error)
	GetCourseCertID(ctx context.Context, courseID int) (*int, error)
	GetPathCertID(ctx context.Context, pathID int) (*int, error)
	GetUserCourseEnded(ctx context.Context, userID, courseID int) (int, error)
	GetUserPathEnded(ctx context.Context, userID, pathID int) (int, error)
	UserExists(ctx context.Context, userID int) (bool, error)
	LoadTemplateMeta(ctx context.Context, certID int) (*CertTemplateMeta, error)
	LoadTemplateRender(ctx context.Context, certID int) (*CertTemplateRender, error)
	FindUserCert(ctx context.Context, userID, pathID, courseID int) (*UserCert, error)
	GetUserCertByID(ctx context.Context, id int) (*UserCert, error)
	UpdateUserCert(ctx context.Context, id int, payload map[string]any) error
	CreateUserCert(ctx context.Context, row *UserCert) error
	UpdateUserCertImage(ctx context.Context, id int, image string, updatedAt time.Time) error
	ListEndedCourseUserIDs(ctx context.Context, courseID, limit int) ([]int, error)
	ListEndedPathUserIDs(ctx context.Context, pathID, limit int) ([]int, error)
	CreateCertNotice(ctx context.Context, userID, certID int, title, content string, now time.Time) error
	ListPathIDsByCourse(ctx context.Context, courseID int) ([]int, error)
	GetUserPathProgress(ctx context.Context, userID, pathID int) (*UserPathProgress, error)
	ListActivePathCourseIDs(ctx context.Context, pathID int) ([]int, error)
	CountEndedUserCourses(ctx context.Context, userID int, courseIDs []int) (int64, error)
	UpdateUserPathProgress(ctx context.Context, id, hasEndCourse, isEnd int) error
	GetUserDisplay(ctx context.Context, userID int) (nickName, account string, err error)
	GetCourseTitle(ctx context.Context, courseID int) (string, error)
	GetPathTitle(ctx context.Context, pathID int) (string, error)
}

type UserCert struct {
	ID        int        `gorm:"column:id;primaryKey" json:"id"`
	UserID    int        `gorm:"column:user_id" json:"user_id"`
	CertID    int        `gorm:"column:cert_id" json:"cert_id"`
	PathID    int        `gorm:"column:path_id" json:"path_id"`
	CourseID  int        `gorm:"column:course_id" json:"course_id"`
	CertNo    string     `gorm:"column:cert_no" json:"cert_no"`
	Title     string     `gorm:"column:title" json:"title"`
	Image     string     `gorm:"column:image" json:"image"`
	IsValid   int        `gorm:"column:is_valid" json:"is_valid"`
	IssuedAt  *time.Time `gorm:"column:issued_at" json:"issued_at"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (UserCert) TableName() string { return "user_cert" }

type CertSubjectBrief struct {
	ID     int    `gorm:"column:id"`
	Title  string `gorm:"column:title"`
	CertID *int   `gorm:"column:cert_id"`
}

type CertTemplateMeta struct {
	ID     int    `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	IsShow int    `gorm:"column:is_show"`
}

type CertTemplateRender struct {
	ID            int             `gorm:"column:id"`
	TemplateImage string          `gorm:"column:template_image"`
	Params        json.RawMessage `gorm:"column:params"`
}

type UserPathProgress struct {
	ID           int `gorm:"column:id"`
	HasEndCourse int `gorm:"column:has_end_course"`
	IsEnd        int `gorm:"column:is_end"`
}

type certRepo struct{ db *gorm.DB }

func NewCert(db *gorm.DB) CertRepository {
	return &certRepo{db: db}
}

func (r *certRepo) GetCourseBrief(ctx context.Context, courseID int) (*CertSubjectBrief, error) {
	var row CertSubjectBrief
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, cert_id").
		Where("id = ? AND deleted_at IS NULL", courseID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *certRepo) GetPathBrief(ctx context.Context, pathID int) (*CertSubjectBrief, error) {
	var row CertSubjectBrief
	err := r.db.WithContext(ctx).Table("paths").
		Select("id, title, cert_id").
		Where("id = ? AND deleted_at IS NULL", pathID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *certRepo) GetCourseCertID(ctx context.Context, courseID int) (*int, error) {
	var certID *int
	err := r.db.WithContext(ctx).Table("courses").Select("cert_id").
		Where("id = ? AND deleted_at IS NULL", courseID).Scan(&certID).Error
	return certID, err
}

func (r *certRepo) GetPathCertID(ctx context.Context, pathID int) (*int, error) {
	var certID *int
	err := r.db.WithContext(ctx).Table("paths").Select("cert_id").
		Where("id = ? AND deleted_at IS NULL", pathID).Scan(&certID).Error
	return certID, err
}

func (r *certRepo) GetUserCourseEnded(ctx context.Context, userID, courseID int) (int, error) {
	var ended int
	err := r.db.WithContext(ctx).Table("user_course").
		Select("is_end").Where("user_id = ? AND course_id = ?", userID, courseID).
		Scan(&ended).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return ended, nil
}

func (r *certRepo) GetUserPathEnded(ctx context.Context, userID, pathID int) (int, error) {
	var ended int
	err := r.db.WithContext(ctx).Table("user_path").
		Select("is_end").Where("user_id = ? AND path_id = ?", userID, pathID).
		Scan(&ended).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return ended, nil
}

func (r *certRepo) UserExists(ctx context.Context, userID int) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", userID).Count(&n).Error
	return n > 0, err
}

func (r *certRepo) LoadTemplateMeta(ctx context.Context, certID int) (*CertTemplateMeta, error) {
	var tpl CertTemplateMeta
	err := r.db.WithContext(ctx).Table("certs").
		Select("id, name, is_show").Where("id = ?", certID).Take(&tpl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (r *certRepo) LoadTemplateRender(ctx context.Context, certID int) (*CertTemplateRender, error) {
	var tpl CertTemplateRender
	err := r.db.WithContext(ctx).Table("certs").
		Select("id, template_image, params").
		Where("id = ?", certID).Take(&tpl).Error
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (r *certRepo) FindUserCert(ctx context.Context, userID, pathID, courseID int) (*UserCert, error) {
	var row UserCert
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if courseID > 0 {
		q = q.Where("course_id = ? AND path_id = 0", courseID)
	} else {
		q = q.Where("path_id = ? AND course_id = 0", pathID)
	}
	err := q.Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *certRepo) GetUserCertByID(ctx context.Context, id int) (*UserCert, error) {
	var row UserCert
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *certRepo) UpdateUserCert(ctx context.Context, id int, payload map[string]any) error {
	return r.db.WithContext(ctx).Table("user_cert").Where("id = ?", id).Updates(payload).Error
}

func (r *certRepo) CreateUserCert(ctx context.Context, row *UserCert) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *certRepo) UpdateUserCertImage(ctx context.Context, id int, image string, updatedAt time.Time) error {
	return r.db.WithContext(ctx).Table("user_cert").Where("id = ?", id).
		Updates(map[string]any{"image": image, "updated_at": updatedAt}).Error
}

func (r *certRepo) ListEndedCourseUserIDs(ctx context.Context, courseID, limit int) ([]int, error) {
	var userIDs []int
	err := r.db.WithContext(ctx).Table("user_course").
		Where("course_id = ? AND is_end = 1", courseID).
		Order("id ASC").Limit(limit).Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *certRepo) ListEndedPathUserIDs(ctx context.Context, pathID, limit int) ([]int, error) {
	var userIDs []int
	err := r.db.WithContext(ctx).Table("user_path").
		Where("path_id = ? AND is_end = 1", pathID).
		Order("id ASC").Limit(limit).Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *certRepo) CreateCertNotice(ctx context.Context, userID, certID int, title, content string, now time.Time) error {
	return r.db.WithContext(ctx).Table("user_messages").Create(map[string]any{
		"user_id": userID, "admin_id": 0, "title": title, "content": content,
		"type": "notice", "related_type": "user_cert", "related_id": certID,
		"is_read": 0, "created_at": now, "updated_at": now,
	}).Error
}

func (r *certRepo) ListPathIDsByCourse(ctx context.Context, courseID int) ([]int, error) {
	var pathIDs []int
	err := r.db.WithContext(ctx).Table("path_courses").
		Where("course_id = ?", courseID).Pluck("path_id", &pathIDs).Error
	return pathIDs, err
}

func (r *certRepo) GetUserPathProgress(ctx context.Context, userID, pathID int) (*UserPathProgress, error) {
	var row UserPathProgress
	err := r.db.WithContext(ctx).Table("user_path").
		Where("user_id = ? AND path_id = ?", userID, pathID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *certRepo) ListActivePathCourseIDs(ctx context.Context, pathID int) ([]int, error) {
	var courseIDs []int
	err := r.db.WithContext(ctx).Table("path_courses").
		Where("path_id = ? AND course_id IN (SELECT id FROM courses WHERE deleted_at IS NULL)", pathID).
		Pluck("course_id", &courseIDs).Error
	return courseIDs, err
}

func (r *certRepo) CountEndedUserCourses(ctx context.Context, userID int, courseIDs []int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id IN ? AND is_end = 1", userID, courseIDs).
		Count(&n).Error
	return n, err
}

func (r *certRepo) UpdateUserPathProgress(ctx context.Context, id, hasEndCourse, isEnd int) error {
	return r.db.WithContext(ctx).Table("user_path").Where("id = ?", id).
		UpdateColumns(map[string]any{"has_end_course": hasEndCourse, "is_end": isEnd}).Error
}

func (r *certRepo) GetUserDisplay(ctx context.Context, userID int) (string, string, error) {
	var row struct {
		NickName *string `gorm:"column:nick_name"`
		Account  *string `gorm:"column:account"`
	}
	err := r.db.WithContext(ctx).Table("users").
		Select("nick_name, account").Where("id = ?", userID).Take(&row).Error
	if err != nil {
		return "", "", err
	}
	nick, account := "", ""
	if row.NickName != nil {
		nick = *row.NickName
	}
	if row.Account != nil {
		account = *row.Account
	}
	return nick, account, nil
}

func (r *certRepo) GetCourseTitle(ctx context.Context, courseID int) (string, error) {
	var title string
	err := r.db.WithContext(ctx).Table("courses").Select("title").
		Where("id = ? AND deleted_at IS NULL", courseID).Scan(&title).Error
	return title, err
}

func (r *certRepo) GetPathTitle(ctx context.Context, pathID int) (string, error) {
	var title string
	err := r.db.WithContext(ctx).Table("paths").Select("title").
		Where("id = ? AND deleted_at IS NULL", pathID).Scan(&title).Error
	return title, err
}
