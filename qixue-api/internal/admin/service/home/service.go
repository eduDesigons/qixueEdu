package adminhome

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/async"
	"qixue-api/internal/platform/cache"
	adminrepo "qixue-api/internal/admin/repo"
	"qixue-api/internal/service/bunny"
	"qixue-api/internal/service/scope"
)

const (
	cacheStatsTTL = 45 * time.Second
	cacheChartTTL = 90 * time.Second
	cacheRankTTL  = 90 * time.Second
)

type Service struct {
	home  adminrepo.AdminHomeRepository
	db    *gorm.DB
	rdb   *redis.Client
	stats *bunny.Stats
}

func New(home adminrepo.AdminHomeRepository, db *gorm.DB, rdb *redis.Client) *Service {
	cfg := bunny.LoadConfig()
	return &Service{home: home, db: db, rdb: rdb, stats: bunny.NewStats(cfg, rdb)}
}

type VerseView struct {
	Verse string `json:"verse"`
}

func (s *Service) GetVerse(ctx context.Context) (VerseView, error) {
	maxID, err := s.home.MaxVerseID(ctx)
	if err != nil {
		return VerseView{}, err
	}
	if maxID <= 0 {
		return VerseView{Verse: ""}, nil
	}
	startID := rand.Intn(maxID) + 1
	verse, _ := s.home.VerseContentGE(ctx, startID)
	if verse == "" {
		verse, _ = s.home.VerseContentLE(ctx, startID)
	}
	return VerseView{Verse: verse}, nil
}

func (s *Service) GetCdnTraffic(ctx context.Context) (map[string]any, error) {
	return s.stats.DashboardPayload(ctx)
}

func (s *Service) GetPendingStats(ctx context.Context, adminID int) (map[string]any, error) {
	return s.buildPendingStats(ctx, adminID)
}

func (s *Service) GetTodayStats(ctx context.Context, adminID int) (map[string]any, error) {
	return s.remember(ctx, adminID, "today_v2", cacheStatsTTL, func() (map[string]any, error) {
		return s.buildTodayStats(ctx, adminID)
	})
}

// GetStats returns pending and today dashboard stats.
func (s *Service) GetStats(ctx context.Context, adminID int) (map[string]any, error) {
	var pending, today map[string]any
	if err := async.Run(
		func() error {
			var callErr error
			pending, callErr = s.buildPendingStats(ctx, adminID)
			return callErr
		},
		func() error {
			var callErr error
			today, callErr = s.remember(ctx, adminID, "today_v2", cacheStatsTTL, func() (map[string]any, error) {
				return s.buildTodayStats(ctx, adminID)
			})
			return callErr
		},
	); err != nil {
		return nil, err
	}
	merged := make(map[string]any, len(pending)+len(today))
	for k, v := range pending {
		merged[k] = v
	}
	for k, v := range today {
		merged[k] = v
	}
	return merged, nil
}

func (s *Service) GetChartData(ctx context.Context, adminID int) (map[string]any, error) {
	return s.remember(ctx, adminID, "chart14", cacheChartTTL, func() (map[string]any, error) {
		return s.buildChartPayload(ctx, adminID)
	})
}

func (s *Service) GetCourseLearningRank(ctx context.Context, adminID int) []map[string]any {
	return s.safeList(func() ([]map[string]any, error) {
		return s.rememberList(ctx, adminID, "rank_courses_learning", cacheRankTTL, s.buildCourseRank)
	})
}

func (s *Service) GetChapterLearningRank(ctx context.Context, adminID int) []map[string]any {
	return s.safeList(func() ([]map[string]any, error) {
		return s.rememberList(ctx, adminID, "rank_chapters_v2", cacheRankTTL, s.buildChapterRank)
	})
}

func (s *Service) GetChapterCompleteUserRank(ctx context.Context, adminID int) []map[string]any {
	return s.safeList(func() ([]map[string]any, error) {
		return s.rememberList(ctx, adminID, "rank_coins_14d_top8", cacheRankTTL, s.buildCreditRank)
	})
}

