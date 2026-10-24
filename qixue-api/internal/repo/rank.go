package repo

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"
)

// RankRepository covers rank snapshot boards and rebuild queries.
type RankRepository interface {
	EnsureSnapshotsTable(ctx context.Context) error
	FindSnapshot(ctx context.Context, siteCode, typ, period, periodKey, snapshotDate string) (*RankSnapshotRef, error)
	UpdateSnapshot(ctx context.Context, id int, periodLabel, payload string, updatedAt time.Time) error
	CreateSnapshot(ctx context.Context, row RankSnapshotCreate) error
	ListSnapshotsByPeriod(ctx context.Context, siteCode, typ, period string) ([]RankSnapRow, error)
	ListSnapshotsByPeriodKeys(ctx context.Context, siteCode string, periodKeys []string) ([]RankSnapRow, error)
	AggregateDuration(ctx context.Context, start, end string, limit int) ([]RankIDValue, error)
	AggregateCourse(ctx context.Context, start, end string, limit int) ([]RankIDValue, error)
	AggregateCoins(ctx context.Context, start, end string, limit int) ([]RankIDValue, error)
	ListCoursesBrief(ctx context.Context, ids []int) ([]RankCourseBrief, error)
	ListUsersBrief(ctx context.Context, ids []int) ([]RankUserBrief, error)
	ListUserCoins(ctx context.Context, ids []int) ([]RankUserCoins, error)
}

type RankSnapshotRef struct {
	ID      int    `gorm:"column:id"`
	Payload string `gorm:"column:payload"`
}

