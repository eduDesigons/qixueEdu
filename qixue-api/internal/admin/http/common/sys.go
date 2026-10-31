package common

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/response"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/service/appconfig"
	"qixue-api/internal/service/site"
)

type Sys struct {
	db          *gorm.DB
	rdb         *redis.Client
	assetDomain string
}

func NewSys(db *gorm.DB, rdb *redis.Client, assetDomain string) *Sys {
	return &Sys{db: db, rdb: rdb, assetDomain: assetDomain}
}

func (h *Sys) Boot(c *gin.Context) {
	ctx := c.Request.Context()
	var data struct {
		Sys    any `json:"sysConfig"`
		Player any `json:"playerConfig"`
		Coins  any `json:"coinConfig"`
	}
	err := cache.Remember(ctx, h.rdb, "sys_info_boot", &data, func() (any, error) {
		sys, err := appconfig.Get(ctx, h.db, h.rdb, "sysConfig")
		if err != nil {
			return nil, err
		}
		player, err := appconfig.Get(ctx, h.db, h.rdb, "playerConfig")
		if err != nil {
			return nil, err
		}
		coins, err := appconfig.Get(ctx, h.db, h.rdb, "coinConfig")
		if err != nil {
			return nil, err
		}
		return struct {
			Sys    any `json:"sysConfig"`
			Player any `json:"playerConfig"`
			Coins  any `json:"coinConfig"`
		}{Sys: sys, Player: player, Coins: coins}, nil
	})
	if response.WriteErr(c, err) {
		return
	}
	sysMap, _ := data.Sys.(map[string]any)
	if sysMap == nil {
		sysMap = map[string]any{}
	}
	codes, err := site.Codes(ctx, h.db, h.rdb)
	if response.WriteErr(c, err) {
		return
	}
	sysMap["assetDomain"] = site.AssetDomain(c.Request.Host, h.assetDomain, codes)
	response.OK(c, "获取成功", gin.H{
		"sysConfig":    sysMap,
		"playerConfig": data.Player,
		"coinConfig":   data.Coins,
	})
}
