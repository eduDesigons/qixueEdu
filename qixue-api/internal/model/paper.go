package model

import "time"

type Paper struct {
	ID         uint       `gorm:"column:id;primaryKey"`
	CourseID   *int       `gorm:"column:course_id"`
	ChapterID  *int       `gorm:"column:chapter_id"`
	Score      int        `gorm:"column:score"`
	PassScore  int        `gorm:"column:pass_score"`
	CreatedID  *int       `gorm:"column:created_id"`
	CreatedAt  *time.Time `gorm:"column:created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
	IsSkipMark int        `gorm:"column:is_skip_mark"`
	Type       string     `gorm:"column:type"`
}

func (Paper) TableName() string { return "papers" }

type PaperQuestion struct {
	ID        uint       `gorm:"column:id;primaryKey"`
	PaperID   *int       `gorm:"column:paper_id"`
	Sort      int        `gorm:"column:sort"`
	Type      string     `gorm:"column:type"`
	Title     *string    `gorm:"column:title"`
	Answer    string     `gorm:"column:answer"`
	Score     int        `gorm:"column:score"`
	Options   *string    `gorm:"column:options"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (PaperQuestion) TableName() string { return "paper_questions" }

type UserPaper struct {
	ID            uint       `gorm:"column:id;primaryKey"`
	UserID        int        `gorm:"column:user_id"`
	PaperID       int        `gorm:"column:paper_id"`
	Score         int        `gorm:"column:score"`
	SubmitCount   int        `gorm:"column:submit_count"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
	IsPass        int        `gorm:"column:is_pass"`
	IsPeopleCheck int        `gorm:"column:is_people_check"`
}

func (UserPaper) TableName() string { return "user_paper" }

type UserPaperQuestionLog struct {
	ID              uint       `gorm:"column:id;primaryKey"`
	UserID          int64      `gorm:"column:user_id"`
	UserPaperID     int        `gorm:"column:user_paper_id"`
	PaperID         int64      `gorm:"column:paper_id"`
	QuestionID      int64      `gorm:"column:question_id"`
	AnswerContent   string     `gorm:"column:answer_content"`
	IsCorrect       int        `gorm:"column:is_correct"`
	Score           int        `gorm:"column:score"`
	CreatedAt       *time.Time `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
	NeedPeopleCheck int        `gorm:"column:need_people_check"`
	GradeRemark     string     `gorm:"column:grade_remark"`
}

func (UserPaperQuestionLog) TableName() string { return "user_paper_question_log" }
