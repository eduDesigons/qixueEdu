package coin

import (
	"context"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/repo"
)

const fieldComment = "发表评论"

var spaceRe = regexp.MustCompile(`\s+`)
var tagRe = regexp.MustCompile(`(?s)<[^>]*>`)

// OnComment grants coins for top-level comments (pid==0).
func OnComment(ctx context.Context, db *gorm.DB, rdb *redis.Client, uid, pid int, content string) {
	if uid <= 0 || pid > 0 {
		return
	}
	text := spaceRe.ReplaceAllString(tagRe.ReplaceAllString(content, ""), "")
	text = strings.TrimSpace(text)
	sum := amount(ctx, db, rdb, "comment", 5)
	capN := amount(ctx, db, rdb, "commentDailyCap", 3)
	minLen := amount(ctx, db, rdb, "commentMinLen", 4)
	if sum <= 0 || capN <= 0 || utf8.RuneCountInString(text) < minLen {
		return
	}
	day := time.Now().Format("2006-01-02")
	start := day + " 00:00:00"
	end := day + " 23:59:59"
	coins := repo.NewCoin(db)
	todayCount, err := coins.CountFieldInRange(ctx, uid, fieldComment, start, end)
	if err != nil {
		slog.Error("积分-评论查询失败", "user_id", uid, "err", err)
		return
	}
	if int(todayCount) >= capN {
		return
	}
	ymd := strings.ReplaceAll(day, "-", "")
	remark := "发表有效评论 #" + ymd + "-" + strconv.FormatInt(todayCount+1, 10) + "#"
	if err := AddCoins(ctx, db, rdb, uid, fieldComment, sum, remark); err != nil {
		slog.Error("积分-评论失败", "user_id", uid, "err", err)
	}
}
