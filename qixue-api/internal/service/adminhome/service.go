package adminhome

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/bunny"
	"qixue-api/internal/service/scope"
)

const (
	cacheStatsTTL = 45 * time.Second
	cacheChartTTL = 90 * time.Second
	cacheRankTTL  = 90 * time.Second
)

type Service struct {
	home  repo.AdminHomeRepository
	db    *gorm.DB
	rdb   *redis.Client
	stats *bunny.Stats
}

func New(home repo.AdminHomeRepository, db *gorm.DB, rdb *redis.Client) *Service {
	cfg := bunny.LoadConfig()
	return &Service{home: home, db: db, rdb: rdb, stats: bunny.NewStats(cfg, rdb)}
}

func (s *Service) GetVerse(c *gin.Context) {
	ctx := c.Request.Context()
	maxID, err := s.home.MaxVerseID(ctx)
	if response.WriteErr(c, err) {
		return
	}
	if maxID <= 0 {
		response.OK(c, "随机一言", gin.H{"verse": ""})
		return
	}
	startID := rand.Intn(maxID) + 1
	verse, _ := s.home.VerseContentGE(ctx, startID)
	if verse == "" {
		verse, _ = s.home.VerseContentLE(ctx, startID)
	}
	response.OK(c, "随机一言", gin.H{"verse": verse})
}

func (s *Service) GetCdnTraffic(c *gin.Context) {
	payload, err := s.stats.DashboardPayload(c.Request.Context())
	if err != nil {
		response.Fail(c, 1, "CDN 流量获取失败："+err.Error())
		return
	}
	response.OK(c, "CDN 流量", payload)
}

func (s *Service) GetPendingStats(c *gin.Context) {
	data, err := s.buildPendingStats(c)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "待办统计", data)
}

func (s *Service) GetTodayStats(c *gin.Context) {
	data, err := s.remember(c, "today_v2", cacheStatsTTL, func() (map[string]any, error) {
		return s.buildTodayStats(c)
	})
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "今日学情", data)
}

// GetStats returns pending and today dashboard stats.
func (s *Service) GetStats(c *gin.Context) {
	pending, err := s.buildPendingStats(c)
	if response.WriteErr(c, err) {
		return
	}
	today, err := s.remember(c, "today_v2", cacheStatsTTL, func() (map[string]any, error) {
		return s.buildTodayStats(c)
	})
	if response.WriteErr(c, err) {
		return
	}
	merged := make(map[string]any, len(pending)+len(today))
	for k, v := range pending {
		merged[k] = v
	}
	for k, v := range today {
		merged[k] = v
	}
	response.OK(c, "看板统计", merged)
}

func (s *Service) GetChartData(c *gin.Context) {
	data, err := s.remember(c, "chart14", cacheChartTTL, func() (map[string]any, error) {
		return s.buildChartPayload(c)
	})
	if err != nil {
		response.Fail(c, 1, "图表数据获取失败："+err.Error())
		return
	}
	response.OK(c, "图表数据", data)
}

func (s *Service) GetCourseLearningRank(c *gin.Context) {
	list := s.safeList(func() ([]map[string]any, error) {
		return s.rememberList(c, "rank_courses_learning", cacheRankTTL, s.buildCourseRank)
	})
	response.OK(c, "课程学习中排行", list)
}

func (s *Service) GetChapterLearningRank(c *gin.Context) {
	list := s.safeList(func() ([]map[string]any, error) {
		return s.rememberList(c, "rank_chapters_v2", cacheRankTTL, s.buildChapterRank)
	})
	response.OK(c, "近30天完成课节排行", list)
}

func (s *Service) GetChapterCompleteUserRank(c *gin.Context) {
	list := s.safeList(func() ([]map[string]any, error) {
		return s.rememberList(c, "rank_coins_14d_top8", cacheRankTTL, s.buildCreditRank)
	})
	response.OK(c, "近两周积分排行", list)
}

