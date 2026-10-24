package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Remember loads a JSON value from Redis (qx:name) or calls load and stores it.
func Remember(ctx context.Context, rdb *redis.Client, name string, dest any, load func() (any, error)) error {
	return RememberTTL(ctx, rdb, name, dest, 0, load)
}

// RememberTTL is Remember with an optional TTL; ttl<=0 means no expiry.
func RememberTTL(ctx context.Context, rdb *redis.Client, name string, dest any, ttl time.Duration, load func() (any, error)) error {
	if rdb != nil {
		raw, err := rdb.Get(ctx, "qx:"+name).Bytes()
		if err == nil && json.Unmarshal(raw, dest) == nil {
			return nil
		}
	}
	value, err := load()
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(encoded, dest); err != nil {
		return err
	}
	if rdb != nil {
		_ = rdb.Set(ctx, "qx:"+name, encoded, ttl).Err()
	}
	return nil
}

// GetJSON reads qx:name; returns false on miss.
func GetJSON(ctx context.Context, rdb *redis.Client, name string, dest any) bool {
	if rdb == nil {
		return false
	}
	raw, err := rdb.Get(ctx, "qx:"+name).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(raw, dest) == nil
}

// SetJSON writes qx:name.
func SetJSON(ctx context.Context, rdb *redis.Client, name string, value any, ttl time.Duration) {
	if rdb == nil {
		return
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = rdb.Set(ctx, "qx:"+name, encoded, ttl).Err()
}

// Drop deletes qx:name and the bare name key for each entry.
func Drop(ctx context.Context, rdb *redis.Client, names ...string) {
	if rdb == nil || len(names) == 0 {
		return
	}
	keys := make([]string, 0, len(names)*2)
	for _, name := range names {
		keys = append(keys, "qx:"+name, name)
	}
	_ = rdb.Del(ctx, keys...).Err()
}
