package repo

import (
	"context"

	"gorm.io/gorm"
)

// AdminHomeRepository covers admin dashboard stats / charts / ranks / verse.
type AdminHomeRepository interface {
	MaxVerseID(ctx context.Context) (int, error)
	VerseContentGE(ctx context.Context, startID int) (string, error)
	VerseContentLE(ctx context.Context, startID int) (string, error)

	ScopedCount(ctx context.Context, table, where string, cityID int, limited bool, args ...any) (int, error)
	ScopedSum(ctx context.Context, table, column, where string, cityID int, limited bool, args ...any) (int, error)

	CountByDate(ctx context.Context, table, dateColumn, start, end, extra string, cityID int, limited bool) (map[string]int, error)
	SumByDate(ctx context.Context, table, dateColumn, sumColumn, start, end, extra string, cityID int, limited bool) (map[string]int, error)

	CourseLearningRanks(ctx context.Context, cityID int, limited bool, limit int) ([]IDCountRow, error)
	CourseTitles(ctx context.Context, ids []int) (map[int]string, error)
	ChapterCompleteRanks(ctx context.Context, since string, cityID int, limited bool, limit int) ([]ChapterRankRow, error)
	ChapterTitles(ctx context.Context, ids []int) (map[int]string, error)
	CoinsRanks(ctx context.Context, since string, cityID int, limited bool, limit int) ([]IDCountRow, error)
	UserAvatars(ctx context.Context, ids []int) (map[int]UserAvatarRow, error)
}

type IDCountRow struct {
	ID    int `gorm:"column:id"`
	Count int `gorm:"column:c"`
}

type ChapterRankRow struct {
	ChapterID int `gorm:"column:chapter_id"`
	CourseID  int `gorm:"column:course_id"`
	Count     int `gorm:"column:c"`
}

type UserAvatarRow struct {
	Avatar   string
	NickName string
}

type adminHomeRepo struct{ db *gorm.DB }

func NewAdminHome(db *gorm.DB) AdminHomeRepository { return &adminHomeRepo{db: db} }

func (r *adminHomeRepo) MaxVerseID(ctx context.Context) (int, error) {
	var maxID int
	err := r.db.WithContext(ctx).Table("verse").Select("COALESCE(MAX(id), 0)").Scan(&maxID).Error
	return maxID, err
}

func (r *adminHomeRepo) VerseContentGE(ctx context.Context, startID int) (string, error) {
	var verse string
	err := r.db.WithContext(ctx).Table("verse").Select("content").
		Where("id >= ?", startID).Order("id ASC").Limit(1).Scan(&verse).Error
	return verse, err
}

func (r *adminHomeRepo) VerseContentLE(ctx context.Context, startID int) (string, error) {
	var verse string
	err := r.db.WithContext(ctx).Table("verse").Select("content").
		Where("id <= ?", startID).Order("id DESC").Limit(1).Scan(&verse).Error
	return verse, err
}

func (r *adminHomeRepo) ScopedCount(ctx context.Context, table, where string, cityID int, limited bool, args ...any) (int, error) {
	q := r.db.WithContext(ctx).Table(table)
	if where != "" {
		q = q.Where(where, args...)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return 0, err
	}
	return int(n), nil
}

