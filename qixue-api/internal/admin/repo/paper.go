package admin

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrPaperCourseMismatch = errors.New("课程不能更换试卷")
	ErrPaperNotFound       = errors.New("试卷不存在")
)

type UserPaperListFilter struct {
	PaperType     string
	Exam          bool
	CourseID      int
	ChapterID     int
	Title         string
	IsPass        *int
	IsPeopleCheck *int
	NickName      string
	CityID        int
	Limited       bool
	Sort          string
	Order         string
	Offset        int
	Limit         int
}

type UserPaperLogFilter struct {
	PaperID  int
	NickName string
	CityID   int
	Limited  bool
	Offset   int
	Limit    int
}

type AdminUserPaperRow struct {
	ID            int        `gorm:"column:id"`
	UserID        int        `gorm:"column:user_id"`
	PaperID       int        `gorm:"column:paper_id"`
	Score         int        `gorm:"column:score"`
	IsPass        int        `gorm:"column:is_pass"`
	IsPeopleCheck int        `gorm:"column:is_people_check"`
	SubmitCount   int        `gorm:"column:submit_count"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
}

type PaperBrief struct {
	ID        int  `gorm:"column:id"`
	CourseID  *int `gorm:"column:course_id"`
	ChapterID *int `gorm:"column:chapter_id"`
}

type UserNickBrief struct {
	ID       int     `gorm:"column:id"`
	NickName *string `gorm:"column:nick_name"`
	Avatar   string  `gorm:"column:avatar"`
	Credit   int     `gorm:"column:credit"`
}

type CourseThumb struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
	Thumb string `gorm:"column:thumb"`
}

type ChapterTitle struct {
	ID    int    `gorm:"column:id"`
	Title string `gorm:"column:title"`
}

type PaperChapterRef struct {
	ID       int  `gorm:"column:id"`
	CourseID int  `gorm:"column:course_id"`
	PaperID  *int `gorm:"column:paper_id"`
}

type PaperGradePaper struct {
	ID        int  `gorm:"column:id"`
	CourseID  *int `gorm:"column:course_id"`
	ChapterID *int `gorm:"column:chapter_id"`
	PassScore int  `gorm:"column:pass_score"`
}

type QuestionLogRow struct {
	ID          int    `gorm:"column:id"`
	QuestionID  int    `gorm:"column:question_id"`
	Score       int    `gorm:"column:score"`
	GradeRemark string `gorm:"column:grade_remark"`
}

type CourseCredit struct {
	Title  string `gorm:"column:title"`
	Credit int    `gorm:"column:credit"`
}

type PaperSaveInput struct {
	ID         int
	CourseID   int
	ChapterID  int
	IsExam     bool
	Score      int
	PassScore  int
	IsSkipMark int
	CreatedID  int
	Questions  []PaperQuestionInput
}

type PaperQuestionInput struct {
	Type    string
	Title   string
	Answer  json.RawMessage
	Score   int
	Options json.RawMessage
}

type AdminPaperRepository interface {
	CountCourse(ctx context.Context, courseID int) (int64, error)
	FindChapterRef(ctx context.Context, chapterID int) (*PaperChapterRef, error)
	SavePaper(ctx context.Context, in PaperSaveInput) error
	CountCourseUsers(ctx context.Context, courseID int) (int64, error)
	DeletePaperCascade(ctx context.Context, paperID, courseID int, chapterID *int) error
	ListUserPapers(ctx context.Context, f UserPaperListFilter) ([]AdminUserPaperRow, int64, error)
	ListPaperLogs(ctx context.Context, f UserPaperLogFilter) ([]AdminUserPaperRow, int64, error)
	ListUsersBrief(ctx context.Context, ids []int) ([]UserNickBrief, error)
	ListPapersBrief(ctx context.Context, ids []int) ([]PaperBrief, error)
	ListCoursesThumb(ctx context.Context, ids []int) ([]CourseThumb, error)
	ListChaptersTitle(ctx context.Context, ids []int) ([]ChapterTitle, error)
	FindUserPaper(ctx context.Context, id int) (*AdminUserPaperRow, error)
	ListQuestionLogsMap(ctx context.Context, userPaperID int) ([]map[string]any, error)
	FindUserPaperForGrade(ctx context.Context, id int) (idOut, userID, paperID, score, isPass int, err error)
	FindPaperForGrade(ctx context.Context, paperID int) (*PaperGradePaper, error)
	ListQuestionLogs(ctx context.Context, userPaperID int) ([]QuestionLogRow, error)
	UpdateQuestionLog(ctx context.Context, id int, updates map[string]any) error
	UpdateUserPaperGrade(ctx context.Context, id, score, isPass int) error
	FindUserChapter(ctx context.Context, userID, chapterID int) (*UserChapterRef, error)
	MarkUserChapterEnd(ctx context.Context, id int) error
	IncHasEndChapter(ctx context.Context, userID, courseID int) error
	ChapterTitle(ctx context.Context, id int) (string, error)
	MarkUserCourseEnd(ctx context.Context, userID, courseID int) error
	FindCourseCredit(ctx context.Context, courseID int) (*CourseCredit, error)
	CoursePaperID(ctx context.Context, courseID int) (*int, error)
}

type UserChapterRef struct {
	ID       int  `gorm:"column:id"`
	CourseID *int `gorm:"column:course_id"`
}

type adminPaperRepo struct{ db *gorm.DB }

func NewAdminPaper(db *gorm.DB) AdminPaperRepository { return &adminPaperRepo{db: db} }

func (r *adminPaperRepo) CountCourse(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("courses").Where("id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n, err
}

func (r *adminPaperRepo) FindChapterRef(ctx context.Context, chapterID int) (*PaperChapterRef, error) {
	var row PaperChapterRef
	err := r.db.WithContext(ctx).Table("course_chapters").Where("id = ? AND deleted_at IS NULL", chapterID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

type paperInsertRow struct {
	ID         int            `gorm:"column:id;primaryKey"`
	CourseID   *int           `gorm:"column:course_id"`
	ChapterID  *int           `gorm:"column:chapter_id"`
	Score      int            `gorm:"column:score"`
	PassScore  int            `gorm:"column:pass_score"`
	CreatedID  *int           `gorm:"column:created_id"`
	IsSkipMark int            `gorm:"column:is_skip_mark"`
	Type       string         `gorm:"column:type"`
	CreatedAt  *time.Time     `gorm:"column:created_at"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (paperInsertRow) TableName() string { return "papers" }

func (r *adminPaperRepo) SavePaper(ctx context.Context, in PaperSaveInput) error {
	now := time.Now()
	paperType := "chapter"
	if in.IsExam {
		paperType = "course"
	}
	data := map[string]any{
		"course_id": in.CourseID, "score": in.Score, "pass_score": in.PassScore,
		"is_skip_mark": in.IsSkipMark, "type": paperType, "updated_at": now,
	}
	if in.IsExam {
		data["chapter_id"] = nil
	} else {
		data["chapter_id"] = in.ChapterID
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		id := in.ID
		if id > 0 {
			if in.IsExam {
				var coursePaperID *int
				if err := tx.Table("courses").Select("paper_id").Where("id = ?", in.CourseID).Scan(&coursePaperID).Error; err != nil {
					return err
				}
				if coursePaperID == nil || *coursePaperID != id {
					return ErrPaperCourseMismatch
				}
			}
			var count int64
			if err := tx.Model(&paperInsertRow{}).Where("id = ?", id).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return ErrPaperNotFound
			}
			if err := tx.Model(&paperInsertRow{}).Where("id = ?", id).Updates(data).Error; err != nil {
				return err
			}
			if err := tx.Where("paper_id = ?", id).Delete(&questionInsertRow{}).Error; err != nil {
				return err
			}
			return insertPaperQuestions(tx, id, in.Questions, now)
		}
		courseID := in.CourseID
		createdID := in.CreatedID
		row := paperInsertRow{
			CourseID: &courseID, Score: in.Score, PassScore: in.PassScore,
			CreatedID: &createdID, IsSkipMark: in.IsSkipMark, Type: paperType,
			CreatedAt: &now, UpdatedAt: &now,
		}
		if !in.IsExam {
			ch := in.ChapterID
			row.ChapterID = &ch
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		id = row.ID
		if err := insertPaperQuestions(tx, id, in.Questions, now); err != nil {
			return err
		}
		if in.IsExam {
			return tx.Table("courses").Where("id = ?", in.CourseID).Update("paper_id", id).Error
		}
		return tx.Table("course_chapters").Where("id = ?", in.ChapterID).Update("paper_id", id).Error
	})
}

