package site

import (
	"context"
	"net"
	"regexp"
	"sort"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
)

var (
	consoleHost = regexp.MustCompile(`^console\.(.+)$`)
	staticHost  = regexp.MustCompile(`^([^.]+)\.static\.(.+)$`)
	appHost     = regexp.MustCompile(`^([^.]+)\.(h5|admin|console)\.(.+)$`)
	siteCode    = regexp.MustCompile(`^[a-z0-9]+$`)
)

func Codes(ctx context.Context, db *gorm.DB, rdb *redis.Client) ([]string, error) {
	var out []string
	err := cache.Remember(ctx, rdb, "area_site_codes", &out, func() (any, error) {
		return loadCodes(ctx, db)
	})
	if out == nil {
		out = []string{}
	}
	return out, err
}

func loadCodes(ctx context.Context, db *gorm.DB) ([]string, error) {
	var codes []string
	err := db.WithContext(ctx).Raw(
		"SELECT code FROM areas WHERE parent_id = 0 AND deleted_at IS NULL AND code IS NOT NULL AND code <> ''",
	).Scan(&codes).Error
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.ToLower(strings.TrimSpace(code))
		if code == "" || !siteCode.MatchString(code) {
			continue
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		out = append(out, code)
	}
	sort.Strings(out)
	return out, nil
}

// RequestCode 从 Host 推导站点 city_code
func RequestCode(host string, codes []string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if host == "" {
		return ""
	}
	site := ""
	if m := staticHost.FindStringSubmatch(host); m != nil {
		site = m[1]
	} else if m := appHost.FindStringSubmatch(host); m != nil {
		site = m[1]
	}
	if site != "" && contains(codes, site) {
		return site
	}
	return ""
}

func AssetDomain(host, fallback string, codes []string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if host == "" {
		return ""
	}
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		return normalize(fallback)
	}
	assetHost := ""
	if m := consoleHost.FindStringSubmatch(host); m != nil {
		assetHost = "static." + m[1]
	} else {
		site, suffix := "", ""
		if m := staticHost.FindStringSubmatch(host); m != nil {
			site, suffix = m[1], m[2]
		} else if m := appHost.FindStringSubmatch(host); m != nil {
			site, suffix = m[1], m[3]
		}
		if site != "" && contains(codes, site) {
			assetHost = site + ".static." + suffix
		}
	}
	if assetHost == "" {
		return ""
	}
	return normalize(assetHost)
}

func normalize(domain string) string {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return ""
	}
	domain = strings.TrimRight(domain, "/")
	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		return domain
	}
	return "https://" + domain
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}
