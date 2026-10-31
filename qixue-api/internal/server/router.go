package server

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/admin/http/common"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/platform/config"
	"qixue-api/internal/service/auth"
	"qixue-api/internal/service/captcha"
)

func New(cfg config.Config, db *gorm.DB, rdb *redis.Client, slide *captcha.Slide) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.LoggerWithFormatter(accessLog), middleware.CORS(), gzip.Gzip(gzip.DefaultCompression))

	adminJWT := auth.NewIssuer(cfg.AdminJWT, cfg.JWTExpire)
	userJWT := auth.NewIssuer(cfg.UserJWT, cfg.JWTExpire)

	info := common.NewSys(db, rdb, cfg.AssetDomain)
	r.GET("/admin/common/sys/getSysInfo", info.Boot)
	r.GET("/webapp/sys/getSysInfo", info.Boot)

	registerAdmin(r, db, rdb, slide, adminJWT)
	registerWebapp(r, db, rdb, slide, userJWT)
	return r
}

// accessLog 访问日志：不记录客户端 IP。
func accessLog(p gin.LogFormatterParams) string {
	path := p.Path
	if p.ErrorMessage != "" {
		return fmt.Sprintf("[GIN] %s | %3d | %13v | %-7s %s | %s\n",
			p.TimeStamp.Format("2006/01/02 - 15:04:05"),
			p.StatusCode,
			p.Latency,
			p.Method,
			path,
			p.ErrorMessage,
		)
	}
	return fmt.Sprintf("[GIN] %s | %3d | %13v | %-7s %s\n",
		p.TimeStamp.Format("2006/01/02 - 15:04:05"),
		p.StatusCode,
		p.Latency,
		p.Method,
		path,
	)
}