type RankSnapshotCreate struct {
	SiteCode     string
	Type         string
	Period       string
	PeriodKey    string
	PeriodLabel  string
	SnapshotDate string
	Payload      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RankSnapRow struct {
	Type         string     `gorm:"column:type"`
	PeriodKey    string     `gorm:"column:period_key"`
	PeriodLabel  string     `gorm:"column:period_label"`
	Payload      string     `gorm:"column:payload"`
	SnapshotDate time.Time  `gorm:"column:snapshot_date"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

type RankIDValue struct {
	ID    int `gorm:"column:id"`
	Value int `gorm:"column:value"`
}

type RankCourseBrief struct {
	ID     int    `gorm:"column:id"`
	Title  string `gorm:"column:title"`
	Thumb  string `gorm:"column:thumb"`
	IsShow int    `gorm:"column:is_show"`
}

type RankUserBrief struct {
	ID       int     `gorm:"column:id"`
	NickName *string `gorm:"column:nick_name"`
	Avatar   string  `gorm:"column:avatar"`
}

type RankUserCoins struct {
	ID    int `gorm:"column:id"`
	Coins int `gorm:"column:coins"`
}

const rankCreateTableSQL = `
CREATE TABLE IF NOT EXISTS ` + "`rank_snapshots`" + ` (
  ` + "`id`" + ` int unsigned NOT NULL AUTO_INCREMENT,
  ` + "`site_code`" + ` varchar(16) NOT NULL DEFAULT '',
  ` + "`type`" + ` varchar(16) NOT NULL,
  ` + "`period`" + ` varchar(8) NOT NULL,
  ` + "`period_key`" + ` varchar(16) NOT NULL,
  ` + "`period_label`" + ` varchar(64) NOT NULL DEFAULT '',
  ` + "`snapshot_date`" + ` date NOT NULL,
  ` + "`payload`" + ` mediumtext NOT NULL,
  ` + "`created_at`" + ` timestamp NULL DEFAULT NULL,
  ` + "`updated_at`" + ` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (` + "`id`" + `),
  UNIQUE KEY ` + "`uk_rank_snap`" + ` (` + "`site_code`" + `,` + "`type`" + `,` + "`period`" + `,` + "`period_key`" + `,` + "`snapshot_date`" + `),
  KEY ` + "`idx_rank_type`" + ` (` + "`site_code`" + `,` + "`type`" + `,` + "`period`" + `,` + "`snapshot_date`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

var (
	rankEnsureOnce sync.Once
	rankEnsureErr  error
)

type rankRepo struct{ db *gorm.DB }

func NewRank(db *gorm.DB) RankRepository {
	return &rankRepo{db: db}
}

func (r *rankRepo) EnsureSnapshotsTable(ctx context.Context) error {
	rankEnsureOnce.Do(func() {
		rankEnsureErr = r.db.WithContext(context.Background()).Exec(rankCreateTableSQL).Error
	})
	return rankEnsureErr
}

func (r *rankRepo) FindSnapshot(ctx context.Context, siteCode, typ, period, periodKey, snapshotDate string) (*RankSnapshotRef, error) {
	var row RankSnapshotRef
	err := r.db.WithContext(ctx).Table("rank_snapshots").
		Select("id, payload").
		Where("site_code = ? AND type = ? AND period = ? AND period_key = ? AND snapshot_date = ?",
			siteCode, typ, period, periodKey, snapshotDate).
		Limit(1).Find(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	return &row, nil
}

func (r *rankRepo) UpdateSnapshot(ctx context.Context, id int, periodLabel, payload string, updatedAt time.Time) error {
	return r.db.WithContext(ctx).Table("rank_snapshots").Where("id = ?", id).Updates(map[string]any{
		"period_label": periodLabel,
		"payload":      payload,
		"updated_at":   updatedAt,
	}).Error
}

func (r *rankRepo) CreateSnapshot(ctx context.Context, row RankSnapshotCreate) error {
	return r.db.WithContext(ctx).Table("rank_snapshots").Create(map[string]any{
		"site_code":     row.SiteCode,
		"type":          row.Type,
		"period":        row.Period,
		"period_key":    row.PeriodKey,
		"period_label":  row.PeriodLabel,
		"snapshot_date": row.SnapshotDate,
		"payload":       row.Payload,
		"created_at":    row.CreatedAt,
		"updated_at":    row.UpdatedAt,
	}).Error
}

func (r *rankRepo) ListSnapshotsByPeriod(ctx context.Context, siteCode, typ, period string) ([]RankSnapRow, error) {
	var rows []RankSnapRow
	err := r.db.WithContext(ctx).Table("rank_snapshots").
		Where("site_code = ? AND type = ? AND period = ?", siteCode, typ, period).
		Order("snapshot_date DESC, id DESC").Find(&rows).Error
	return rows, err
}

func (r *rankRepo) ListSnapshotsByPeriodKeys(ctx context.Context, siteCode string, periodKeys []string) ([]RankSnapRow, error) {
	var rows []RankSnapRow
	err := r.db.WithContext(ctx).Table("rank_snapshots").
		Where("site_code = ? AND period_key IN ?", siteCode, periodKeys).
		Order("snapshot_date DESC, id DESC").Find(&rows).Error
	return rows, err
}

func (r *rankRepo) AggregateDuration(ctx context.Context, start, end string, limit int) ([]RankIDValue, error) {
	var rows []RankIDValue
	err := r.db.WithContext(ctx).Raw(`
SELECT ucc.user_id AS id, COALESCE(SUM(mf.duration), 0) AS value
FROM user_course_chapter ucc
JOIN course_chapters cc ON cc.id = ucc.chapter_id
JOIN media_file mf ON mf.id = cc.media_id
WHERE ucc.has_end_chapter = 1
  AND ucc.updated_at >= ? AND ucc.updated_at < ?
  AND ucc.user_id NOT IN (
    SELECT id FROM users WHERE is_lock = 1 OR is_active = -1
  )
GROUP BY ucc.user_id
HAVING COALESCE(SUM(mf.duration), 0) > 0
ORDER BY value DESC
LIMIT ?
`, start, end, limit).Scan(&rows).Error
	return rows, err
}

func (r *rankRepo) AggregateCourse(ctx context.Context, start, end string, limit int) ([]RankIDValue, error) {
	var rows []RankIDValue
	err := r.db.WithContext(ctx).Raw(`
SELECT ucc.course_id AS id, COUNT(DISTINCT ucc.user_id) AS value
FROM user_course_chapter ucc
WHERE ucc.has_end_chapter = 1
  AND ucc.updated_at >= ? AND ucc.updated_at < ?
  AND ucc.user_id NOT IN (
    SELECT id FROM users WHERE is_lock = 1 OR is_active = -1
  )
GROUP BY ucc.course_id
HAVING COUNT(DISTINCT ucc.user_id) > 0
ORDER BY value DESC
LIMIT ?
`, start, end, limit).Scan(&rows).Error
	return rows, err
}

func (r *rankRepo) AggregateCoins(ctx context.Context, start, end string, limit int) ([]RankIDValue, error) {
	var rows []RankIDValue
	err := r.db.WithContext(ctx).Raw(`
SELECT uco.user_id AS id, SUM(uco.`+"`sum`"+`) AS value
FROM user_coins_records uco
WHERE uco.`+"`sum`"+` > 0 AND uco.deleted_at IS NULL
  AND uco.created_at >= ? AND uco.created_at < ?
  AND uco.user_id NOT IN (
    SELECT id FROM users WHERE is_lock = 1 OR is_active = -1
  )
GROUP BY uco.user_id
HAVING SUM(uco.`+"`sum`"+`) > 0
ORDER BY value DESC
LIMIT ?
`, start, end, limit).Scan(&rows).Error
	return rows, err
}

func (r *rankRepo) ListCoursesBrief(ctx context.Context, ids []int) ([]RankCourseBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []RankCourseBrief
	err := r.db.WithContext(ctx).Table("courses").
		Select("id, title, thumb, is_show").
		Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (r *rankRepo) ListUsersBrief(ctx context.Context, ids []int) ([]RankUserBrief, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []RankUserBrief
	err := r.db.WithContext(ctx).Table("users").Select("id, nick_name, avatar").Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (r *rankRepo) ListUserCoins(ctx context.Context, ids []int) ([]RankUserCoins, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []RankUserCoins
	err := r.db.WithContext(ctx).Table("users").Select("id, coins").Where("id IN ?", ids).Find(&rows).Error
	return rows, err
}
