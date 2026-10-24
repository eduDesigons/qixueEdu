package scope

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
)

type cityCache struct {
	City int `json:"city"`
}

// CityID 返回管理员所属大区。未限定时 ok 为 false。
func CityID(ctx context.Context, db *gorm.DB, rdb *redis.Client, adminID int) (int, bool, error) {
	if adminID <= 0 {
		return 0, false, nil
	}
	var hit cityCache
	err := cache.Remember(ctx, rdb, "admin_"+strconv.Itoa(adminID), &hit, func() (any, error) {
		var city sql.NullInt64
		err := db.WithContext(ctx).Raw("SELECT city_id FROM administrator WHERE id = ? LIMIT 1", adminID).Scan(&city).Error
		if err != nil {
			return nil, err
		}
		if !city.Valid {
			return cityCache{}, nil
		}
		return cityCache{City: int(city.Int64)}, nil
	})
	if err != nil {
		return 0, false, err
	}
	if hit.City == 0 {
		return 0, false, nil
	}
	return hit.City, true, nil
}
