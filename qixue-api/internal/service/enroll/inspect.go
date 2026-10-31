package enroll

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/repo"
	"qixue-api/internal/service/cert"
	"qixue-api/internal/service/coin"
)

const (
	KindMissingChapters     = "course_missing_chapters"
	KindOrphanChapters      = "course_orphan_chapters"
	KindCourseIDMismatch    = "course_id_mismatch"
	KindMissingEnroll       = "course_missing_enroll"
	KindProgressMismatch    = "course_progress_mismatch"
	KindEndedIncomplete     = "course_ended_incomplete"
	KindChapterFlagMismatch = "chapter_flag_mismatch"
	KindPathMissingCourses  = "path_missing_courses"
	KindPathProgress        = "path_progress_mismatch"
	KindPathEndedIncomplete = "path_ended_incomplete"
	KindMissingCredit       = "course_missing_credit"
	KindMissingCoins        = "course_missing_coins"
	KindMissingCourseCert   = "course_missing_cert"
	KindPathShouldEnd       = "path_should_end"
	KindMissingPathCert     = "path_missing_cert"

	inspectKindLimit = 300
	repairBatchLimit = 50
	inspectErrLimit  = 20
)

var ErrUnknownKind = errors.New("未知巡检类型")

var kindLabels = map[string]string{
	KindMissingChapters:     "课节与课程不一致",
	KindOrphanChapters:      "课节记录孤立",
	KindCourseIDMismatch:    "课节课程不一致",
	KindMissingEnroll:       "缺少入课记录",
	KindProgressMismatch:    "课程完成数异常",
	KindEndedIncomplete:     "已结业但课节未齐",
	KindChapterFlagMismatch: "课节学完标记异常",
	KindPathMissingCourses:  "系列课程缺失",
	KindPathProgress:        "系列完成数异常",
	KindPathEndedIncomplete: "系列已结业但课程未齐",
	KindMissingCredit:       "结业未发学分",
	KindMissingCoins:        "结业未发积分",
	KindMissingCourseCert:   "结业未发证书",
	KindPathShouldEnd:       "系列课已学完未结业",
	KindMissingPathCert:     "系列未发证书",
}

// Issue is one abnormal student/course/chapter/path row.
type Issue struct {
	Key       string `json:"key"`
	Kind      string `json:"kind"`
	Label     string `json:"label"`
	Detail    string `json:"detail"`
	UserID    int    `json:"user_id"`
	CourseID  int    `json:"course_id,omitempty"`
	PathID    int    `json:"path_id,omitempty"`
	ChapterID int    `json:"chapter_id,omitempty"`
	RecordID  int    `json:"record_id,omitempty"`
	Expect    int    `json:"expect,omitempty"`
	Actual    int    `json:"actual,omitempty"`
}

// KindCount is a per-kind tally for the scan summary.
type KindCount struct {
	Kind  string `json:"kind"`
	Label string `json:"label"`
	Count int    `json:"count,omitempty"`
}

// ScanResult is the inspect payload returned to admin.
type ScanResult struct {
	Kind      string      `json:"kind"`
	Label     string      `json:"label"`
	Total     int         `json:"total"`
	Truncated bool        `json:"truncated"`
	Summary   []KindCount `json:"summary"`
	Issues    []Issue     `json:"issues"`
}

// RepairResult reports how many issues of one kind were fixed.
type RepairResult struct {
	Kind      string   `json:"kind"`
	Label     string   `json:"label"`
	Scanned   int      `json:"scanned"`
	Repaired  int      `json:"repaired"`
	Failed    int      `json:"failed"`
	Remaining int      `json:"remaining"`
	Truncated bool     `json:"truncated"`
	Errors    []string `json:"errors"`
}

