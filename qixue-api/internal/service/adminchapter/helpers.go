package adminchapter

import (
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
