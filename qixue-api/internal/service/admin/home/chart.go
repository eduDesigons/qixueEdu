package adminhome

import (
	"context"
	"time"
)

func (s *Service) buildChartPayload(ctx context.Context, adminID int) (map[string]any, error) {
	days := 14
	dates := make([]string, 0, days)
	loc := shanghai()
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	for i := days - 1; i >= 0; i-- {
		dates = append(dates, today.AddDate(0, 0, -i).Format("2006-01-02"))
	}
	rangeStart := today.AddDate(0, 0, -(days - 1)).Format("2006-01-02 15:04:05")
	rangeEnd := today.Add(24*time.Hour - time.Second).Format("2006-01-02 15:04:05")
	zeros := make([]int, days)

	courseEnd := s.safeSeries(func() ([]int, error) {
		m, err := s.queryCountByDate(ctx, adminID, "user_course", "updated_at", rangeStart, rangeEnd, "is_end = 1")
		return s.alignSeries(dates, m, err)
	}, zeros)
	chapterEnd := s.safeSeries(func() ([]int, error) {
		m, err := s.queryCountByDate(ctx, adminID, "user_course_chapter", "updated_at", rangeStart, rangeEnd, "is_end = 1")
		return s.alignSeries(dates, m, err)
	}, zeros)
	paperPass := s.safeSeries(func() ([]int, error) {
		m, err := s.queryCountByDate(ctx, adminID, "user_paper", "created_at", rangeStart, rangeEnd, "is_pass = 1")
		return s.alignSeries(dates, m, err)
	}, zeros)
	creditSum := s.safeSeries(func() ([]int, error) {
		m, err := s.querySumByDate(ctx, adminID, "user_credit_records", "created_at", "sum", rangeStart, rangeEnd, "deleted_at IS NULL")
		return s.alignSeries(dates, m, err)
	}, zeros)

	return map[string]any{
		"line": map[string]any{
			"dates":       dates,
			"course_end":  courseEnd,
			"chapter_end": chapterEnd,
			"paper_pass":  paperPass,
			"credit_sum":  creditSum,
		},
	}, nil
}

func (s *Service) safeSeries(fn func() ([]int, error), fallback []int) []int {
	vals, err := fn()
	if err != nil || vals == nil {
		out := make([]int, len(fallback))
		copy(out, fallback)
		return out
	}
	return vals
}

func (s *Service) queryCountByDate(ctx context.Context, adminID int, table, dateColumn, start, end, extra string) (map[string]int, error) {
	cityID, limited, err := s.cityScope(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return s.home.CountByDate(ctx, table, dateColumn, start, end, extra, cityID, limited)
}

func (s *Service) querySumByDate(ctx context.Context, adminID int, table, dateColumn, sumColumn, start, end, extra string) (map[string]int, error) {
	cityID, limited, err := s.cityScope(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return s.home.SumByDate(ctx, table, dateColumn, sumColumn, start, end, extra, cityID, limited)
}

func (s *Service) alignSeries(dates []string, dateToValue map[string]int, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}
	out := make([]int, len(dates))
	for i, d := range dates {
		out[i] = dateToValue[d]
	}
	return out, nil
}