type inspectRow struct {
	UserID    int `gorm:"column:user_id"`
	CourseID  int `gorm:"column:course_id"`
	PathID    int `gorm:"column:path_id"`
	ChapterID int `gorm:"column:chapter_id"`
	RecordID  int `gorm:"column:record_id"`
	Expect    int `gorm:"column:expect_n"`
	Actual    int `gorm:"column:actual_n"`
}

type kindSpec struct {
	kind   string
	sql    string
	detail func(inspectRow) string
}

type scanner struct {
	ctx   context.Context
	db    *gorm.DB
	trunc bool
	items []Issue
	sum   []KindCount
}

var kindSpecs = []kindSpec{
	{
		kind: KindMissingChapters,
		sql: `
SELECT uc.user_id, uc.course_id, uc.id AS record_id,
  COUNT(cc.id) AS expect_n,
  COUNT(ucc.id) AS actual_n
FROM user_course uc
INNER JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL
INNER JOIN course_chapters cc ON cc.course_id = uc.course_id AND cc.deleted_at IS NULL
LEFT JOIN user_course_chapter ucc ON ucc.user_id = uc.user_id AND ucc.chapter_id = cc.id
GROUP BY uc.id, uc.user_id, uc.course_id
HAVING COUNT(cc.id) > COUNT(ucc.id)`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("课程有 %d 课节，学员有 %d 条，修复时补齐", row.Expect, row.Actual)
		},
	},
	{
		kind: KindOrphanChapters,
		sql: `
SELECT ucc.id AS record_id, ucc.user_id, ucc.course_id, ucc.chapter_id
FROM user_course_chapter ucc
LEFT JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL
WHERE cc.id IS NULL`,
		detail: func(inspectRow) string { return "课节已删除，学习记录悬空" },
	},
	{
		kind: KindCourseIDMismatch,
		sql: `
SELECT ucc.id AS record_id, ucc.user_id, ucc.course_id, ucc.chapter_id, cc.course_id AS expect_n, IFNULL(ucc.course_id, 0) AS actual_n
FROM user_course_chapter ucc
INNER JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL
WHERE IFNULL(ucc.course_id, 0) <> cc.course_id`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("课节属于课程 %d，记录写成了 %d", row.Expect, row.Actual)
		},
	},
	{
		kind: KindMissingEnroll,
		sql: `
SELECT ucc.user_id, ucc.course_id, MIN(ucc.id) AS record_id, COUNT(*) AS actual_n
FROM user_course_chapter ucc
INNER JOIN courses c ON c.id = ucc.course_id AND c.deleted_at IS NULL
LEFT JOIN user_course uc ON uc.user_id = ucc.user_id AND uc.course_id = ucc.course_id
WHERE uc.id IS NULL AND IFNULL(ucc.course_id, 0) > 0
GROUP BY ucc.user_id, ucc.course_id`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("有 %d 条课节记录，但没有入课记录", row.Actual)
		},
	},
	{
		kind: KindProgressMismatch,
		sql: `
SELECT uc.user_id, uc.course_id, uc.id AS record_id, IFNULL(e.ended, 0) AS expect_n, uc.has_end_chapter AS actual_n
FROM user_course uc
INNER JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL
LEFT JOIN (
  SELECT ucc.user_id, ucc.course_id, SUM(ucc.is_end = 1) AS ended
  FROM user_course_chapter ucc
  INNER JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL
  GROUP BY ucc.user_id, ucc.course_id
) e ON e.user_id = uc.user_id AND e.course_id = uc.course_id
WHERE uc.has_end_chapter <> IFNULL(e.ended, 0)`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("已学完课节 %d，入课统计为 %d", row.Expect, row.Actual)
		},
	},
	{
		kind: KindEndedIncomplete,
		sql: `
SELECT uc.user_id, uc.course_id, uc.id AS record_id, t.expect_n, IFNULL(e.ended, 0) AS actual_n
FROM user_course uc
INNER JOIN (
  SELECT course_id, COUNT(*) AS expect_n FROM course_chapters WHERE deleted_at IS NULL GROUP BY course_id
) t ON t.course_id = uc.course_id
LEFT JOIN (
  SELECT ucc.user_id, ucc.course_id, SUM(ucc.is_end = 1) AS ended
  FROM user_course_chapter ucc
  INNER JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL
  GROUP BY ucc.user_id, ucc.course_id
) e ON e.user_id = uc.user_id AND e.course_id = uc.course_id
INNER JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL
WHERE uc.is_end = 1 AND IFNULL(e.ended, 0) < t.expect_n`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("课程已结业，但仅 %d/%d 课节标记学完", row.Actual, row.Expect)
		},
	},
	{
		kind: KindMissingCredit,
		sql: `
SELECT uc.user_id, uc.course_id, uc.id AS record_id, c.credit AS expect_n
FROM user_course uc
INNER JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL
WHERE uc.is_end = 1 AND c.credit > 0
  AND NOT EXISTS (
    SELECT 1 FROM user_credit_records r
    WHERE r.user_id = uc.user_id
      AND r.field = '课程发放'
      AND r.remark = CONCAT(c.title, '课程学习结束发放')
  )`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("课程已结业，应得学分 %d，尚未发放", row.Expect)
		},
	},
	{
		kind: KindMissingCoins,
		sql: `
SELECT uc.user_id, uc.course_id, uc.id AS record_id
FROM user_course uc
INNER JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL
WHERE uc.is_end = 1
  AND NOT EXISTS (
    SELECT 1 FROM user_coins_records r
    WHERE r.user_id = uc.user_id
      AND r.field = '课程结业'
      AND r.deleted_at IS NULL
      AND r.remark LIKE CONCAT('%#course', uc.course_id, '#%')
  )`,
		detail: func(inspectRow) string { return "课程已结业，尚未发放结业积分" },
	},
	{
		kind: KindMissingCourseCert,
		sql: `
SELECT uc.user_id, uc.course_id, uc.id AS record_id
FROM user_course uc
INNER JOIN courses c ON c.id = uc.course_id AND c.deleted_at IS NULL
WHERE uc.is_end = 1 AND IFNULL(c.cert_id, 0) > 0
  AND NOT EXISTS (
    SELECT 1 FROM user_cert ucert
    WHERE ucert.user_id = uc.user_id AND ucert.course_id = uc.course_id AND IFNULL(ucert.path_id, 0) = 0
  )`,
		detail: func(inspectRow) string { return "课程已结业且配置了证书，尚未发放" },
	},
	{
		kind: KindChapterFlagMismatch,
		sql: `
SELECT ucc.id AS record_id, ucc.user_id, ucc.course_id, ucc.chapter_id
FROM user_course_chapter ucc
WHERE ucc.is_end = 1 AND IFNULL(ucc.has_end_chapter, 0) = 0`,
		detail: func(inspectRow) string { return "课节已结业，学完标记仍为 0" },
	},
	{
		kind: KindPathMissingCourses,
		sql: `
SELECT up.user_id, up.path_id, pc.course_id, up.id AS record_id
FROM user_path up
INNER JOIN paths p ON p.id = up.path_id AND p.deleted_at IS NULL
INNER JOIN path_courses pc ON pc.path_id = up.path_id AND pc.deleted_at IS NULL
INNER JOIN courses c ON c.id = pc.course_id AND c.deleted_at IS NULL
LEFT JOIN user_course uc ON uc.user_id = up.user_id AND uc.course_id = pc.course_id
WHERE up.deleted_at IS NULL AND uc.id IS NULL`,
		detail: func(inspectRow) string { return "已加入系列，但系列内课程未入课" },
	},
	{
		kind: KindPathProgress,
		sql: `
SELECT up.user_id, up.path_id, up.id AS record_id, IFNULL(e.ended, 0) AS expect_n, up.has_end_course AS actual_n
FROM user_path up
INNER JOIN paths p ON p.id = up.path_id AND p.deleted_at IS NULL
LEFT JOIN (
  SELECT up2.user_id, up2.path_id, COUNT(uc.id) AS ended
  FROM user_path up2
  INNER JOIN path_courses pc ON pc.path_id = up2.path_id AND pc.deleted_at IS NULL
  INNER JOIN user_course uc ON uc.user_id = up2.user_id AND uc.course_id = pc.course_id AND uc.is_end = 1
  WHERE up2.deleted_at IS NULL
  GROUP BY up2.user_id, up2.path_id
) e ON e.user_id = up.user_id AND e.path_id = up.path_id
WHERE up.deleted_at IS NULL AND up.has_end_course <> IFNULL(e.ended, 0)`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("已完成课程 %d，系列统计为 %d", row.Expect, row.Actual)
		},
	},
	{
		kind: KindPathEndedIncomplete,
		sql: `
SELECT up.user_id, up.path_id, up.id AS record_id, t.expect_n, IFNULL(e.ended, 0) AS actual_n
FROM user_path up
INNER JOIN paths p ON p.id = up.path_id AND p.deleted_at IS NULL
INNER JOIN (
  SELECT path_id, COUNT(*) AS expect_n FROM path_courses WHERE deleted_at IS NULL GROUP BY path_id
) t ON t.path_id = up.path_id
LEFT JOIN (
  SELECT up2.user_id, up2.path_id, COUNT(uc.id) AS ended
  FROM user_path up2
  INNER JOIN path_courses pc ON pc.path_id = up2.path_id AND pc.deleted_at IS NULL
  INNER JOIN user_course uc ON uc.user_id = up2.user_id AND uc.course_id = pc.course_id AND uc.is_end = 1
  WHERE up2.deleted_at IS NULL
  GROUP BY up2.user_id, up2.path_id
) e ON e.user_id = up.user_id AND e.path_id = up.path_id
WHERE up.deleted_at IS NULL AND up.is_end = 1 AND IFNULL(e.ended, 0) < t.expect_n`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("系列已结业，但仅 %d/%d 课程完成", row.Actual, row.Expect)
		},
	},
	{
		kind: KindPathShouldEnd,
		sql: `
SELECT up.user_id, up.path_id, up.id AS record_id, t.expect_n, IFNULL(e.ended, 0) AS actual_n
FROM user_path up
INNER JOIN paths p ON p.id = up.path_id AND p.deleted_at IS NULL
INNER JOIN (
  SELECT path_id, COUNT(*) AS expect_n FROM path_courses WHERE deleted_at IS NULL GROUP BY path_id
) t ON t.path_id = up.path_id
LEFT JOIN (
  SELECT up2.user_id, up2.path_id, COUNT(uc.id) AS ended
  FROM user_path up2
  INNER JOIN path_courses pc ON pc.path_id = up2.path_id AND pc.deleted_at IS NULL
  INNER JOIN user_course uc ON uc.user_id = up2.user_id AND uc.course_id = pc.course_id AND uc.is_end = 1
  WHERE up2.deleted_at IS NULL
  GROUP BY up2.user_id, up2.path_id
) e ON e.user_id = up.user_id AND e.path_id = up.path_id
WHERE up.deleted_at IS NULL AND up.is_end = 0 AND t.expect_n > 0 AND IFNULL(e.ended, 0) >= t.expect_n`,
		detail: func(row inspectRow) string {
			return fmt.Sprintf("系列内 %d 门课程均已结业，系列未标记完成", row.Expect)
		},
	},
	{
		kind: KindMissingPathCert,
		sql: `
SELECT up.user_id, up.path_id, up.id AS record_id
FROM user_path up
INNER JOIN paths p ON p.id = up.path_id AND p.deleted_at IS NULL
WHERE up.deleted_at IS NULL AND up.is_end = 1 AND IFNULL(p.cert_id, 0) > 0
  AND NOT EXISTS (
    SELECT 1 FROM user_cert ucert
    WHERE ucert.user_id = up.user_id AND ucert.path_id = up.path_id AND IFNULL(ucert.course_id, 0) = 0
  )`,
		detail: func(inspectRow) string { return "系列已结业且配置了证书，尚未发放" },
	},
}

