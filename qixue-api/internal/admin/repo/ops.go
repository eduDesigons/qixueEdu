package admin

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// AdminOpsRepository covers ops enroll/learn/amount lists and chart aggregates.
type AdminOpsRepository interface {
	ListCourseEnrolls(ctx context.Context, f OpsListFilter) ([]OpsCourseEnrollRow, int64, error)
	ListPathEnrolls(ctx context.Context, f OpsListFilter) ([]OpsPathEnrollRow, int64, error)
	ListLearns(ctx context.Context, f OpsListFilter) ([]OpsLearnRow, int64, error)
	ListAmountRecords(ctx context.Context, table string, f OpsAmountFilter) ([]OpsAmountRow, int64, error)

	LoadUsers(ctx context.Context, ids []int) (map[int]OpsUserBrief, error)
	LoadTitleThumbs(ctx context.Context, table string, ids []int) (map[int]OpsTitleThumb, error)
	LoadTitles(ctx context.Context, table string, ids []int) (map[int]OpsTitleOnly, error)

	ScopedCount(ctx context.Context, table string, cityID int, limited bool, extra string) (int, error)
	ScopedSum(ctx context.Context, table, column string, cityID int, limited bool, extra string) (int, error)
	CountSince(ctx context.Context, table, column, from string, cityID int, limited bool, extra string) (int, error)
	SumSince(ctx context.Context, table, dateCol, sumCol, from string, cityID int, limited bool, extra string) (int, error)
	CountByDate(ctx context.Context, table, column string, dates []string, cityID int, limited bool, extra string) ([]int, error)
	SumByDate(ctx context.Context, table, column, sumColumn string, dates []string, cityID int, limited bool, extra string) ([]int, error)
	CountByExpr(ctx context.Context, table, expr string, dates []string, cityID int, limited bool, extra string) ([]int, error)
	CountExprSince(ctx context.Context, table, expr, day string, cityID int, limited bool, extra string) (int, error)
	TopCount(ctx context.Context, table, idColumn, titleTable string, cityID int, limited bool, extra string) ([]OpsNamedValue, error)
	TopText(ctx context.Context, table, column string, cityID int, limited bool, extra string) ([]OpsNamedValue, error)
	TopAmount(ctx context.Context, table string, cityID int, limited bool, extra string) ([]OpsNamedValue, error)
}

type OpsListFilter struct {
	NickName string
	Title    string
	IsEnd    *int
	CityID   int
	Limited  bool
	Offset   int
	Limit    int
}

type OpsAmountFilter struct {
	NickName string
	Field    string
	CityID   int
	Limited  bool
	Offset   int
	Limit    int
}

