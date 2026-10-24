package adminpaper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/repo"
)

type Service struct {
	papers repo.AdminPaperRepository
	db     *gorm.DB
	rdb    *redis.Client
}

func New(papers repo.AdminPaperRepository, db *gorm.DB, rdb *redis.Client) *Service {
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

func (s *Service) GetPaperDetail(c *gin.Context) {
	id := bind.Query(c, "id")
	var paper paperRow
	err := s.db.WithContext(c.Request.Context()).
		Select("id, course_id, chapter_id, score, pass_score").
		Where("id = ?", id).Take(&paper).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "试卷不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	var questions []questionRow
	if response.WriteErr(c, s.db.WithContext(c.Request.Context()).
		Where("paper_id = ?", id).Order("sort ASC, id ASC").Find(&questions).Error) {
		return
	}
	if questions == nil {
		questions = []questionRow{}
	}
	paper.Questions = questions
	response.OK(c, "", paper)
}

func (s *Service) PaperSave(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	questions, _ := raw["questions"].([]any)
	if len(questions) == 0 {
		response.Fail(c, 1, "请添加问题")
		return
	}
	isSkipMark := 1
	qInputs := make([]repo.PaperQuestionInput, 0, len(questions))
	for _, item := range questions {
		q, _ := item.(map[string]any)
		typ := asString(q["type"])
		if typ == "" {
			response.Fail(c, 1, "请添加问题")
			return
		}
		if typ != "single" && typ != "multiple" && typ != "completion" && typ != "text" {
			isSkipMark = 0
		}
		answer, _ := json.Marshal(q["answer"])
		options, _ := json.Marshal(q["options"])
		if string(options) == "null" {
			options = []byte("null")
		}
		qInputs = append(qInputs, repo.PaperQuestionInput{
			Type: typ, Title: asString(q["title"]), Answer: answer,
			Score: asInt(q["score"]), Options: options,
		})
	}
	courseID := asInt(raw["course_id"])
	chapterID := asInt(raw["chapter_id"])
	isExam := chapterID == 0
	ctx := c.Request.Context()
	courseCount, err := s.papers.CountCourse(ctx, courseID)
	if response.WriteErr(c, err) {
		return
	}
	if courseCount == 0 {
		response.Fail(c, 1, "课程不存在")
		return
	}
	id := asInt(raw["id"])
	if !isExam {
		chapter, err := s.papers.FindChapterRef(ctx, chapterID)
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, 1, "课节不存在")
			return
		}
		if response.WriteErr(c, err) {
			return
		}
		if chapter.CourseID != courseID {
			response.Fail(c, 1, "课节与课程不匹配")
			return
		}
		if id > 0 && chapter.PaperID != nil && *chapter.PaperID != id {
			response.Fail(c, 1, "课节已经设置过测评了，不能变更")
			return
		}
	}
	err = s.papers.SavePaper(ctx, repo.PaperSaveInput{
		ID: id, CourseID: courseID, ChapterID: chapterID, IsExam: isExam,
		Score: asInt(raw["score"]), PassScore: asInt(raw["pass_score"]),
		IsSkipMark: isSkipMark, CreatedID: middleware.UID(c), Questions: qInputs,
	})
	if errors.Is(err, repo.ErrPaperCourseMismatch) {
		response.Fail(c, 1, err.Error())
		return
	}
	if errors.Is(err, repo.ErrPaperNotFound) {
		response.Fail(c, 1, err.Error())
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if asInt(raw["id"]) > 0 {
		response.OK(c, "更新成功", []any{})
		return
	}
	response.OK(c, "成功", []any{})
}

func (s *Service) PaperDelete(c *gin.Context) {
	id := bind.Query(c, "id")
	ctx := c.Request.Context()
	var paper paperRow
	err := s.db.WithContext(ctx).Where("id = ?", id).Take(&paper).Error
	if err == gorm.ErrRecordNotFound {
		response.Fail(c, 1, "试卷不存在")
		return
	}
	if response.WriteErr(c, err) {
		return
	}
	if paper.CourseID == nil || *paper.CourseID <= 0 {
		response.Fail(c, 1, "课程不存在")
		return
	}
	courseID := *paper.CourseID
	courseExists, err := s.papers.CountCourse(ctx, courseID)
	if response.WriteErr(c, err) {
		return
	}
	if courseExists == 0 {
		response.Fail(c, 1, "课程不存在")
		return
	}
	users, err := s.papers.CountCourseUsers(ctx, courseID)
	if response.WriteErr(c, err) {
		return
	}
	if users > 0 {
		response.Fail(c, 1, "已经有学员加入的可能，不能变更试卷")
		return
	}
	if response.WriteErr(c, s.papers.DeletePaperCascade(ctx, id, courseID, paper.ChapterID)) {
		return
	}
	response.OK(c, "成功", []any{})
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}

func asInt(v any) int { return bind.AsInt(v) }