type questionInsertRow struct {
	ID        int64           `gorm:"column:id;primaryKey"`
	PaperID   int             `gorm:"column:paper_id"`
	Sort      int             `gorm:"column:sort"`
	Type      string          `gorm:"column:type"`
	Title     string          `gorm:"column:title"`
	Answer    json.RawMessage `gorm:"column:answer"`
	Score     int             `gorm:"column:score"`
	Options   json.RawMessage `gorm:"column:options"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at"`
}

func (questionInsertRow) TableName() string { return "paper_questions" }

func insertPaperQuestions(tx *gorm.DB, paperID int, questions []PaperQuestionInput, now time.Time) error {
	for index, q := range questions {
		options := q.Options
		if len(options) == 0 || string(options) == "null" {
			options = []byte("null")
		}
		row := questionInsertRow{
			PaperID: paperID, Sort: index, Type: q.Type, Title: q.Title,
			Answer: q.Answer, Score: q.Score, Options: options, CreatedAt: &now, UpdatedAt: &now,
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *adminPaperRepo) CountCourseUsers(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course").Where("course_id = ?", courseID).Count(&n).Error
	return n, err
}

func (r *adminPaperRepo) DeletePaperCascade(ctx context.Context, paperID, courseID int, chapterID *int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if chapterID != nil && *chapterID > 0 {
			if err := tx.Table("course_chapters").Where("id = ?", *chapterID).Update("paper_id", nil).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Table("courses").Where("id = ?", courseID).Update("paper_id", nil).Error; err != nil {
				return err
			}
		}
		var userPaperIDs []int
		if err := tx.Table("user_paper").Where("paper_id = ?", paperID).Pluck("id", &userPaperIDs).Error; err != nil {
			return err
		}
		if len(userPaperIDs) > 0 {
			if err := tx.Exec("DELETE FROM user_paper_question_log WHERE user_paper_id IN ?", userPaperIDs).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_paper WHERE id IN ?", userPaperIDs).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec("DELETE FROM paper_questions WHERE paper_id = ?", paperID).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Table("papers").Where("id = ?", paperID).Update("deleted_at", now).Error
	})
}

func (r *adminPaperRepo) ListUserPapers(ctx context.Context, f UserPaperListFilter) ([]AdminUserPaperRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_paper").
		Joins("JOIN papers ON papers.id = user_paper.paper_id AND papers.deleted_at IS NULL").
		Where("user_paper.deleted_at IS NULL AND papers.type = ?", f.PaperType)
	if f.Exam {
		query = query.Where("papers.chapter_id IS NULL")
	}
	if f.CourseID > 0 {
		query = query.Where("papers.course_id = ?", f.CourseID)
	}
	if f.ChapterID > 0 {
		query = query.Where("papers.chapter_id = ?", f.ChapterID)
	}
	if f.Title != "" {
		query = query.Where("papers.course_id IN (SELECT id FROM courses WHERE title LIKE ? AND deleted_at IS NULL)", "%"+f.Title+"%")
	}
	if f.IsPass != nil {
		query = query.Where("user_paper.is_pass = ?", *f.IsPass)
	}
	if f.IsPeopleCheck != nil {
		query = query.Where("user_paper.is_people_check = ?", *f.IsPeopleCheck)
	}
	if f.NickName != "" {
		query = query.Where("user_paper.user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Limited {
		query = query.Where("user_paper.user_id IN (SELECT id FROM users WHERE city_id = ?)", f.CityID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []AdminUserPaperRow
	err := query.Select("user_paper.id, user_paper.user_id, user_paper.paper_id, user_paper.score, user_paper.is_pass, user_paper.is_people_check, user_paper.submit_count, user_paper.created_at").
		Order("user_paper." + f.Sort + " " + f.Order).
		Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminPaperRepo) ListPaperLogs(ctx context.Context, f UserPaperLogFilter) ([]AdminUserPaperRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_paper").Where("paper_id = ? AND deleted_at IS NULL", f.PaperID)
	if f.NickName != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Limited {
		query = query.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", f.CityID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []AdminUserPaperRow
	err := query.Select("id, user_id, paper_id, score, submit_count, created_at, is_pass, is_people_check").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminPaperRepo) ListUsersBrief(ctx context.Context, ids []int) ([]UserNickBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []UserNickBrief
	err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, avatar, credit").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPaperRepo) ListPapersBrief(ctx context.Context, ids []int) ([]PaperBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []PaperBrief
	err := r.db.WithContext(ctx).Table("papers").
		Select("id, course_id, chapter_id").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPaperRepo) ListCoursesThumb(ctx context.Context, ids []int) ([]CourseThumb, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []CourseThumb
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPaperRepo) ListChaptersTitle(ctx context.Context, ids []int) ([]ChapterTitle, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []ChapterTitle
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, title").Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error
	return rows, err
}

func (r *adminPaperRepo) FindUserPaper(ctx context.Context, id int) (*AdminUserPaperRow, error) {
	var row AdminUserPaperRow
	err := r.db.WithContext(ctx).Table("user_paper").Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPaperRepo) ListQuestionLogsMap(ctx context.Context, userPaperID int) ([]map[string]any, error) {
	var logs []map[string]any
	err := r.db.WithContext(ctx).Table("user_paper_question_log").
		Where("user_paper_id = ?", userPaperID).Find(&logs).Error
	return logs, err
}

func (r *adminPaperRepo) FindUserPaperForGrade(ctx context.Context, id int) (int, int, int, int, int, error) {
	var row struct {
		ID      int `gorm:"column:id"`
		UserID  int `gorm:"column:user_id"`
		PaperID int `gorm:"column:paper_id"`
		Score   int `gorm:"column:score"`
		IsPass  int `gorm:"column:is_pass"`
	}
	err := r.db.WithContext(ctx).Table("user_paper").Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	return row.ID, row.UserID, row.PaperID, row.Score, row.IsPass, err
}

func (r *adminPaperRepo) FindPaperForGrade(ctx context.Context, paperID int) (*PaperGradePaper, error) {
	var row PaperGradePaper
	err := r.db.WithContext(ctx).Table("papers").
		Select("id, course_id, chapter_id, pass_score").
		Where("id = ? AND deleted_at IS NULL", paperID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPaperRepo) ListQuestionLogs(ctx context.Context, userPaperID int) ([]QuestionLogRow, error) {
	var rows []QuestionLogRow
	err := r.db.WithContext(ctx).Table("user_paper_question_log").
		Where("user_paper_id = ?", userPaperID).Find(&rows).Error
	return rows, err
}

func (r *adminPaperRepo) UpdateQuestionLog(ctx context.Context, id int, updates map[string]any) error {
	return r.db.WithContext(ctx).Table("user_paper_question_log").Where("id = ?", id).Updates(updates).Error
}

func (r *adminPaperRepo) UpdateUserPaperGrade(ctx context.Context, id, score, isPass int) error {
	return r.db.WithContext(ctx).Table("user_paper").Where("id = ?", id).Updates(map[string]any{
		"score": score, "is_pass": isPass, "is_people_check": 0,
	}).Error
}

func (r *adminPaperRepo) FindUserChapter(ctx context.Context, userID, chapterID int) (*UserChapterRef, error) {
	var row UserChapterRef
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Select("id, course_id").
		Where("user_id = ? AND chapter_id = ?", userID, chapterID).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPaperRepo) MarkUserChapterEnd(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Table("user_course_chapter").Where("id = ?", id).
		Updates(map[string]any{"is_end": 1, "has_end_chapter": 1}).Error
}

func (r *adminPaperRepo) IncHasEndChapter(ctx context.Context, userID, courseID int) error {
	return r.db.WithContext(ctx).Exec(
		"UPDATE user_course SET has_end_chapter = has_end_chapter + 1 WHERE user_id = ? AND course_id = ?",
		userID, courseID,
	).Error
}

func (r *adminPaperRepo) ChapterTitle(ctx context.Context, id int) (string, error) {
	var title string
	err := r.db.WithContext(ctx).Table("course_chapters").Select("title").Where("id = ? AND deleted_at IS NULL", id).Scan(&title).Error
	return title, err
}

func (r *adminPaperRepo) MarkUserCourseEnd(ctx context.Context, userID, courseID int) error {
	return r.db.WithContext(ctx).Table("user_course").
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Update("is_end", 1).Error
}

func (r *adminPaperRepo) FindCourseCredit(ctx context.Context, courseID int) (*CourseCredit, error) {
	var row CourseCredit
	err := r.db.WithContext(ctx).Table("courses").Select("title, credit").Where("id = ? AND deleted_at IS NULL", courseID).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminPaperRepo) CoursePaperID(ctx context.Context, courseID int) (*int, error) {
	var paperID *int
	err := r.db.WithContext(ctx).Table("courses").Select("paper_id").Where("id = ?", courseID).Scan(&paperID).Error
	return paperID, err
}