type OpsCourseEnrollRow struct {
	ID            int        `gorm:"column:id"`
	UserID        int        `gorm:"column:user_id"`
	CourseID      int        `gorm:"column:course_id"`
	HasEndChapter int        `gorm:"column:has_end_chapter"`
	IsEnd         int        `gorm:"column:is_end"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

type OpsPathEnrollRow struct {
	ID           int        `gorm:"column:id"`
	UserID       int        `gorm:"column:user_id"`
	PathID       int        `gorm:"column:path_id"`
	HasEndCourse int        `gorm:"column:has_end_course"`
	IsEnd        int        `gorm:"column:is_end"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

type OpsLearnRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	CourseID  int        `gorm:"column:course_id"`
	ChapterID int        `gorm:"column:chapter_id"`
	IsEnd     int        `gorm:"column:is_end"`
	UnLock    int        `gorm:"column:un_lock"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type OpsAmountRow struct {
	ID        int        `gorm:"column:id"`
	UserID    int        `gorm:"column:user_id"`
	Field     string     `gorm:"column:field"`
	Sum       int        `gorm:"column:sum"`
	Remark    string     `gorm:"column:remark"`
	CreatedAt *time.Time `gorm:"column:created_at"`
}

type OpsUserBrief struct {
	ID       int     `json:"id"`
	NickName *string `json:"nick_name"`
	Avatar   string  `json:"avatar"`
	Mobile   string  `json:"mobile"`
	CityID   *int    `json:"city_id"`
}

type OpsTitleThumb struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Thumb string `json:"thumb"`
}

type OpsTitleOnly struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type OpsNamedValue struct {
	Name  string
	Value int
}

type adminOpsRepo struct{ db *gorm.DB }

func NewAdminOps(db *gorm.DB) AdminOpsRepository { return &adminOpsRepo{db: db} }

func (r *adminOpsRepo) ListCourseEnrolls(ctx context.Context, f OpsListFilter) ([]OpsCourseEnrollRow, int64, error) {
	q := r.db.WithContext(ctx).Table("user_course")
	if f.NickName != "" {
		q = q.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Title != "" {
		q = q.Where("course_id IN (SELECT id FROM courses WHERE title LIKE ? AND deleted_at IS NULL)", "%"+f.Title+"%")
	}
	if f.IsEnd != nil {
		q = q.Where("is_end = ?", *f.IsEnd)
	}
	if f.Limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", f.CityID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []OpsCourseEnrollRow
	err := q.Select("id, user_id, course_id, has_end_chapter, is_end, created_at, updated_at").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminOpsRepo) ListPathEnrolls(ctx context.Context, f OpsListFilter) ([]OpsPathEnrollRow, int64, error) {
	q := r.db.WithContext(ctx).Table("user_path")
	if f.NickName != "" {
		q = q.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Title != "" {
		q = q.Where("path_id IN (SELECT id FROM paths WHERE title LIKE ? AND deleted_at IS NULL)", "%"+f.Title+"%")
	}
	if f.IsEnd != nil {
		q = q.Where("is_end = ?", *f.IsEnd)
	}
	if f.Limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", f.CityID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []OpsPathEnrollRow
	err := q.Select("id, user_id, path_id, has_end_course, is_end, created_at, updated_at").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminOpsRepo) ListLearns(ctx context.Context, f OpsListFilter) ([]OpsLearnRow, int64, error) {
	q := r.db.WithContext(ctx).Table("user_course_chapter")
	if f.NickName != "" {
		q = q.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Title != "" {
		like := "%" + f.Title + "%"
		q = q.Where(`(
			course_id IN (SELECT id FROM courses WHERE title LIKE ? AND deleted_at IS NULL)
			OR chapter_id IN (SELECT id FROM course_chapters WHERE title LIKE ? AND deleted_at IS NULL)
		)`, like, like)
	}
	if f.IsEnd != nil {
		q = q.Where("is_end = ?", *f.IsEnd)
	}
	if f.Limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", f.CityID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []OpsLearnRow
	err := q.Select("id, user_id, course_id, chapter_id, is_end, un_lock, created_at, updated_at").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminOpsRepo) ListAmountRecords(ctx context.Context, table string, f OpsAmountFilter) ([]OpsAmountRow, int64, error) {
	q := r.db.WithContext(ctx).Table(table).Where("deleted_at IS NULL")
	if f.NickName != "" {
		q = q.Where("user_id IN (SELECT id FROM users WHERE nick_name LIKE ?)", "%"+f.NickName+"%")
	}
	if f.Field != "" {
		q = q.Where("field LIKE ?", "%"+f.Field+"%")
	}
	if f.Limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", f.CityID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []OpsAmountRow
	err := q.Select("id, user_id, field, `sum`, remark, created_at").
		Order("id DESC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error
	return rows, total, err
}

func (r *adminOpsRepo) LoadUsers(ctx context.Context, ids []int) (map[int]OpsUserBrief, error) {
	out := map[int]OpsUserBrief{}
	if len(ids) == 0 {
		return out, nil
	}
	type row struct {
		ID       int     `gorm:"column:id"`
		NickName *string `gorm:"column:nick_name"`
		Avatar   string  `gorm:"column:avatar"`
		Mobile   string  `gorm:"column:mobile"`
		CityID   *int    `gorm:"column:city_id"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).Table("users").
		Select("id, nick_name, avatar, mobile, city_id").Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, item := range rows {
		out[item.ID] = OpsUserBrief{
			ID: item.ID, NickName: item.NickName, Avatar: item.Avatar, Mobile: item.Mobile, CityID: item.CityID,
		}
	}
	return out, nil
}

