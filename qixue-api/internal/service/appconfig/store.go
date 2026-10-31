package appconfig

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"qixue-api/internal/platform/cache"
)

type stored struct {
	V any `json:"v"`
}

func Get(ctx context.Context, db *gorm.DB, rdb *redis.Client, key string) (any, error) {
	var hit stored
	err := cache.Remember(ctx, rdb, "app_config_"+key, &hit, func() (any, error) {
		var raw sql.NullString
		err := db.WithContext(ctx).Raw("SELECT `value` FROM app_config WHERE `key` = ? LIMIT 1", key).Scan(&raw).Error
		if err != nil {
			return nil, err
		}
		if !raw.Valid {
			return stored{V: nil}, nil
		}
		return stored{V: decode(raw.String)}, nil
	})
	if err != nil {
		return nil, err
	}
	return hit.V, nil
}

func Save(ctx context.Context, db *gorm.DB, rdb *redis.Client, key string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	now := time.Now()
	result := db.WithContext(ctx).Exec("UPDATE app_config SET `value` = ?, updated_at = ? WHERE `key` = ?", string(encoded), now, key)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		err = db.WithContext(ctx).Exec(
			"INSERT INTO app_config (`key`, `value`, created_at, updated_at) VALUES (?, ?, ?, ?)",
			key, string(encoded), now, now,
		).Error
		if err != nil {
			return err
		}
	}
	names := []string{"app_config_" + key}
	if key == "sysConfig" || key == "playerConfig" || key == "coinConfig" {
		names = append(names, "sys_info_boot")
	}
	cache.Drop(ctx, rdb, names...)
	return nil
}

func decode(raw string) any {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return nil
	}
	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return raw
	}
	if text, ok := value.(string); ok {
		var again any
		if json.Unmarshal([]byte(text), &again) == nil {
			return again
		}
	}
	return value
}
