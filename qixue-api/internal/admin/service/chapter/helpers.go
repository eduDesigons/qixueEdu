package adminchapter

import (
	"context"
	"strconv"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/service/scope"
)

// BizError is a business failure mapped to HTTP Fail(code=1) by the handler.
type BizError struct{ Message string }

func (e *BizError) Error() string { return e.Message }

func biz(msg string) error { return &BizError{Message: msg} }

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

func (s *Service) resolveCity(ctx context.Context, adminID, requestCityID int) (int, error) {
	cityID, limited, err := scope.CityID(ctx, s.db, s.rdb, adminID)
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
