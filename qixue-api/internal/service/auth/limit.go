package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const failLimit = 5

type Limiter struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewLimiter(rdb *redis.Client) *Limiter {
	return &Limiter{rdb: rdb, ttl: 15 * time.Minute}
}

func (l *Limiter) Blocked(ctx context.Context, scope, account string) bool {
	return l.Count(ctx, scope, account) >= failLimit
}

func (l *Limiter) Count(ctx context.Context, scope, account string) int {
	n, err := l.rdb.Get(ctx, failKey(scope, account)).Int()
	if err != nil {
		return 0
	}
	return n
}

func (l *Limiter) Record(ctx context.Context, scope, account string) {
	l.RecordTTL(ctx, scope, account, l.ttl)
}

func (l *Limiter) RecordTTL(ctx context.Context, scope, account string, ttl time.Duration) {
	if ttl <= 0 {
		ttl = l.ttl
	}
	key := failKey(scope, account)
	n, err := l.rdb.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return
	}
	_ = l.rdb.Set(ctx, key, n+1, ttl).Err()
}

func (l *Limiter) Clear(ctx context.Context, scope, account string) {
	_ = l.rdb.Del(ctx, failKey(scope, account)).Err()
}

func failKey(scope, account string) string {
	account = strings.ToLower(strings.TrimSpace(account))
	if len(account) > 64 {
		account = account[:64]
	}
	sum := sha256.Sum256([]byte(account))
	return "qx:login_fail:" + scope + ":" + hex.EncodeToString(sum[:])
}
