package adminhome

import (
	"context"
	"time"

	"qixue-api/internal/platform/async"
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
	cityID, limited, err := s.cityScope(ctx, adminID)
	if err != nil {
		return nil, err
	}

	series := func(m map[string]int, err error) []int {
		return s.safeSeries(func() ([]int, error) { return s.alignSeries(dates, m, err) }, zeros)
	}
	var courseEnd, chapterEnd, paperPass, creditSum []int
	_ = async.Run(
		func() error {
			m, callErr := s.home.CountByDate(ctx, "user_course", "updated_at", rangeStart, rangeEnd, "is_end = 1", cityID, limited)
			courseEnd = series(m, callErr)
			return nil
		},
		func() error {
			m, callErr := s.home.CountByDate(ctx, "user_course_chapter", "updated_at", rangeStart, rangeEnd, "is_end = 1", cityID, limited)
			chapterEnd = series(m, callErr)
			return nil
		},
		func() error {
			m, callErr := s.home.CountByDate(ctx, "user_paper", "created_at", rangeStart, rangeEnd, "is_pass = 1", cityID, limited)
			paperPass = series(m, callErr)
			return nil
		},
		func() error {
			m, callErr := s.home.SumByDate(ctx, "user_credit_records", "created_at", "sum", rangeStart, rangeEnd, "deleted_at IS NULL", cityID, limited)
			creditSum = series(m, callErr)
			return nil
		},
	)

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