func KindLabel(kind string) string {
	if s := kindLabels[kind]; s != "" {
		return s
	}
	return kind
}

func lookupKind(kind string) (kindSpec, bool) {
	kind = strings.TrimSpace(kind)
	for _, spec := range kindSpecs {
		if spec.kind == kind {
			return spec, true
		}
	}
	return kindSpec{}, false
}

// KindSteps returns inspect/repair kinds in execution order.
func KindSteps() []KindCount {
	out := make([]KindCount, 0, len(kindSpecs))
	for _, spec := range kindSpecs {
		out = append(out, KindCount{Kind: spec.kind, Label: KindLabel(spec.kind)})
	}
	return out
}

func queryKindRows(ctx context.Context, db *gorm.DB, sql string, limit int) ([]inspectRow, error) {
	sql = strings.TrimSpace(sql)
	if limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	var rows []inspectRow
	err := db.WithContext(ctx).Raw(sql).Scan(&rows).Error
	return rows, err
}

func (s *scanner) collect(spec kindSpec) error {
	rows, err := queryKindRows(s.ctx, s.db, spec.sql, inspectKindLimit+1)
	if err != nil {
		return err
	}
	if len(rows) > inspectKindLimit {
		s.trunc = true
		rows = rows[:inspectKindLimit]
	}
	if len(rows) == 0 {
		return nil
	}
	label := KindLabel(spec.kind)
	s.sum = append(s.sum, KindCount{Kind: spec.kind, Label: label, Count: len(rows)})
	for _, row := range rows {
		text := ""
		if spec.detail != nil {
			text = spec.detail(row)
		}
		s.items = append(s.items, Issue{
			Key:       fmt.Sprintf("%s|%d|%d|%d|%d|%d", spec.kind, row.UserID, row.CourseID, row.PathID, row.ChapterID, row.RecordID),
			Kind:      spec.kind,
			Label:     label,
			Detail:    text,
			UserID:    row.UserID,
			CourseID:  row.CourseID,
			PathID:    row.PathID,
			ChapterID: row.ChapterID,
			RecordID:  row.RecordID,
			Expect:    row.Expect,
			Actual:    row.Actual,
		})
	}
	return nil
}