func (s *Service) buildPendingStats(ctx context.Context, adminID int) (map[string]any, error) {
	start, end := todayRange()
	cityID, limited, err := s.cityScope(ctx, adminID)
	if err != nil {
		return nil, err
	}
	var applyPending, pathApply, paperPending, courseComment, chapterComment, feedback int
	if err := async.Run(
		func() error {
			var callErr error
			applyPending, callErr = s.home.ScopedCount(ctx, "user_course_apply", "status = 0", cityID, limited)
			return callErr
		},
		func() error {
			pathApply, _ = s.home.ScopedCount(ctx, "user_path_apply", "status = 0", cityID, limited)
			return nil
		},
		func() error {
			var callErr error
			paperPending, callErr = s.home.ScopedCount(ctx, "user_paper", "is_people_check = 1", cityID, limited)
			return callErr
		},
		func() error {
			var callErr error
			courseComment, callErr = s.home.ScopedCount(ctx, "course_comments", "created_at BETWEEN ? AND ?", cityID, limited, start, end)
			return callErr
		},
		func() error {
			var callErr error
			chapterComment, callErr = s.home.ScopedCount(ctx, "course_chapter_comments", "created_at BETWEEN ? AND ?", cityID, limited, start, end)
			return callErr
		},
		func() error {
			feedback, _ = s.home.ScopedCount(ctx, "user_feedbacks", "status = 0", cityID, limited)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return map[string]any{
		"apply_pending_count":      applyPending,
		"path_apply_pending_count": pathApply,
		"paper_pending_count":      paperPending,
		"course_comment_count":     courseComment,
		"chapter_comment_count":    chapterComment,
		"feedback_pending_count":   feedback,
	}, nil
}

func (s *Service) buildTodayStats(ctx context.Context, adminID int) (map[string]any, error) {
	start, end := todayRange()
	cityID, limited, err := s.cityScope(ctx, adminID)
	if err != nil {
		return nil, err
	}
	var courseEnd, chapterEnd, paperCount, paperPass, creditSum, coinsSum int
	if err := async.Run(
		func() error {
			var callErr error
			courseEnd, callErr = s.home.ScopedCount(ctx, "user_course", "is_end = 1 AND updated_at BETWEEN ? AND ?", cityID, limited, start, end)
			return callErr
		},
		func() error {
			var callErr error
			chapterEnd, callErr = s.home.ScopedCount(ctx, "user_course_chapter", "is_end = 1 AND updated_at BETWEEN ? AND ?", cityID, limited, start, end)
			return callErr
		},
		func() error {
			var callErr error
			paperCount, callErr = s.home.ScopedCount(ctx, "user_paper", "created_at BETWEEN ? AND ?", cityID, limited, start, end)
			return callErr
		},
		func() error {
			var callErr error
			paperPass, callErr = s.home.ScopedCount(ctx, "user_paper", "created_at BETWEEN ? AND ? AND is_pass = 1", cityID, limited, start, end)
			return callErr
		},
		func() error {
			var callErr error
			creditSum, callErr = s.home.ScopedSum(ctx, "user_credit_records", "sum", "deleted_at IS NULL AND created_at BETWEEN ? AND ?", cityID, limited, start, end)
			return callErr
		},
		func() error {
			coinsSum, _ = s.home.ScopedSum(ctx, "user_coins_records", "sum", "deleted_at IS NULL AND `sum` > 0 AND created_at BETWEEN ? AND ?", cityID, limited, start, end)
			return nil
		},
	); err != nil {
		return nil, err
	}
	rate := 0
	if paperCount > 0 {
		rate = int(float64(paperPass)/float64(paperCount)*100 + 0.5)
	}
	return map[string]any{
		"today_course_end_count":  courseEnd,
		"today_chapter_end_count": chapterEnd,
		"today_paper_count":       paperCount,
		"today_paper_pass_count":  paperPass,
		"today_paper_pass_rate":   rate,
		"today_credit_sum":        creditSum,
		"today_coins_sum":         coinsSum,
	}, nil
}

func (s *Service) cityScope(ctx context.Context, adminID int) (int, bool, error) {
	return scope.CityID(ctx, s.db, s.rdb, adminID)
}

func (s *Service) cacheScope(ctx context.Context, adminID int) string {
	cityID, limited, _ := s.cityScope(ctx, adminID)
	if !limited {
		return "all"
	}
	return "city_" + strconv.Itoa(cityID)
}

func (s *Service) cacheKey(ctx context.Context, adminID int, suffix string) string {
	return "admin_dash_" + s.cacheScope(ctx, adminID) + "_" + suffix
}

func (s *Service) remember(ctx context.Context, adminID int, suffix string, ttl time.Duration, builder func() (map[string]any, error)) (map[string]any, error) {
	key := s.cacheKey(ctx, adminID, suffix)
	var hit map[string]any
	if cache.GetJSON(ctx, s.rdb, key, &hit) && hit != nil {
		return hit, nil
	}
	data, err := builder()
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, s.rdb, key, data, ttl)
	return data, nil
}

func (s *Service) rememberList(ctx context.Context, adminID int, suffix string, ttl time.Duration, builder func(context.Context, int) ([]map[string]any, error)) ([]map[string]any, error) {
	key := s.cacheKey(ctx, adminID, suffix)
	var hit []map[string]any
	if cache.GetJSON(ctx, s.rdb, key, &hit) && hit != nil {
		return hit, nil
	}
	data, err := builder(ctx, adminID)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, s.rdb, key, data, ttl)
	return data, nil
}

func (s *Service) safeList(fn func() ([]map[string]any, error)) []map[string]any {
	list, err := fn()
	if err != nil || list == nil {
		return []map[string]any{}
	}
	return list
}

func todayRange() (string, string) {
	loc := shanghai()
	now := time.Now().In(loc)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24*time.Hour - time.Second)
	return start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05")
}

func shanghai() *time.Location {
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("CST", 8*3600)
	}
	return l
}
