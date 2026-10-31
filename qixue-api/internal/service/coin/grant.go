package coin

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
)

const (
	fieldTest   = "通过测评"
	fieldCourse = "课程结业"
	fieldWatch  = "观看视频"
	fieldRead   = "阅读图文"
)

// OnLearnEnd grants coins once per chapter (watch/read) when first completed.
func OnLearnEnd(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, chapterID, chapterType int, title string) {
	if uid <= 0 || chapterID <= 0 {
		return
	}
	isMedia := chapterType == 1 || chapterType == 2
	field := fieldRead
	sum := amount(ctx, db, rdb, "read", 5)
	prefix := "读完课节「"
	if isMedia {
		field = fieldWatch
		sum = amount(ctx, db, rdb, "watch", 10)
		prefix = "看完课节「"
	}
	if title == "" {
		title = "课节"
	}
	uniq := "#c" + strconv.Itoa(chapterID) + "#"
	grantOnce(ctx, repo.NewCoin(db), rdb, uid, field, sum, prefix+title+"」 "+uniq, uniq)
}

func OnPassTest(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, paperID int, title string) {
	if title == "" {
		title = "测评"
	}
	uniq := "#p" + strconv.Itoa(paperID) + "#"
	grantOnce(ctx, repo.NewCoin(db), rdb, uid, fieldTest, amount(ctx, db, rdb, "test", 10), "通过测评「"+title+"」 "+uniq, uniq)
}

func OnCourseEnd(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, courseID int, title string) {
	if title == "" {
		title = "课程"
	}
	uniq := "#course" + strconv.Itoa(courseID) + "#"
	grantOnce(ctx, repo.NewCoin(db), rdb, uid, fieldCourse, amount(ctx, db, rdb, "course", 20), "完成课程「"+title+"」 "+uniq, uniq)
}

func AddCredit(ctx context.Context, db *gorm.DB, uid int, field string, credit int, remark string) error {
	return repo.NewCoin(db).AddCredit(ctx, uid, field, credit, remark)
}

func AddCoins(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid int, field string, coins int, remark string) error {
	if uid <= 0 || coins <= 0 {
		return nil
	}
	if err := repo.NewCoin(db).AddCoins(ctx, uid, field, coins, remark); err != nil {
		return err
	}
	cache.Drop(ctx, rdb, "user_my_count_"+strconv.Itoa(uid))
	return nil
}

func grantOnce(ctx context.Context, coins repo.CoinRepository, rdb *redis.Client, uid int, field string, sum int, remark, uniq string) {
	if uid <= 0 || sum <= 0 || uniq == "" {
		return
	}
	count, err := coins.CountFieldRemarkLike(ctx, uid, field, "%"+uniq+"%")
	if err != nil {
		slog.Error("积分发放查询失败", "user_id", uid, "field", field, "err", err)
		return
	}
	if count > 0 {
		return
	}
	if err := coins.AddCoins(ctx, uid, field, sum, remark); err != nil {
		slog.Error("积分发放失败", "user_id", uid, "field", field, "err", err)
		return
	}
	cache.Drop(ctx, rdb, "user_my_count_"+strconv.Itoa(uid))
}