func (r *adminOpsRepo) LoadTitleThumbs(ctx context.Context, table string, ids []int) (map[int]OpsTitleThumb, error) {
	out := map[int]OpsTitleThumb{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []OpsTitleThumb
	q := r.db.WithContext(ctx).Table(table).Select("id, title, thumb").Where("id IN ?", ids)
	if table == "courses" || table == "paths" {
		q = q.Where("deleted_at IS NULL")
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		cp := rows[i]
		out[cp.ID] = cp
	}
	return out, nil
}

func (r *adminOpsRepo) LoadTitles(ctx context.Context, table string, ids []int) (map[int]OpsTitleOnly, error) {
	out := map[int]OpsTitleOnly{}
	if len(ids) == 0 {
		return out, nil
	}
	var rows []OpsTitleOnly
	q := r.db.WithContext(ctx).Table(table).Select("id, title").Where("id IN ?", ids)
	if table == "courses" || table == "course_chapters" || table == "paths" {
		q = q.Where("deleted_at IS NULL")
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		cp := rows[i]
		out[cp.ID] = cp
	}
	return out, nil
}

func (r *adminOpsRepo) ScopedCount(ctx context.Context, table string, cityID int, limited bool, extra string) (int, error) {
	q := r.db.WithContext(ctx).Table(table)
	if extra != "" {
		q = q.Where(extra)
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

func (r *adminOpsRepo) ScopedSum(ctx context.Context, table, column string, cityID int, limited bool, extra string) (int, error) {
	q := r.db.WithContext(ctx).Table(table)
	if extra != "" {
		q = q.Where(extra)
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

func (r *adminOpsRepo) CountSince(ctx context.Context, table, column, from string, cityID int, limited bool, extra string) (int, error) {
	q := r.db.WithContext(ctx).Table(table).Where(column+" >= ?", from)
	if extra != "" {
		q = q.Where(extra)
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

func (r *adminOpsRepo) SumSince(ctx context.Context, table, dateCol, sumCol, from string, cityID int, limited bool, extra string) (int, error) {
	q := r.db.WithContext(ctx).Table(table).Where(dateCol+" >= ?", from)
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var n *float64
	if err := q.Select("COALESCE(SUM(`" + sumCol + "`), 0)").Scan(&n).Error; err != nil {
		return 0, err
	}
	if n == nil {
		return 0, nil
	}
	return int(*n), nil
}

func (r *adminOpsRepo) CountByDate(ctx context.Context, table, column string, dates []string, cityID int, limited bool, extra string) ([]int, error) {
	return r.series(ctx, table, "DATE(`"+column+"`)", column, dates, "COUNT(*)", cityID, limited, extra)
}

func (r *adminOpsRepo) SumByDate(ctx context.Context, table, column, sumColumn string, dates []string, cityID int, limited bool, extra string) ([]int, error) {
	return r.series(ctx, table, "DATE(`"+column+"`)", column, dates, "COALESCE(SUM(`"+sumColumn+"`), 0)", cityID, limited, extra)
}

func (r *adminOpsRepo) CountByExpr(ctx context.Context, table, expr string, dates []string, cityID int, limited bool, extra string) ([]int, error) {
	start, end := dates[0], dates[len(dates)-1]
	q := r.db.WithContext(ctx).Table(table).Where(expr+" BETWEEN ? AND ?", start, end)
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var rows []dateCountRow
	if err := q.Select(expr + " as d, COUNT(*) as c").Group(expr).Find(&rows).Error; err != nil {
		return nil, err
	}
	return alignDates(dates, dateCountMap(rows)), nil
}

func (r *adminOpsRepo) CountExprSince(ctx context.Context, table, expr, day string, cityID int, limited bool, extra string) (int, error) {
	q := r.db.WithContext(ctx).Table(table).Where(expr+" = ?", day)
	if extra != "" {
		q = q.Where(extra)
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

func (r *adminOpsRepo) series(ctx context.Context, table, expr, column string, dates []string, aggregate string, cityID int, limited bool, extra string) ([]int, error) {
	start := dates[0] + " 00:00:00"
	end := dates[len(dates)-1] + " 23:59:59"
	q := r.db.WithContext(ctx).Table(table).Where(column+" BETWEEN ? AND ?", start, end)
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	var rows []dateCountRow
	if err := q.Select(expr + " as d, " + aggregate + " as c").Group(expr).Find(&rows).Error; err != nil {
		return nil, err
	}
	return alignDates(dates, dateCountMap(rows)), nil
}

func alignDates(dates []string, mapped map[string]int) []int {
	out := make([]int, len(dates))
	for i, d := range dates {
		out[i] = mapped[d]
	}
	return out
}

func (r *adminOpsRepo) TopCount(ctx context.Context, table, idColumn, titleTable string, cityID int, limited bool, extra string) ([]OpsNamedValue, error) {
	q := r.db.WithContext(ctx).Table(table).Where(idColumn + " > 0")
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	type row struct {
		GID int `gorm:"column:gid"`
		C   int `gorm:"column:c"`
	}
	var rows []row
	if err := q.Select("`" + idColumn + "` as gid, COUNT(*) as c").Group(idColumn).Order("c DESC").Limit(8).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(rows))
	for _, item := range rows {
		if item.GID > 0 {
			ids = append(ids, item.GID)
		}
	}
	titles := map[int]string{}
	if len(ids) > 0 {
		type trow struct {
			ID    int    `gorm:"column:id"`
			Title string `gorm:"column:title"`
		}
		var trows []trow
		tq := r.db.WithContext(ctx).Table(titleTable).Select("id, title").Where("id IN ?", ids)
		if titleTable == "courses" || titleTable == "paths" {
			tq = tq.Where("deleted_at IS NULL")
		}
		if err := tq.Find(&trows).Error; err != nil {
			return nil, err
		}
		for _, t := range trows {
			titles[t.ID] = t.Title
		}
	}
	items := make([]OpsNamedValue, 0, len(rows))
	for _, item := range rows {
		name := titles[item.GID]
		if name == "" {
			name = "-"
		}
		items = append(items, OpsNamedValue{Name: name, Value: item.C})
	}
	return items, nil
}

func (r *adminOpsRepo) TopText(ctx context.Context, table, column string, cityID int, limited bool, extra string) ([]OpsNamedValue, error) {
	q := r.db.WithContext(ctx).Table(table)
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	type row struct {
		Name string `gorm:"column:name"`
		C    int    `gorm:"column:c"`
	}
	var rows []row
	if err := q.Select("`" + column + "` as name, COUNT(*) as c").Group(column).Order("c DESC").Limit(8).Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]OpsNamedValue, 0, len(rows))
	for _, item := range rows {
		name := item.Name
		if name == "" {
			name = "未命名"
		}
		items = append(items, OpsNamedValue{Name: name, Value: item.C})
	}
	return items, nil
}

func (r *adminOpsRepo) TopAmount(ctx context.Context, table string, cityID int, limited bool, extra string) ([]OpsNamedValue, error) {
	q := r.db.WithContext(ctx).Table(table)
	if extra != "" {
		q = q.Where(extra)
	}
	if limited {
		q = q.Where("user_id IN (SELECT id FROM users WHERE city_id = ?)", cityID)
	}
	type row struct {
		Name string `gorm:"column:name"`
		C    int    `gorm:"column:c"`
	}
	var rows []row
	if err := q.Select("field as name, COALESCE(SUM(`sum`), 0) as c").Group("field").Order("c DESC").Limit(8).Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]OpsNamedValue, 0, len(rows))
	for _, item := range rows {
		if item.C == 0 {
			continue
		}
		name := item.Name
		if name == "" {
			name = "未分类"
		}
		items = append(items, OpsNamedValue{Name: name, Value: item.C})
	}
	return items, nil
}
