package coin

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
	"qixue-api/internal/service/appconfig"
)

const fieldDaily = "每日登录"

func TryDailyEnter(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid int) {
	if uid <= 0 || rdb == nil {
		return
	}
	day := time.Now().Format("2006-01-02")
	key := "qx:coin_daily:" + strconv.Itoa(uid) + ":" + strings.ReplaceAll(day, "-", "")
	ok, err := rdb.SetNX(ctx, key, "1", 26*time.Hour).Result()
	if err != nil || !ok {
		return
	}
	if err := grantDaily(ctx, repo.NewCoin(db), db, rdb, uid, day); err != nil {
		slog.Error("积分-每日进入失败", "user_id", uid, "err", err)
		_ = rdb.Del(ctx, key).Err()
	}
}

func grantDaily(ctx context.Context, coins repo.CoinRepository, db *gorm.DB, rdb *redis.Client, uid int, day string) error {
	start := day + " 00:00:00"
	end := day + " 23:59:59"
	count, err := coins.CountFieldInRange(ctx, uid, fieldDaily, start, end)
	if err != nil {
		return err
	}
	sum := amount(ctx, db, rdb, "daily", 5)
	if count > 0 || sum <= 0 {
		return nil
	}
	if err := coins.AddCoins(ctx, uid, fieldDaily, sum, "每日首次进入 #"+day+"#"); err != nil {
		return err
	}
	cache.Drop(ctx, rdb, "user_my_count_"+strconv.Itoa(uid))
	return nil
}

func amount(ctx context.Context, db *gorm.DB, rdb *redis.Client, key string, fallback int) int {
	raw, err := appconfig.Get(ctx, db, rdb, "coinConfig")
	if err != nil || raw == nil {
		return fallback
	}
	values, _ := raw.(map[string]any)
	if values == nil {
		return fallback
	}
	value, ok := values[key]
	if !ok || value == nil || value == "" {
		return fallback
	}
	switch n := value.(type) {
	case float64:
		if n < 0 {
			return 0
		}
		return int(n)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(n))
		if err != nil || parsed < 0 {
			return 0
		}
		return parsed
	default:
		return fallback
	}
}
