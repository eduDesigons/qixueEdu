package model

import "time"

type Path struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	Title         string     `gorm:"column:title"`
	Thumb         string     `gorm:"column:thumb"`
	Desc          string     `gorm:"column:desc"`
	Content       *string    `gorm:"column:content"`
	IsOpen        int        `gorm:"column:is_open"`
	IsShow        int        `gorm:"column:is_show"`
	IsRec         int        `gorm:"column:is_rec"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	CreatedID     *int       `gorm:"column:created_id"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
	CategoryID    int        `gorm:"column:category_id"`
	CertID        *uint      `gorm:"column:cert_id"`
	MinNeedCredit int        `gorm:"column:min_need_credit"`
	ViewCount     int        `gorm:"column:view_count"`
}

func (Path) TableName() string { return "paths" }

type PathCategory struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	ParentID  int        `gorm:"column:parent_id"`
	Name      string     `gorm:"column:name"`
	Sort      int        `gorm:"column:sort"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (PathCategory) TableName() string { return "path_categories" }

type PathCourse struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	PathID    int        `gorm:"column:path_id"`
	CourseID  int        `gorm:"column:course_id"`
	Sort      int        `gorm:"column:sort"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (PathCourse) TableName() string { return "path_courses" }
