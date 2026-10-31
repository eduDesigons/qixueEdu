package model

import "time"

type Course struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	Title         string     `gorm:"column:title"`
	Thumb         string     `gorm:"column:thumb"`
	Desc          string     `gorm:"column:desc"`
	Type          int        `gorm:"column:type"`
	Content       *string    `gorm:"column:content"`
	IsShow        int        `gorm:"column:is_show"`
	CategoryID    int        `gorm:"column:category_id"`
	IsRec         int        `gorm:"column:is_rec"`
	MinNeedCredit int        `gorm:"column:min_need_credit"`
	Credit        int        `gorm:"column:credit"`
	PaperID       *int       `gorm:"column:paper_id"`
	CertID        *uint      `gorm:"column:cert_id"`
	CreatedID     int        `gorm:"column:created_id"`
	ViewCount     int        `gorm:"column:view_count"`
	ChapterCount  *int       `gorm:"column:chapter_count"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (Course) TableName() string { return "courses" }

type CourseChapter struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	CourseID  int        `gorm:"column:course_id"`
	Title     string     `gorm:"column:title"`
	MediaID   *int       `gorm:"column:media_id"`
	MediaURL  *string    `gorm:"column:media_url"`
	Type      int        `gorm:"column:type"`
	Content   *string    `gorm:"column:content"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	CreatedID int        `gorm:"column:created_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	PaperID   *int       `gorm:"column:paper_id"`
	Sort      int        `gorm:"column:sort"`
	ViewCount int        `gorm:"column:view_count"`
}

func (CourseChapter) TableName() string { return "course_chapters" }

type CourseCategory struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	Name          string     `gorm:"column:name"`
	Icon          *string    `gorm:"column:icon"`
	ParentID      int        `gorm:"column:parent_id"`
	IsShow        int        `gorm:"column:is_show"`
	Sort          int        `gorm:"column:sort"`
	MinNeedCredit int        `gorm:"column:min_need_credit"`
	BaseCredit    int        `gorm:"column:base_credit"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (CourseCategory) TableName() string { return "course_categories" }

type CourseAttach struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	CourseID      int        `gorm:"column:course_id"`
	Name          string     `gorm:"column:name"`
	Path          string     `gorm:"column:path"`
	Size          int        `gorm:"column:size"`
	DownloadTimes int        `gorm:"column:download_times"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (CourseAttach) TableName() string { return "course_attach" }
