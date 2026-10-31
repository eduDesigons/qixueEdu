package model

import "time"

type Cert struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	Name          string     `gorm:"column:name"`
	TemplateImage string     `gorm:"column:template_image"`
	PreviewImage  *string    `gorm:"column:preview_image"`
	Params        *string    `gorm:"column:params"`
	IsShow        uint8      `gorm:"column:is_show"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

func (Cert) TableName() string { return "certs" }

type UserCert struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	CertID    int64      `gorm:"column:cert_id"`
	UserID    int64      `gorm:"column:user_id"`
	CertNo    string     `gorm:"column:cert_no"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	URL       *string    `gorm:"column:url"`
	PathID    *int       `gorm:"column:path_id"`
	CourseID  uint       `gorm:"column:course_id"`
	Title     string     `gorm:"column:title"`
	Image     string     `gorm:"column:image"`
	IsValid   uint8      `gorm:"column:is_valid"`
	IssuedAt  *time.Time `gorm:"column:issued_at"`
}

func (UserCert) TableName() string { return "user_cert" }
