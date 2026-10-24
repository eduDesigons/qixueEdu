package rank

import (
	"fmt"
	"time"
)

type periodInfo struct {
	Period       string
	Key          string
	Label        string
	Start        string
	End          string
	LiveEnd      string
	SnapshotDate string
}

func periodMeta(period, which string) periodInfo {
	now := time.Now()
	if period != "month" {
		period = "week"
	}
	loc := now.Location()
	if period == "month" {
		var start time.Time
		var endExclusive time.Time
		if which == "previous" {
			firstThis := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
			start = firstThis.AddDate(0, -1, 0)
			endExclusive = firstThis
		} else {
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
			endExclusive = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		}
		liveEnd := now
		snap := now.Format("2006-01-02")
		if which == "previous" {
			snap = endExclusive.Format("2006-01-02")
		}
		return periodInfo{
			Period: period, Key: start.Format("2006-01"),
			Label: fmt.Sprintf("%d年%d月", start.Year(), int(start.Month())),
			Start: fmtTime(start), End: fmtTime(endExclusive), LiveEnd: fmtTime(liveEnd),
			SnapshotDate: snap,
		}
	}

	weekStart := startOfWeekMonday(now)
	var start time.Time
	var endExclusive time.Time
	if which == "previous" {
		start = weekStart.AddDate(0, 0, -7)
		endExclusive = weekStart
	} else {
		start = weekStart
		endExclusive = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	}
	liveEnd := now
	weekEnd := start.AddDate(0, 0, 6)
	isoYear, isoWeek := start.ISOWeek()
	snap := now.Format("2006-01-02")
	if which == "previous" {
		snap = endExclusive.Format("2006-01-02")
	}
	return periodInfo{
		Period: period,
		Key:    fmt.Sprintf("%d-W%02d", isoYear, isoWeek),
		Label: fmt.Sprintf("%d月%d日-%d月%d日",
			int(start.Month()), start.Day(), int(weekEnd.Month()), weekEnd.Day()),
		Start: fmtTime(start), End: fmtTime(endExclusive), LiveEnd: fmtTime(liveEnd),
		SnapshotDate: snap,
	}
}

func startOfWeekMonday(t time.Time) time.Time {
	loc := t.Location()
	day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	wd := int(day.Weekday())
	if wd == 0 {
		wd = 7
	}
	return day.AddDate(0, 0, -(wd - 1))
}

func fmtTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func boardsCacheKey(site string) string {
	return "rank_boards_v4_" + site
}

func historyCacheKey(site, typ, period string) string {
	return "rank_history_v4_" + site + "_" + typ + "_" + period
}

func normalizeType(typ string) string {
	switch typ {
	case TypeDuration, TypeCourse, TypeCoins:
		return typ
	default:
		return TypeDuration
	}
}
