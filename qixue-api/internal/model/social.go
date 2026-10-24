package model

import "time"

type CourseComment struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	PID       int        `gorm:"column:pid"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	Content   string     `gorm:"column:content"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (CourseComment) TableName() string { return "course_comments" }

type CourseChapterComment struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	PID       int        `gorm:"column:pid"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	ChapterID int        `gorm:"column:chapter_id"`
	Content   string     `gorm:"column:content"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (CourseChapterComment) TableName() string { return "course_chapter_comments" }

type UserLikeCourse struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (UserLikeCourse) TableName() string { return "user_like_courses" }

type UserLikeChapter struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	ChapterID int        `gorm:"column:chapter_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (UserLikeChapter) TableName() string { return "user_like_chapters" }

type UserMessage struct {
	ID          uint       `gorm:"column:id;primaryKey"`
	UserID      uint       `gorm:"column:user_id"`
	AdminID     uint       `gorm:"column:admin_id"`
	Title       string     `gorm:"column:title"`
	Content     string     `gorm:"column:content"`
	Type        string     `gorm:"column:type"`
	RelatedType *string    `gorm:"column:related_type"`
	RelatedID   *uint      `gorm:"column:related_id"`
	IsRead      uint8      `gorm:"column:is_read"`
	ReadAt      *time.Time `gorm:"column:read_at"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (UserMessage) TableName() string { return "user_messages" }

type UserFeedback struct {
	ID           uint       `gorm:"column:id;primaryKey"`
	UserID       uint       `gorm:"column:user_id"`
	Content      string     `gorm:"column:content"`
	Images       *string    `gorm:"column:images"`
	Status       uint8      `gorm:"column:status"`
	Reply        *string    `gorm:"column:reply"`
	ReplyAdminID uint       `gorm:"column:reply_admin_id"`
	RepliedAt    *time.Time `gorm:"column:replied_at"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (UserFeedback) TableName() string { return "user_feedbacks" }

type UserCoinsRecord struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	UserID    int        `gorm:"column:user_id"`
	Field     string     `gorm:"column:field"`
	Sum       int        `gorm:"column:sum"`
	Remark    string     `gorm:"column:remark"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (UserCoinsRecord) TableName() string { return "user_coins_records" }

type UserCreditRecord struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	UserID    int        `gorm:"column:user_id"`
	Field     string     `gorm:"column:field"`
	Sum       int        `gorm:"column:sum"`
	Remark    string     `gorm:"column:remark"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (UserCreditRecord) TableName() string { return "user_credit_records" }
