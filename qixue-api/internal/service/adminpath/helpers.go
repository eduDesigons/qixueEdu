package adminpath

import (
	"html"
	"strconv"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/service/scope"
)

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	default:
		return ""
	}
}

func asInt(v any) int { return bind.AsInt(v) }

func toIntSlice(v any) []int {
	switch list := v.(type) {
	case []any:
		out := make([]int, 0, len(list))
		for _, item := range list {
			if n := asInt(item); n > 0 {
				out = append(out, n)
			}
		}
		return out
	case []int:
		return list
	case []float64:
		out := make([]int, 0, len(list))
		for _, item := range list {
			if n := int(item); n > 0 {
				out = append(out, n)
			}
		}
		return out
	default:
		return nil
	}
}

func uniquePositive(ids []int) []int {
	seen := map[int]struct{}{}
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func keys(set map[int]struct{}) []int {
	out := make([]int, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func (s *Service) resolveCity(c *gin.Context, requestCityID int) (int, error) {
	cityID, limited, err := scope.CityID(c.Request.Context(), s.db, s.rdb, middleware.UID(c))
	if err != nil {
		return 0, err
	}
	if limited {
		return cityID, nil
	}
	if requestCityID > 0 {
		return requestCityID, nil
	}
	return 0, nil
}

func escapeTitle(title string) string {
	return html.EscapeString(title)
}

var allowedSort = map[string]bool{
	"id": true, "title": true, "created_at": true, "updated_at": true,
	"is_show": true, "is_rec": true, "is_open": true, "category_id": true,
	"min_need_credit": true,
}