func (s *Service) buildPendingStats(c *gin.Context) (map[string]any, error) {
	start, end := todayRange()
	applyPending, err := s.scopedCount(c, "user_course_apply", "status = 0")
	if err != nil {
		return nil, err
	}
	pathApply := s.safeCount(func() (int, error) {
		return s.scopedCount(c, "user_path_apply", "status = 0")
	})
	paperPending, err := s.scopedCount(c, "user_paper", "is_people_check = 1")
	if err != nil {
		return nil, err
	}
	courseComment, err := s.scopedCount(c, "course_comments", "created_at BETWEEN ? AND ?", start, end)
	if err != nil {
		return nil, err
	}
	chapterComment, err := s.scopedCount(c, "course_chapter_comments", "created_at BETWEEN ? AND ?", start, end)
	if err != nil {
		return nil, err
	}
	feedback := s.safeCount(func() (int, error) {
		return s.scopedCount(c, "user_feedbacks", "status = 0")
	})
	return map[string]any{
		"apply_pending_count":      applyPending,
		"path_apply_pending_count": pathApply,
		"paper_pending_count":      paperPending,
		"course_comment_count":     courseComment,
		"chapter_comment_count":    chapterComment,
		"feedback_pending_count":   feedback,
	}, nil
}

func (s *Service) buildTodayStats(c *gin.Context) (map[string]any, error) {
	start, end := todayRange()
	courseEnd, err := s.scopedCount(c, "user_course", "is_end = 1 AND updated_at BETWEEN ? AND ?", start, end)
	if err != nil {
		return nil, err
	}
	chapterEnd, err := s.scopedCount(c, "user_course_chapter", "is_end = 1 AND updated_at BETWEEN ? AND ?", start, end)
	if err != nil {
		return nil, err
	}
	paperCount, err := s.scopedCount(c, "user_paper", "created_at BETWEEN ? AND ?", start, end)
	if err != nil {
		return nil, err
	}
	paperPass, err := s.scopedCount(c, "user_paper", "created_at BETWEEN ? AND ? AND is_pass = 1", start, end)
	if err != nil {
		return nil, err
	}
	rate := 0
	if paperCount > 0 {
		rate = int(float64(paperPass)/float64(paperCount)*100 + 0.5)
	}
	creditSum, err := s.scopedSum(c, "user_credit_records", "sum", "deleted_at IS NULL AND created_at BETWEEN ? AND ?", start, end)
	if err != nil {
		return nil, err
	}
	coinsSum := s.safeCount(func() (int, error) {
		return s.scopedSum(c, "user_coins_records", "sum", "deleted_at IS NULL AND `sum` > 0 AND created_at BETWEEN ? AND ?", start, end)
	})
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

func (s *Service) cityScope(c *gin.Context) (int, bool, error) {
	return scope.CityID(c.Request.Context(), s.db, s.rdb, middleware.UID(c))
}

func (s *Service) cacheScope(c *gin.Context) string {
	cityID, limited, _ := s.cityScope(c)
	if !limited {
		return "all"
	}
	return "city_" + strconv.Itoa(cityID)
}

func (s *Service) cacheKey(c *gin.Context, suffix string) string {
	return "admin_dash_" + s.cacheScope(c) + "_" + suffix
}

func (s *Service) remember(c *gin.Context, suffix string, ttl time.Duration, builder func() (map[string]any, error)) (map[string]any, error) {
	key := s.cacheKey(c, suffix)
	var hit map[string]any
	if cache.GetJSON(c.Request.Context(), s.rdb, key, &hit) && hit != nil {
		return hit, nil
	}
	data, err := builder()
	if err != nil {
		return nil, err
	}
	cache.SetJSON(c.Request.Context(), s.rdb, key, data, ttl)
	return data, nil
}

func (s *Service) rememberList(c *gin.Context, suffix string, ttl time.Duration, builder func(*gin.Context) ([]map[string]any, error)) ([]map[string]any, error) {
	key := s.cacheKey(c, suffix)
	var hit []map[string]any
	if cache.GetJSON(c.Request.Context(), s.rdb, key, &hit) && hit != nil {
		return hit, nil
	}
	data, err := builder(c)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(c.Request.Context(), s.rdb, key, data, ttl)
	return data, nil
}

func (s *Service) scopedCount(c *gin.Context, table, where string, args ...any) (int, error) {
	cityID, limited, err := s.cityScope(c)
	if err != nil {
		return 0, err
	}
	return s.home.ScopedCount(c.Request.Context(), table, where, cityID, limited, args...)
}

func (s *Service) scopedSum(c *gin.Context, table, column, where string, args ...any) (int, error) {
	cityID, limited, err := s.cityScope(c)
	if err != nil {
		return 0, err
	}
	return s.home.ScopedSum(c.Request.Context(), table, column, where, cityID, limited, args...)
}

func (s *Service) safeCount(fn func() (int, error)) int {
	n, err := fn()
	if err != nil {
		return 0
	}
	return n
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