// ScanKind finds issues of a single inspect kind.
func ScanKind(ctx context.Context, db *gorm.DB, kind string) (ScanResult, error) {
	spec, ok := lookupKind(kind)
	if !ok {
		return ScanResult{}, ErrUnknownKind
	}
	s := &scanner{ctx: ctx, db: db}
	if err := s.collect(spec); err != nil {
		return ScanResult{}, err
	}
	return ScanResult{
		Kind:      spec.kind,
		Label:     KindLabel(spec.kind),
		Total:     len(s.items),
		Truncated: s.trunc,
		Summary:   s.sum,
		Issues:    s.items,
	}, nil
}

func pairKey(uid, id int) string {
	return fmt.Sprintf("%d:%d", uid, id)
}

func noteErr(errs *[]string, failed *int, err error) {
	*failed++
	if err == nil || len(*errs) >= inspectErrLimit {
		return
	}
	*errs = append(*errs, err.Error())
}

func chunkInts(ids []int, size int) [][]int {
	if size <= 0 {
		size = 400
	}
	var out [][]int
	for start := 0; start < len(ids); start += size {
		end := start + size
		if end > len(ids) {
			end = len(ids)
		}
		out = append(out, ids[start:end])
	}
	return out
}

func markRepaired(out *RepairResult) {
	out.Repaired++
}

