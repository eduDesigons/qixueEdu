package pms

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/response"
	"qixue-api/internal/service/appconfig"
)

type Config struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewConfig(db *gorm.DB, rdb *redis.Client) *Config {
	return &Config{db: db, rdb: rdb}
}

func (h *Config) Get(c *gin.Context) {
	key := c.Query("type")
	if key == "" {
		response.Fail(c, 1, "参数错误")
		return
	}
	value, err := appconfig.Get(c.Request.Context(), h.db, h.rdb, key)
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "获取成功", value)
}

func (h *Config) Save(c *gin.Context) {
	key := c.Query("type")
	if key == "" {
		response.Fail(c, 1, "参数错误")
		return
	}
	var body any
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	if response.WriteErr(c, appconfig.Save(c.Request.Context(), h.db, h.rdb, key, body)) {
		return
	}
	response.OK(c, "保存成功", nil)
}
