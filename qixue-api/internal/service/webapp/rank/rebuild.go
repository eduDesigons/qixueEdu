package rank

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
)

// RebuildAll refreshes current + previous period snapshots for all board types
func (s *Service) RebuildAll(ctx context.Context, refreshDuration bool) error {
	if err := s.ranks.EnsureSnapshotsTable(ctx); err != nil {
		return err
	}
	return s.rebuildSite(ctx, Site, refreshDuration)
}

// RefreshDurationAll force-refreshes current week/month snapshots for all types
func (s *Service) RefreshDurationAll(ctx context.Context) error {
	if err := s.ranks.EnsureSnapshotsTable(ctx); err != nil {
		return err
	}
	for _, typ := range []string{TypeDuration, TypeCourse, TypeCoins} {
		if err := s.persistSnapshot(ctx, Site, typ, "week", "current", true); err != nil {
			return err
		}
		if err := s.persistSnapshot(ctx, Site, typ, "month", "current", true); err != nil {
			return err
		}
	}
	payload, err := s.assembleFromSnapshots(ctx)
	if err != nil {
		return err
	}
	cache.SetJSON(ctx, s.rdb, boardsCacheKey(Site), payload, CacheTTL)
	return nil
}

// MaybeDailyRebuild runs a full rebuild once in the first 10 minutes after midnight
func (s *Service) MaybeDailyRebuild(ctx context.Context) error {
	now := time.Now()
	if now.Hour() != 0 || now.Minute() > 10 {
		return nil
	}
	flag := "rank_daily_" + now.Format("20060102")
	var exists int
	if cache.GetJSON(ctx, s.rdb, flag, &exists) && exists == 1 {
		return nil
	}
	return s.withSnapLock(ctx, "snap", 5*time.Minute, func(c context.Context) error {
		if cache.GetJSON(c, s.rdb, flag, &exists) && exists == 1 {
			return nil
		}
		cache.SetJSON(c, s.rdb, flag, 1, 24*time.Hour)
		return s.RebuildAll(c, true)
	})
}

func (s *Service) rebuildSite(ctx context.Context, site string, refreshDuration bool) error {
	site = trimSite(site)
	if site == "" {
		return nil
	}
	for _, typ := range []string{TypeDuration, TypeCourse, TypeCoins} {
		for _, period := range []string{"week", "month"} {
			if err := s.persistSnapshot(ctx, site, typ, period, "current", refreshDuration); err != nil {
				return err
			}
			if err := s.persistSnapshot(ctx, site, typ, period, "previous", false); err != nil {
				return err
			}
		}
	}
	payload, err := s.assembleFromSnapshots(ctx)
	if err != nil {
		return err
	}
	cache.SetJSON(ctx, s.rdb, boardsCacheKey(site), payload, CacheTTL)
	for _, typ := range []string{TypeDuration, TypeCourse, TypeCoins} {
		for _, period := range []string{"week", "month"} {
			cache.Drop(ctx, s.rdb, historyCacheKey(site, typ, period))
		}
	}
	return nil
}

func (s *Service) persistSnapshot(ctx context.Context, siteCode, typ, period, which string, force bool) error {
	siteCode = trimSite(siteCode)
	if siteCode == "" {
		return nil
	}
	meta := periodMeta(period, which)

	row, err := s.ranks.FindSnapshot(ctx, siteCode, typ, period, meta.Key, meta.SnapshotDate)
	if err != nil {
		return err
	}
	if row != nil && !force && len(decodeList(row.Payload)) > 0 {
		return nil
	}

	end := meta.End
	if which == "current" {
		end = meta.LiveEnd
	}
	list := s.compute(ctx, typ, meta.Start, end)
	payload, err := json.Marshal(list)
	if err != nil {
		return err
	}
	now := time.Now()
	if row != nil {
		return s.ranks.UpdateSnapshot(ctx, row.ID, meta.Label, string(payload), now)
	}
	createErr := s.ranks.CreateSnapshot(ctx, repo.RankSnapshotCreate{
		SiteCode:     siteCode,
		Type:         typ,
		Period:       period,
		PeriodKey:    meta.Key,
		PeriodLabel:  meta.Label,
		SnapshotDate: meta.SnapshotDate,
		Payload:      string(payload),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	// Unique race: another writer may have inserted the same snapshot.
	if createErr != nil {
		return nil
	}
	return nil
}

func (s *Service) compute(ctx context.Context, typ, start, end string) []map[string]any {
	switch typ {
	case TypeCourse:
		return s.computeCourse(ctx, start, end)
	case TypeCoins:
		return s.computeCoins(ctx, start, end)
	default:
		return s.computeDuration(ctx, start, end)
	}
}

func (s *Service) computeDuration(ctx context.Context, start, end string) []map[string]any {
	rows, err := s.ranks.AggregateDuration(ctx, start, end, Limit)
	if err != nil || len(rows) == 0 {
		return []map[string]any{}
	}
	aggs := make([]coinAgg, 0, len(rows))
	for _, r := range rows {
		mins := r.Value / 60
		if mins < 1 {
			mins = 1
		}
		aggs = append(aggs, coinAgg{ID: r.ID, Value: mins})
	}
	return s.hydrateUsers(ctx, aggs, "分钟")
}

func (s *Service) computeCourse(ctx context.Context, start, end string) []map[string]any {
	rows, err := s.ranks.AggregateCourse(ctx, start, end, Limit)
	if err != nil || len(rows) == 0 {
		return []map[string]any{}
	}
	ids := make([]int, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	courses := map[int]repo.RankCourseBrief{}
	list, _ := s.ranks.ListCoursesBrief(ctx, ids)
	for _, c := range list {
		courses[c.ID] = c
	}
	out := make([]map[string]any, 0, Limit)
	rank := 0
	for _, r := range rows {
		c, ok := courses[r.ID]
		if !ok || c.IsShow != 1 {
			continue
		}
		rank++
		out = append(out, map[string]any{
			"rank": rank, "id": r.ID, "name": c.Title, "thumb": c.Thumb,
			"value": r.Value, "unit": "人在学",
		})
		if rank >= Limit {
			break
		}
	}
	return out
}

func trimSite(site string) string {
	return strings.TrimSpace(site)
}

func (s *Service) withSnapLock(ctx context.Context, name string, ttl time.Duration, fn func(context.Context) error) error {
	if s.rdb == nil {
		return fn(ctx)
	}
	key := "qx:rank_lock:" + name
	ok, err := s.rdb.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return fmt.Errorf("rank lock: %w", err)
	}
	if !ok {
		return nil
	}
	defer func() { _ = s.rdb.Del(context.Background(), key).Err() }()
	return fn(ctx)
}
