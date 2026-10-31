package adminops

import (
	"context"
	"errors"
	"time"

	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/scope"
)

var ErrBadScene = errors.New("参数错误")

func (s *Service) Chart(ctx context.Context, adminID int, scene string) (map[string]any, error) {
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
	if err != nil {
		return nil, err
	}
	switch scene {
	case "course":
		return s.chartCourse(ctx, cityID, limited)
	case "path":
		return s.chartPath(ctx, cityID, limited)
	case "learn":
		return s.chartLearn(ctx, cityID, limited)
	case "cert":
		return s.chartCert(ctx, cityID, limited)
	case "credit":
		return s.chartAmount(ctx, "user_credit_records", "学分", cityID, limited)
	case "coins":
		return s.chartAmount(ctx, "user_coins_records", "金币", cityID, limited)
	default:
		return nil, ErrBadScene
	}
}

func (s *Service) chartCourse(ctx context.Context, cityID int, limited bool) (map[string]any, error) {
	dates := dayDates(14)
	trend, err := s.ops.CountByDate(ctx, "user_course", "created_at", dates, cityID, limited, "")
	if err != nil {
		return nil, err
	}
	total, err := s.ops.ScopedCount(ctx, "user_course", cityID, limited, "")
	if err != nil {
		return nil, err
	}
	ended, err := s.ops.ScopedCount(ctx, "user_course", cityID, limited, "is_end = 1")
	if err != nil {
		return nil, err
	}
	today, err := s.ops.CountSince(ctx, "user_course", "created_at", todayStart(), cityID, limited, "")
	if err != nil {
		return nil, err
	}
	rank, err := s.ops.TopCount(ctx, "user_course", "course_id", "courses", cityID, limited, "")
	if err != nil {
		return nil, err
	}
	return pack("新增入课", dates, trend, []map[string]any{
		{"name": "学习中", "value": max(0, total-ended)},
		{"name": "已结业", "value": ended},
	}, "课程", namedValues(rank), []map[string]any{
		{"label": "累计入课", "value": total},
		{"label": "今日新增", "value": today},
		{"label": "近14天", "value": sumInts(trend)},
		{"label": "已结业", "value": ended},
	}), nil
}

func (s *Service) chartPath(ctx context.Context, cityID int, limited bool) (map[string]any, error) {
	dates := dayDates(14)
	trend, err := s.ops.CountByDate(ctx, "user_path", "created_at", dates, cityID, limited, "")
	if err != nil {
		return nil, err
	}
	total, err := s.ops.ScopedCount(ctx, "user_path", cityID, limited, "")
	if err != nil {
		return nil, err
	}
	ended, err := s.ops.ScopedCount(ctx, "user_path", cityID, limited, "is_end = 1")
	if err != nil {
		return nil, err
	}
	today, err := s.ops.CountSince(ctx, "user_path", "created_at", todayStart(), cityID, limited, "")
	if err != nil {
		return nil, err
	}
	rank, err := s.ops.TopCount(ctx, "user_path", "path_id", "paths", cityID, limited, "")
	if err != nil {
		return nil, err
	}
	return pack("新增入系", dates, trend, []map[string]any{
		{"name": "学习中", "value": max(0, total-ended)},
		{"name": "已结业", "value": ended},
	}, "系列", namedValues(rank), []map[string]any{
		{"label": "累计入系", "value": total},
		{"label": "今日新增", "value": today},
		{"label": "近14天", "value": sumInts(trend)},
		{"label": "已结业", "value": ended},
	}), nil
}

func (s *Service) chartLearn(ctx context.Context, cityID int, limited bool) (map[string]any, error) {
	dates := dayDates(14)
	trend, err := s.ops.CountByDate(ctx, "user_course_chapter", "updated_at", dates, cityID, limited, "is_end = 1")
	if err != nil {
		return nil, err
	}
	total, err := s.ops.ScopedCount(ctx, "user_course_chapter", cityID, limited, "")
	if err != nil {
		return nil, err
	}
	ended, err := s.ops.ScopedCount(ctx, "user_course_chapter", cityID, limited, "is_end = 1")
	if err != nil {
		return nil, err
	}
	today, err := s.ops.CountSince(ctx, "user_course_chapter", "updated_at", todayStart(), cityID, limited, "is_end = 1")
	if err != nil {
		return nil, err
	}
	rank, err := s.ops.TopCount(ctx, "user_course_chapter", "course_id", "courses", cityID, limited, "is_end = 1")
	if err != nil {
		return nil, err
	}
	return pack("学完课节", dates, trend, []map[string]any{
		{"name": "未学完", "value": max(0, total-ended)},
		{"name": "已学完", "value": ended},
	}, "课程", namedValues(rank), []map[string]any{
		{"label": "课节记录", "value": total},
		{"label": "今日学完", "value": today},
		{"label": "近14天学完", "value": sumInts(trend)},
		{"label": "已学完", "value": ended},
	}), nil
}