func syncCourse(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, courseID int, out *RepairResult) {
	if uid <= 0 || courseID <= 0 {
		return
	}
	n, err := EnsureUserCourseChapters(ctx, db, rdb, uid, courseID)
	if err != nil {
		noteErr(&out.Errors, &out.Failed, err)
		return
	}
	ForgetListCache(ctx, rdb, uid)
	if n > 0 {
		out.Repaired += n
		return
	}
	markRepaired(out)
}

func syncPath(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, pathID int, out *RepairResult) {
	if uid <= 0 || pathID <= 0 {
		return
	}
	if err := SyncUserPathStatus(ctx, db, rdb, uid, pathID); err != nil {
		noteErr(&out.Errors, &out.Failed, err)
		return
	}
	markRepaired(out)
}

func repairKindRows(ctx context.Context, db *gorm.DB, rdb *redis.Client, kind string, rows []inspectRow, out *RepairResult) {
	switch kind {
	case KindMissingChapters, KindProgressMismatch, KindEndedIncomplete:
		seen := map[string]struct{}{}
		for _, row := range rows {
			key := pairKey(row.UserID, row.CourseID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			syncCourse(ctx, db, rdb, row.UserID, row.CourseID, out)
		}
	case KindOrphanChapters:
		var ids []int
		for _, row := range rows {
			if row.RecordID > 0 {
				ids = append(ids, row.RecordID)
			}
		}
		for _, chunk := range chunkInts(ids, 400) {
			if err := db.WithContext(ctx).Exec("DELETE FROM user_course_chapter WHERE id IN ?", chunk).Error; err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			out.Repaired += len(chunk)
		}
	case KindCourseIDMismatch:
		var ids []int
		seen := map[string]struct{}{}
		for _, row := range rows {
			if row.RecordID > 0 {
				ids = append(ids, row.RecordID)
			}
		}
		for _, chunk := range chunkInts(ids, 400) {
			if err := db.WithContext(ctx).Exec(`
UPDATE user_course_chapter ucc
INNER JOIN course_chapters cc ON cc.id = ucc.chapter_id AND cc.deleted_at IS NULL
SET ucc.course_id = cc.course_id
WHERE ucc.id IN ?`, chunk).Error; err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			out.Repaired += len(chunk)
		}
		for _, row := range rows {
			for _, courseID := range []int{row.Expect, row.CourseID} {
				key := pairKey(row.UserID, courseID)
				if courseID <= 0 {
					continue
				}
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				syncCourse(ctx, db, rdb, row.UserID, courseID, out)
			}
		}
	case KindMissingEnroll:
		cr := courses(db)
		now := time.Now()
		seen := map[string]struct{}{}
		for _, row := range rows {
			key := pairKey(row.UserID, row.CourseID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			if err := cr.CreateUserCoursesIgnore(ctx, []map[string]any{repo.NewUserCourseRow(row.UserID, row.CourseID, now)}); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			syncCourse(ctx, db, rdb, row.UserID, row.CourseID, out)
		}
	case KindChapterFlagMismatch:
		var ids []int
		seen := map[string]struct{}{}
		for _, row := range rows {
			if row.RecordID > 0 {
				ids = append(ids, row.RecordID)
			}
		}
		for _, chunk := range chunkInts(ids, 400) {
			if err := db.WithContext(ctx).Table("user_course_chapter").Where("id IN ?", chunk).
				Update("has_end_chapter", 1).Error; err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			out.Repaired += len(chunk)
		}
		for _, row := range rows {
			key := pairKey(row.UserID, row.CourseID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			if err := SyncUserCourseStatus(ctx, db, rdb, row.UserID, row.CourseID); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
			}
		}
	case KindMissingCredit:
		type courseCredit struct {
			Title  string `gorm:"column:title"`
			Credit int    `gorm:"column:credit"`
		}
		for _, row := range rows {
			var course courseCredit
			if err := db.WithContext(ctx).Table("courses").
				Select("title, credit").Where("id = ? AND deleted_at IS NULL", row.CourseID).
				Take(&course).Error; err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			if course.Credit <= 0 {
				continue
			}
			remark := course.Title + "课程学习结束发放"
			granted, err := CourseCreditAlreadyGranted(ctx, db, row.UserID, remark)
			if err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			if granted {
				continue
			}
			if err := coin.AddCredit(ctx, db, row.UserID, "课程发放", course.Credit, remark); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			markRepaired(out)
		}
	case KindMissingCoins:
		for _, row := range rows {
			var title string
			if err := db.WithContext(ctx).Table("courses").
				Select("title").Where("id = ?", row.CourseID).Scan(&title).Error; err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			coin.OnCourseEnd(ctx, db, rdb, row.UserID, row.CourseID, title)
			markRepaired(out)
		}
	case KindMissingCourseCert:
		issuer := cert.NewIssuer(repo.NewCert(db))
		for _, row := range rows {
			if _, err := issuer.IssueForCourse(ctx, row.UserID, row.CourseID, false); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			markRepaired(out)
		}
	case KindPathMissingCourses:
		seen := map[string]struct{}{}
		pathSeen := map[string]struct{}{}
		for _, row := range rows {
			key := pairKey(row.UserID, row.CourseID)
			if _, ok := seen[key]; !ok && row.CourseID > 0 {
				seen[key] = struct{}{}
				if err := EnrollUsersToCourse(ctx, db, rdb, row.CourseID, []int{row.UserID}); err != nil {
					noteErr(&out.Errors, &out.Failed, err)
					continue
				}
				markRepaired(out)
			}
			pkey := pairKey(row.UserID, row.PathID)
			if _, ok := pathSeen[pkey]; ok || row.PathID <= 0 {
				continue
			}
			pathSeen[pkey] = struct{}{}
			if err := SyncUserPathStatus(ctx, db, rdb, row.UserID, row.PathID); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
			}
		}
	case KindPathProgress, KindPathEndedIncomplete:
		seen := map[string]struct{}{}
		for _, row := range rows {
			key := pairKey(row.UserID, row.PathID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			syncPath(ctx, db, rdb, row.UserID, row.PathID, out)
		}
	case KindPathShouldEnd:
		issuer := cert.NewIssuer(repo.NewCert(db))
		seen := map[string]struct{}{}
		for _, row := range rows {
			key := pairKey(row.UserID, row.PathID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			if err := SyncUserPathStatus(ctx, db, rdb, row.UserID, row.PathID); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			if _, err := issuer.IssueForPath(ctx, row.UserID, row.PathID, false); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			markRepaired(out)
		}
	case KindMissingPathCert:
		issuer := cert.NewIssuer(repo.NewCert(db))
		for _, row := range rows {
			if _, err := issuer.IssueForPath(ctx, row.UserID, row.PathID, false); err != nil {
				noteErr(&out.Errors, &out.Failed, err)
				continue
			}
			markRepaired(out)
		}
	}
}

// RepairKind fixes one inspect kind in a bounded batch.
func RepairKind(ctx context.Context, db *gorm.DB, rdb *redis.Client, kind string) (RepairResult, error) {
	spec, ok := lookupKind(kind)
	if !ok {
		return RepairResult{}, ErrUnknownKind
	}
	out := RepairResult{Kind: spec.kind, Label: KindLabel(spec.kind)}
	rows, err := queryKindRows(ctx, db, spec.sql, repairBatchLimit)
	if err != nil {
		return out, err
	}
	out.Scanned = len(rows)
	if len(rows) > 0 {
		repairKindRows(ctx, db, rdb, spec.kind, rows, &out)
	}
	left, err := ScanKind(ctx, db, spec.kind)
	if err != nil {
		return out, err
	}
	out.Remaining = left.Total
	out.Truncated = left.Truncated
	return out, nil
}
