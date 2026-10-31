package adminpaper

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	adminrepo "qixue-api/internal/admin/repo"
)

// BizError is a business failure mapped to HTTP Fail(code=1) by the handler.
type BizError struct{ Message string }

func (e *BizError) Error() string { return e.Message }

func biz(msg string) error { return &BizError{Message: msg} }

type Service struct {
	papers adminrepo.AdminPaperRepository
	db     *gorm.DB
	rdb    *redis.Client
}

func New(papers adminrepo.AdminPaperRepository, db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{papers: papers, db: db, rdb: rdb}
}

type paperRow struct {
	ID         int            `gorm:"column:id;primaryKey" json:"id"`
	CourseID   *int           `gorm:"column:course_id" json:"course_id"`
	ChapterID  *int           `gorm:"column:chapter_id" json:"chapter_id"`
	Score      int            `gorm:"column:score" json:"score"`
	PassScore  int            `gorm:"column:pass_score" json:"pass_score"`
	CreatedID  *int           `gorm:"column:created_id" json:"created_id"`
	IsSkipMark int            `gorm:"column:is_skip_mark" json:"is_skip_mark"`
	Type       string         `gorm:"column:type" json:"type"`
	CreatedAt  *time.Time     `gorm:"column:created_at" json:"-"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at" json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
	Questions  []questionRow  `gorm:"-" json:"questions,omitempty"`
}

func (paperRow) TableName() string { return "papers" }

type questionRow struct {
	ID        int64           `gorm:"column:id;primaryKey" json:"id"`
	PaperID   int             `gorm:"column:paper_id" json:"paper_id"`
	Sort      int             `gorm:"column:sort" json:"sort"`
	Type      string          `gorm:"column:type" json:"type"`
	Title     string          `gorm:"column:title" json:"title"`
	Answer    json.RawMessage `gorm:"column:answer" json:"answer"`
	Score     int             `gorm:"column:score" json:"score"`
	Options   json.RawMessage `gorm:"column:options" json:"options"`
	CreatedAt *time.Time      `gorm:"column:created_at" json:"-"`
	UpdatedAt *time.Time      `gorm:"column:updated_at" json:"-"`
}

func (questionRow) TableName() string { return "paper_questions" }

func (s *Service) GetPaperDetail(ctx context.Context, id int) (paperRow, error) {
	var paper paperRow
	err := s.db.WithContext(ctx).
		Select("id, course_id, chapter_id, score, pass_score").
		Where("id = ?", id).Take(&paper).Error
	if err == gorm.ErrRecordNotFound {
		return paperRow{}, biz("试卷不存在")
	}
	if err != nil {
		return paperRow{}, err
	}
	var questions []questionRow
	if err := s.db.WithContext(ctx).
		Where("paper_id = ?", id).Order("sort ASC, id ASC").Find(&questions).Error; err != nil {
		return paperRow{}, err
	}
	if questions == nil {
		questions = []questionRow{}
	}
	paper.Questions = questions
	return paper, nil
}

func (s *Service) PaperSave(ctx context.Context, adminID int, raw map[string]any) (string, error) {
	if raw == nil {
		return "", biz("参数错误")
	}
	questions, _ := raw["questions"].([]any)
	if len(questions) == 0 {
		return "", biz("请添加问题")
	}
	isSkipMark := 1
	qInputs := make([]adminrepo.PaperQuestionInput, 0, len(questions))
	for _, item := range questions {
		q, _ := item.(map[string]any)
		typ := asString(q["type"])
		if typ == "" {
			return "", biz("请添加问题")
		}
		if typ != "single" && typ != "multiple" && typ != "completion" && typ != "text" {
			isSkipMark = 0
		}
		answer, _ := json.Marshal(q["answer"])
		options, _ := json.Marshal(q["options"])
		if string(options) == "null" {
			options = []byte("null")
		}
		qInputs = append(qInputs, adminrepo.PaperQuestionInput{
			Type: typ, Title: asString(q["title"]), Answer: answer,
			Score: asInt(q["score"]), Options: options,
		})
	}
	courseID := asInt(raw["course_id"])
	chapterID := asInt(raw["chapter_id"])
	isExam := chapterID == 0
	courseCount, err := s.papers.CountCourse(ctx, courseID)
	if err != nil {
		return "", err
	}
	if courseCount == 0 {
		return "", biz("课程不存在")
	}
	id := asInt(raw["id"])
	if !isExam {
		chapter, err := s.papers.FindChapterRef(ctx, chapterID)
		if err == gorm.ErrRecordNotFound {
			return "", biz("课节不存在")
		}
		if err != nil {
			return "", err
		}
		if chapter.CourseID != courseID {
			return "", biz("课节与课程不匹配")
		}
		if id > 0 && chapter.PaperID != nil && *chapter.PaperID != id {
			return "", biz("课节已经设置过测评了，不能变更")
		}
	}
	err = s.papers.SavePaper(ctx, adminrepo.PaperSaveInput{
		ID: id, CourseID: courseID, ChapterID: chapterID, IsExam: isExam,
		Score: asInt(raw["score"]), PassScore: asInt(raw["pass_score"]),
		IsSkipMark: isSkipMark, CreatedID: adminID, Questions: qInputs,
	})
	if errors.Is(err, adminrepo.ErrPaperCourseMismatch) {
		return "", biz(err.Error())
	}
	if errors.Is(err, adminrepo.ErrPaperNotFound) {
		return "", biz(err.Error())
	}
	if err != nil {
		return "", err
	}
	if asInt(raw["id"]) > 0 {
		return "更新成功", nil
	}
	return "成功", nil
}

func (s *Service) PaperDelete(ctx context.Context, id int) error {
	var paper paperRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&paper).Error
	if err == gorm.ErrRecordNotFound {
		return biz("试卷不存在")
	}
	if err != nil {
		return err
	}
	if paper.CourseID == nil || *paper.CourseID <= 0 {
		return biz("课程不存在")
	}
	courseID := *paper.CourseID
	courseExists, err := s.papers.CountCourse(ctx, courseID)
	if err != nil {
		return err
	}
	if courseExists == 0 {
		return biz("课程不存在")
	}
	users, err := s.papers.CountCourseUsers(ctx, courseID)
	if err != nil {
		return err
	}
	if users > 0 {
		return biz("已经有学员加入的可能，不能变更试卷")
	}
	return s.papers.DeletePaperCascade(ctx, id, courseID, paper.ChapterID)
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}

func asInt(v any) int { return bind.AsInt(v) }