func (s *Service) chartCert(ctx context.Context, cityID int, limited bool) (map[string]any, error) {
	dates := dayDates(14)
	expr := "DATE(COALESCE(issued_at, created_at))"
	trend, err := s.ops.CountByExpr(ctx, "user_cert", expr, dates, cityID, limited, "")
	if err != nil {
		return nil, err
	}
	total, err := s.ops.ScopedCount(ctx, "user_cert", cityID, limited, "")
	if err != nil {
		return nil, err
	}
	valid, err := s.ops.ScopedCount(ctx, "user_cert", cityID, limited, "is_valid = 1")
	if err != nil {
		return nil, err
	}
	course, err := s.ops.ScopedCount(ctx, "user_cert", cityID, limited, "course_id > 0")
	if err != nil {
		return nil, err
	}
	today, err := s.ops.CountExprSince(ctx, "user_cert", expr, time.Now().In(loc()).Format("2006-01-02"), cityID, limited, "")
	if err != nil {
		return nil, err
	}
	rank, err := s.ops.TopText(ctx, "user_cert", "title", cityID, limited, "")
	if err != nil {
		return nil, err
	}
	return pack("颁发证书", dates, trend, []map[string]any{
		{"name": "有效", "value": valid},
		{"name": "已撤回", "value": max(0, total-valid)},
	}, "证书", namedValues(rank), []map[string]any{
		{"label": "累计颁发", "value": total},
		{"label": "今日颁发", "value": today},
		{"label": "近14天", "value": sumInts(trend)},
		{"label": "来自课程", "value": course},
	}), nil
}

func (s *Service) chartAmount(ctx context.Context, table, unit string, cityID int, limited bool) (map[string]any, error) {
	dates := dayDates(14)
	extra := "deleted_at IS NULL"
	trend, err := s.ops.SumByDate(ctx, table, "created_at", "sum", dates, cityID, limited, extra)
	if err != nil {
		return nil, err
	}
	total, err := s.ops.ScopedSum(ctx, table, "sum", cityID, limited, extra)
	if err != nil {
		return nil, err
	}
	count, err := s.ops.ScopedCount(ctx, table, cityID, limited, extra)
	if err != nil {
		return nil, err
	}
	todaySum, err := s.ops.SumSince(ctx, table, "created_at", "sum", todayStart(), cityID, limited, extra)
	if err != nil {
		return nil, err
	}
	slices, err := s.ops.TopText(ctx, table, "field", cityID, limited, extra)
	if err != nil {
		return nil, err
	}
	fields, err := s.ops.TopAmount(ctx, table, cityID, limited, extra)
	if err != nil {
		return nil, err
	}
	return pack("发放"+unit, dates, trend, namedValues(slices), "发放", namedValues(fields), []map[string]any{
		{"label": "累计" + unit, "value": total},
		{"label": "今日发放", "value": todaySum},
		{"label": "近14天", "value": sumInts(trend)},
		{"label": "流水笔数", "value": count},
	}), nil
}

func namedValues(items []adminrepo.OpsNamedValue) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		out = append(out, map[string]any{"name": item.Name, "value": item.Value})
	}
	return out
}

func pack(trendName string, dates []string, values []int, slices []map[string]any, rankName string, rankItems []map[string]any, stats []map[string]any) map[string]any {
	shortDates := make([]string, len(dates))
	for i, d := range dates {
		if len(d) >= 10 {
			shortDates[i] = d[5:]
		} else {
			shortDates[i] = d
		}
	}
	filtered := make([]map[string]any, 0, len(slices))
	for _, item := range slices {
		if v, ok := item["value"].(int); ok && v > 0 {
			filtered = append(filtered, item)
		}
	}
	return map[string]any{
		"stats": stats,
		"trend": map[string]any{
			"name":   trendName,
			"dates":  shortDates,
			"values": values,
		},
		"slices": filtered,
		"rank": map[string]any{
			"name":  rankName,
			"items": rankItems,
		},
	}
}

func dayDates(days int) []string {
	loc := loc()
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	out := make([]string, 0, days)
	for i := days - 1; i >= 0; i-- {
		out = append(out, today.AddDate(0, 0, -i).Format("2006-01-02"))
	}
	return out
}

func todayStart() string {
	loc := loc()
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Format("2006-01-02 15:04:05")
}

func loc() *time.Location {
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("CST", 8*3600)
	}
	return l
}

func sumInts(vals []int) int {
	n := 0
	for _, v := range vals {
		n += v
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