func (r *adminHomeRepo) ScopedSum(ctx context.Context, table, column, where string, cityID int, limited bool, args ...any) (int, error) {
	q := r.db.WithContext(ctx).Table(table)
	if where != "" {
		q = q.Where(where, args...)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var n *float64
	if err := q.Select("COALESCE(SUM(`" + column + "`), 0)").Scan(&n).Error; err != nil {
		return 0, err
	}
	if n == nil {
		return 0, nil
	}
	return int(*n), nil
}

func (r *adminHomeRepo) CountByDate(ctx context.Context, table, dateColumn, start, end, extra string, cityID int, limited bool) (map[string]int, error) {
	q := r.db.WithContext(ctx).Table(table).Where(dateColumn+" BETWEEN ? AND ?", start, end)
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var rows []dateCountRow
	if err := q.Select("DATE(" + dateColumn + ") as d, COUNT(*) as c").Group("d").Find(&rows).Error; err != nil {
		return nil, err
	}
	return dateCountMap(rows), nil
}

func (r *adminHomeRepo) SumByDate(ctx context.Context, table, dateColumn, sumColumn, start, end, extra string, cityID int, limited bool) (map[string]int, error) {
	q := r.db.WithContext(ctx).Table(table).Where(dateColumn+" BETWEEN ? AND ?", start, end)
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var rows []dateCountRow
	if err := q.Select("DATE(" + dateColumn + ") as d, COALESCE(SUM(`" + sumColumn + "`), 0) as c").Group("d").Find(&rows).Error; err != nil {
		return nil, err
	}
	return dateCountMap(rows), nil
}

type dateCountRow struct {
	D string `gorm:"column:d"`
	C int    `gorm:"column:c"`
}

func dateCountMap(rows []dateCountRow) map[string]int {
	out := map[string]int{}
	for _, r := range rows {
		day := r.D
		if len(day) > 10 {
			day = day[:10]
		}
		out[day] = r.C
	}
	return out
}

func (r *adminHomeRepo) CourseLearningRanks(ctx context.Context, cityID int, limited bool, limit int) ([]IDCountRow, error) {
	q := r.db.WithContext(ctx).Table("user_course").Where("is_end = 0")
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	type row struct {
		CourseID    int `gorm:"column:course_id"`
		EnrollCount int `gorm:"column:enroll_count"`
	}
	var rows []row
	if err := q.Select("course_id, COUNT(*) as enroll_count").Group("course_id").
		Order("enroll_count DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]IDCountRow, 0, len(rows))
	for _, item := range rows {
		out = append(out, IDCountRow{ID: item.CourseID, Count: item.EnrollCount})
	}
	return out, nil
}

func (r *adminHomeRepo) CourseTitles(ctx context.Context, ids []int) (map[int]string, error) {
	return r.idTitleMap(ctx, "courses", ids)
}

func (r *adminHomeRepo) ChapterCompleteRanks(ctx context.Context, since string, cityID int, limited bool, limit int) ([]ChapterRankRow, error) {
	q := r.db.WithContext(ctx).Table("user_course_chapter AS ucc").
		Where("ucc.is_end = 1 AND ucc.updated_at >= ?", since)
	if limited {
		q = q.Where("ucc.user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var rows []ChapterRankRow
	err := q.Select("ucc.chapter_id, ucc.course_id, COUNT(*) as c").
		Group("ucc.chapter_id, ucc.course_id").Order("c DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *adminHomeRepo) ChapterTitles(ctx context.Context, ids []int) (map[int]string, error) {
	return r.idTitleMap(ctx, "course_chapters", ids)
}

func (r *adminHomeRepo) CoinsRanks(ctx context.Context, since string, cityID int, limited bool, limit int) ([]IDCountRow, error) {
	q := r.db.WithContext(ctx).Table("user_coins_records").
		Where("deleted_at IS NULL AND `sum` > 0 AND created_at >= ?", since)
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	type row struct {
		UserID   int `gorm:"column:user_id"`
		CoinsSum int `gorm:"column:coins_sum"`
	}
	var rows []row
	if err := q.Select("user_id, SUM(`sum`) as coins_sum").Group("user_id").
		Order("coins_sum DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]IDCountRow, 0, len(rows))
	for _, item := range rows {
		out = append(out, IDCountRow{ID: item.UserID, Count: item.CoinsSum})
	}
	return out, nil
}

func (r *adminHomeRepo) UserAvatars(ctx context.Context, ids []int) (map[int]UserAvatarRow, error) {
	out := map[int]UserAvatarRow{}
	if len(ids) == 0 {
		return out, nil
	}
	type urow struct {
		ID       int     `gorm:"column:id"`
		Avatar   *string `gorm:"column:avatar"`
		NickName *string `gorm:"column:nick_name"`
	}
	var urows []urow
	if err := r.db.WithContext(ctx).Table("users").
		Select("id, avatar, nick_name").Where("id IN ?", ids).Find(&urows).Error; err != nil {
		return nil, err
	}
	for _, u := range urows {
		av, nick := "", ""
		if u.Avatar != nil {
			av = *u.Avatar
		}
		if u.NickName != nil {
			nick = *u.NickName
		}
		out[u.ID] = UserAvatarRow{Avatar: av, NickName: nick}
	}
	return out, nil
}

func (r *adminHomeRepo) idTitleMap(ctx context.Context, table string, ids []int) (map[int]string, error) {
	out := map[int]string{}
	if len(ids) == 0 {
		return out, nil
	}
	type trow struct {
		ID    int    `gorm:"column:id"`
		Title string `gorm:"column:title"`
	}
	var trows []trow
	if err := r.db.WithContext(ctx).Table(table).
		Select("id, title").Where("id IN ? AND deleted_at IS NULL", ids).Find(&trows).Error; err != nil {
		return nil, err
	}
	for _, t := range trows {
		out[t.ID] = t.Title
	}
	return out, nil
}
