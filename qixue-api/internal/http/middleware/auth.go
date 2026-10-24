package middleware

import (
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/response"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/coin"
)

const (
	ctxUID   = "uid"
	ctxToken = "token"
	ctxExp   = "tokenExp"
)

func UID(c *gin.Context) int {
	id, _ := c.Get(ctxUID)
	n, _ := id.(int)
	return n
}

func Token(c *gin.Context) string {
	raw, _ := c.Get(ctxToken)
	s, _ := raw.(string)
	return s
}

func TokenExp(c *gin.Context) int64 {
	raw, _ := c.Get(ctxExp)
	n, _ := raw.(int64)
	return n
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		method := c.GetHeader("Access-Control-Request-Method")
		if method == "" {
			method = "*"
		}
		headers := c.GetHeader("Access-Control-Request-Headers")
		if headers == "" {
			headers = "*"
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", method)
		c.Header("Access-Control-Allow-Headers", headers)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func AdminAuth(rdb *redis.Client, issuer auth.Issuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := bearer(c)
		if !ok {
			return
		}
		id, exp, err := issuer.Parse(raw)
		if err != nil {
			if errors.Is(err, auth.ErrExpired) {
				response.Fail(c, 401, "登录过期，请重新登录")
			} else {
				response.Fail(c, 401, "token无效")
			}
			c.Abort()
			return
		}
		blocked, err := jwtBlocked(c, rdb, raw)
		if err != nil {
			response.Fail(c, 500, "服务器发生异常")
			c.Abort()
			return
		}
		if blocked {
			response.Fail(c, 401, "登录已退出，请重新登录")
			c.Abort()
			return
		}
		c.Set(ctxUID, id)
		c.Set(ctxToken, raw)
		c.Set(ctxExp, exp)
		c.Next()
	}
}

func WebappAuth(db *gorm.DB, rdb *redis.Client, issuer auth.Issuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := bearer(c)
		if !ok {
			return
		}
		if !signOK(c) {
			return
		}
		id, _, err := issuer.Parse(raw)
		if err != nil {
			if errors.Is(err, auth.ErrExpired) {
				response.Fail(c, 401, "登录过期，请重新登录")
			} else {
				response.Fail(c, 401, "token无效")
			}
			c.Abort()
			return
		}
		if userLocked(c, db, rdb, id) {
			response.Fail(c, 401, "账号已被拉黑")
			c.Abort()
			return
		}
		if !skipDailyEnter(c.Request.URL.Path) {
			coin.TryDailyEnter(c.Request.Context(), db, rdb, id)
		}
		c.Set(ctxUID, id)
		c.Next()
	}
}

func skipDailyEnter(path string) bool {
	switch {
	case strings.HasPrefix(path, "/webapp/home/"):
		return true
	case path == "/webapp/sys/getSysInfo":
		return true
	default:
		return false
	}
}

func bearer(c *gin.Context) (string, bool) {
	header := c.GetHeader("Authorization")
	if header == "" {
		response.Fail(c, 401, "需要登录")
		c.Abort()
		return "", false
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		response.Fail(c, 401, "token is error")
		c.Abort()
		return "", false
	}
	return parts[1], true
}

func jwtBlocked(c *gin.Context, rdb *redis.Client, token string) (bool, error) {
	sum := sha256.Sum256([]byte(token))
	key := "qx:jwt_block:" + hex.EncodeToString(sum[:])
	n, err := rdb.Exists(c.Request.Context(), key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func BlockToken(c *gin.Context, rdb *redis.Client, token string, exp int64) error {
	if token == "" {
		return nil
	}
	ttl := time.Until(time.Unix(exp, 0))
	if ttl < time.Second {
		ttl = time.Second
	}
	sum := sha256.Sum256([]byte(token))
	key := "qx:jwt_block:" + hex.EncodeToString(sum[:])
	return rdb.Set(c.Request.Context(), key, "1", ttl).Err()
}

func signOK(c *gin.Context) bool {
	version := c.GetHeader("version-xcx")
	stamp := c.GetHeader("timestmp-xcx")
	random := c.GetHeader("random-xcx")
	sign := c.GetHeader("sign-xcx")
	if version == "" || stamp == "" || random == "" || sign == "" {
		response.Fail(c, 401, "头部参数缺失")
		c.Abort()
		return false
	}
	sum := md5.Sum([]byte("xq" + version + stamp + random))
	if sign != hex.EncodeToString(sum[:]) {
		response.Fail(c, 401, "签名校验失败")
		c.Abort()
		return false
	}
	ms, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil || abs(ms/1000-time.Now().Unix()) > 10 {
		response.Fail(c, 401, "非法访问")
		c.Abort()
		return false
	}
	return true
}

func userLocked(c *gin.Context, db *gorm.DB, rdb *redis.Client, uid int) bool {
	var lock int
	err := cache.Remember(c.Request.Context(), rdb, "user_is_lock_"+strconv.Itoa(uid), &lock, func() (any, error) {
		var value sql.NullInt64
		err := db.WithContext(c.Request.Context()).Raw("SELECT is_lock FROM users WHERE id = ? LIMIT 1", uid).Scan(&value).Error
		if err != nil {
			return nil, err
		}
		if !value.Valid {
			return -1, nil
		}
		return int(value.Int64), nil
	})
	if err != nil {
		return false
	}
	return lock == 1
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
