package sys

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/http/response"
	"qixue-api/internal/service/appconfig"
)

var httpURL = regexp.MustCompile(`(?i)^https?://`)

type Handler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func New(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{db: db, rdb: rdb}
}

func (h *Handler) GetQuestion(c *gin.Context) {
	data, err := appconfig.Get(c.Request.Context(), h.db, h.rdb, "faqConfig")
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "获取成功", data)
}

func (h *Handler) About(c *gin.Context) {
	data, err := appconfig.Get(c.Request.Context(), h.db, h.rdb, "aboutConfig")
	if response.WriteErr(c, err) {
		return
	}
	response.OK(c, "获取成功", data)
}

func (h *Handler) GetTools(c *gin.Context) {
	raw, err := appconfig.Get(c.Request.Context(), h.db, h.rdb, "toolConfig")
	if response.WriteErr(c, err) {
		return
	}
	list, ok := raw.([]any)
	if !ok {
		response.OK(c, "获取成功", []any{})
		return
	}
	out := make([]gin.H, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		title := strings.TrimSpace(asString(m["title"]))
		url := strings.TrimSpace(asString(m["url"]))
		if title == "" || !httpURL.MatchString(url) {
			continue
		}
		out = append(out, gin.H{
			"title": title,
			"desc":  strings.TrimSpace(asString(m["desc"])),
			"url":   url,
			"icon":  strings.TrimSpace(asString(m["icon"])),
		})
	}
	response.OK(c, "获取成功", out)
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		return ""
	}
}
