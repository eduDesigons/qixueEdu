package model

import "time"

type UserCourse struct {
	ID            int        `gorm:"column:id;primaryKey"`
	UserID        *int       `gorm:"column:user_id"`
	CourseID      *int       `gorm:"column:course_id"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	HasEndChapter int        `gorm:"column:has_end_chapter"`
	IsEnd         int        `gorm:"column:is_end"`
}

func (UserCourse) TableName() string { return "user_course" }

type UserCourseChapter struct {
	ID            int        `gorm:"column:id;primaryKey"`
	UserID        int        `gorm:"column:user_id"`
	ChapterID     int        `gorm:"column:chapter_id"`
	CourseID      *int       `gorm:"column:course_id"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	HasEndChapter *int       `gorm:"column:has_end_chapter"`
	IsEnd         int        `gorm:"column:is_end"`
	UnLock        int        `gorm:"column:un_lock"`
	LearnTime     int        `gorm:"column:learn_time"`
	ExpireTime    *time.Time `gorm:"column:expire_time"`
}

func (UserCourseChapter) TableName() string { return "user_course_chapter" }

type UserCourseApply struct {
	ID        int        `gorm:"column:id;primaryKey"`
	UserID    *int       `gorm:"column:user_id"`
	CourseID  *int       `gorm:"column:course_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	Reason    *string    `gorm:"column:reason"`
	Status    int        `gorm:"column:status"`
}

func (UserCourseApply) TableName() string { return "user_course_apply" }

type UserPath struct {
	ID           uint       `gorm:"column:id;primaryKey"`
	UserID       int        `gorm:"column:user_id"`
	PathID       int        `gorm:"column:path_id"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	HasEndCourse int        `gorm:"column:has_end_course"`
	IsEnd        int        `gorm:"column:is_end"`
}

func (UserPath) TableName() string { return "user_path" }

type UserPathApply struct {
	ID        int        `gorm:"column:id;primaryKey"`
	UserID    *int       `gorm:"column:user_id"`
	PathID    *int       `gorm:"column:path_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	Reason    *string    `gorm:"column:reason"`
	Status    int        `gorm:"column:status"`
}

func (UserPathApply) TableName() string { return "user_path_apply" }
