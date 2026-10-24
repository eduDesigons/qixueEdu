package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type MediaBrief struct {
	ID    int64   `gorm:"column:id"`
	Title string  `gorm:"column:title"`
	Thumb string  `gorm:"column:thumb"`
	URL   *string `gorm:"column:url"`
	Type  string  `gorm:"column:type"`
	Size  int64   `gorm:"column:size"`
}

type ChapterUserRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	ChapterID int        `gorm:"column:chapter_id"`
	CourseID  int        `gorm:"column:course_id"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	UnLock    int        `gorm:"column:un_lock"`
	IsEnd     int        `gorm:"column:is_end"`
}

type ChapterUserListFilter struct {
	ChapterID int
	CityID    int
	NickName  string
	Mobile    string
	Offset    int
	Limit     int
}

type UserChapterUnlockRow struct {
	ID        int `gorm:"column:id"`
	ChapterID int `gorm:"column:chapter_id"`
	UnLock    int `gorm:"column:un_lock"`
}

type ChapterLimitDay struct {
	ID       int  `gorm:"column:id"`
	LimitDay *int `gorm:"column:limit_day"`
}

type AdminChapterRepository interface {
	CommentCountsByChapterIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	UserCountsByChapterIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	EndedUserCountsByChapterIDs(ctx context.Context, ids []int) ([]CountByInt, error)
	FindAdmin(ctx context.Context, id int) (*AdminNick, error)
	FindPaperScore(ctx context.Context, id int) (*PaperScoreBrief, error)
	FindMedia(ctx context.Context, id int) (*MediaBrief, error)
	CountComments(ctx context.Context, chapterID int) (int64, error)
	CountUsers(ctx context.Context, chapterID int) (int64, error)
	CountEndedUsers(ctx context.Context, chapterID int) (int64, error)
	CountCourse(ctx context.Context, courseID int) (int64, error)
	IncCourseChapterCount(ctx context.Context, courseID int) error
	CreateChapterAndInc(ctx context.Context, row any, courseID int) error
	DeleteChapterCascade(ctx context.Context, chapterID, courseID int, paperID *int) error
	ListChapterUsers(ctx context.Context, f ChapterUserListFilter) ([]ChapterUserRow, int64, error)
	ListUsersFull(ctx context.Context, ids []int) ([]UserBriefFull, error)
	FindUserChapter(ctx context.Context, id int) (*UserChapterUnlockRow, error)
	FindChapterLimitDay(ctx context.Context, id int) (*ChapterLimitDay, error)
	UnlockUserChapter(ctx context.Context, id int, expire time.Time) error
}

type adminChapterRepo struct{ db *gorm.DB }

func NewAdminChapter(db *gorm.DB) AdminChapterRepository { return &adminChapterRepo{db: db} }

func (r *adminChapterRepo) CommentCountsByChapterIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		ChapterID int `gorm:"column:chapter_id"`
		Total     int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("course_chapter_comments").
		Select("chapter_id, COUNT(*) AS total").
		Where("chapter_id IN ? AND deleted_at IS NULL", ids).
		Group("chapter_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.ChapterID, Total: row.Total})
	}
	return out, err
}

func (r *adminChapterRepo) UserCountsByChapterIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		ChapterID int `gorm:"column:chapter_id"`
		Total     int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Select("chapter_id, COUNT(*) AS total").
		Where("chapter_id IN ?", ids).
		Group("chapter_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.ChapterID, Total: row.Total})
	}
	return out, err
}

func (r *adminChapterRepo) EndedUserCountsByChapterIDs(ctx context.Context, ids []int) ([]CountByInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []struct {
		ChapterID int `gorm:"column:chapter_id"`
		Total     int `gorm:"column:total"`
	}
	err := r.db.WithContext(ctx).Table("user_course_chapter").
		Select("chapter_id, COUNT(*) AS total").
		Where("chapter_id IN ? AND is_end = 1", ids).
		Group("chapter_id").Find(&rows).Error
	out := make([]CountByInt, 0, len(rows))
	for _, row := range rows {
		out = append(out, CountByInt{Key: row.ChapterID, Total: row.Total})
	}
	return out, err
}

func (r *adminChapterRepo) FindAdmin(ctx context.Context, id int) (*AdminNick, error) {
	var row AdminNick
	err := r.db.WithContext(ctx).Table("administrator").Select("id, nickName").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminChapterRepo) FindPaperScore(ctx context.Context, id int) (*PaperScoreBrief, error) {
	var row PaperScoreBrief
	err := r.db.WithContext(ctx).Table("papers").
		Select("id, score, pass_score").Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminChapterRepo) FindMedia(ctx context.Context, id int) (*MediaBrief, error) {
	var row MediaBrief
	err := r.db.WithContext(ctx).Table("media_file").
		Select("id, title, thumb, url, type, size").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminChapterRepo) CountComments(ctx context.Context, chapterID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("course_chapter_comments").
		Where("chapter_id = ? AND deleted_at IS NULL", chapterID).Count(&n).Error
	return n, err
}

func (r *adminChapterRepo) CountUsers(ctx context.Context, chapterID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course_chapter").Where("chapter_id = ?", chapterID).Count(&n).Error
	return n, err
}

func (r *adminChapterRepo) CountEndedUsers(ctx context.Context, chapterID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("user_course_chapter").Where("chapter_id = ? AND is_end = 1", chapterID).Count(&n).Error
	return n, err
}

func (r *adminChapterRepo) CountCourse(ctx context.Context, courseID int) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Table("courses").Where("id = ? AND deleted_at IS NULL", courseID).Count(&n).Error
	return n, err
}

func (r *adminChapterRepo) IncCourseChapterCount(ctx context.Context, courseID int) error {
	return r.db.WithContext(ctx).Table("courses").Where("id = ?", courseID).
		UpdateColumn("chapter_count", gorm.Expr("chapter_count + 1")).Error
}

func (r *adminChapterRepo) CreateChapterAndInc(ctx context.Context, row any, courseID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		return tx.Table("courses").Where("id = ?", courseID).
			UpdateColumn("chapter_count", gorm.Expr("chapter_count + 1")).Error
	})
}

func (r *adminChapterRepo) DeleteChapterCascade(ctx context.Context, chapterID, courseID int, paperID *int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Table("course_chapter_comments").
			Where("chapter_id = ? AND deleted_at IS NULL", chapterID).
			Update("deleted_at", now).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_like_chapters WHERE chapter_id = ?", chapterID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_course_chapter WHERE chapter_id = ?", chapterID).Error; err != nil {
			return err
		}
		if paperID != nil && *paperID > 0 {
			var userPaperIDs []int64
			if err := tx.Table("user_paper").Where("paper_id = ?", *paperID).
				Pluck("id", &userPaperIDs).Error; err != nil {
				return err
			}
			if len(userPaperIDs) > 0 {
				if err := tx.Exec("DELETE FROM user_paper_question_log WHERE user_paper_id IN ?", userPaperIDs).Error; err != nil {
					return err
				}
			}
			if err := tx.Exec("DELETE FROM user_paper WHERE paper_id = ?", *paperID).Error; err != nil {
				return err
			}
		}
		if err := tx.Table("courses").Where("id = ?", courseID).
			UpdateColumn("chapter_count", gorm.Expr("GREATEST(chapter_count - 1, 0)")).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM course_chapters WHERE id = ?", chapterID).Error; err != nil {
			return err
		}
		var left int64
		if err := tx.Table("course_chapters").Where("course_id = ?", courseID).Count(&left).Error; err != nil {
			return err
		}
		if left == 0 {
			return tx.Table("courses").Where("id = ?", courseID).
				Updates(map[string]any{"is_show": 0, "chapter_count": 0}).Error
		}
		return nil
	})
}

func (r *adminChapterRepo) ListChapterUsers(ctx context.Context, f ChapterUserListFilter) ([]ChapterUserRow, int64, error) {
	query := r.db.WithContext(ctx).Table("user_course_chapter").Where("chapter_id = ?", f.ChapterID)
	query = applyUserCityFilter(query, f.CityID)
	if f.NickName != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Mobile != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE mobile = ?)", f.Mobile)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []ChapterUserRow
	err := query.Select("id, user_id, chapter_id, course_id, created_at, updated_at, un_lock, is_end").
		Order("id ASC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminChapterRepo) ListUsersFull(ctx context.Context, ids []int) ([]UserBriefFull, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []UserBriefFull
	err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, credit, avatar, mobile").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}

func (r *adminChapterRepo) FindUserChapter(ctx context.Context, id int) (*UserChapterUnlockRow, error) {
	var row UserChapterUnlockRow
	err := r.db.WithContext(ctx).Table("user_course_chapter").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminChapterRepo) FindChapterLimitDay(ctx context.Context, id int) (*ChapterLimitDay, error) {
	var row ChapterLimitDay
	err := r.db.WithContext(ctx).Table("course_chapters").
		Select("id, limit_day").Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *adminChapterRepo) UnlockUserChapter(ctx context.Context, id int, expire time.Time) error {
	return r.db.WithContext(ctx).Table("user_course_chapter").Where("id = ?", id).
		Updates(map[string]any{"expire_time": expire, "un_lock": 1}).Error
}
