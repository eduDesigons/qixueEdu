package middleware

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
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
	"qixue-api/internal/platform/bg"
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
		c.Set(ctxUID, id)
		bg.Go(c.Request.Context(), func(ctx context.Context) {
			coin.TryDailyEnter(ctx, db, rdb, id)
		})
		c.Next()
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
	delta := ms/1000 - time.Now().Unix()
	if delta < 0 {
		delta = -delta
	}
	if err != nil || delta > 30 {
		response.Fail(c, 401, "非法访问")
		c.Abort()
		return false
	}
	return true
}
