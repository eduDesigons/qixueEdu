package rank

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"qixue-api/internal/platform/cache"
	"qixue-api/internal/repo"
)

const (
	Limit    = 5
	CacheTTL = 24 * time.Hour
	Site     = "all"

	TypeDuration = "duration"
	TypeCourse   = "course"
	TypeCoins    = "coins"
)

type Service struct {
	ranks repo.RankRepository
	rdb   *redis.Client
}

func New(ranks repo.RankRepository, rdb *redis.Client) *Service {
	return &Service{ranks: ranks, rdb: rdb}
}

type coinAgg struct {
	ID    int
	Value int
}

// Boards returns current week/month boards (read from snapshots; ticker writes them).
func (s *Service) Boards(ctx context.Context, uid int) (map[string]any, error) {
	payload, err := s.cachedBoards(ctx)
	if err != nil {
		return nil, err
	}
	payload = s.attachBoardCoins(ctx, payload)
	return s.withMe(payload, uid), nil
}

// History returns past period boards for a type (excludes current period key).
func (s *Service) History(ctx context.Context, typ, period string, uid int) ([]map[string]any, error) {
	typ = normalizeType(typ)
	if period != "month" {
		period = "week"
	}
	list, err := s.cachedHistory(ctx, typ, period)
	if err != nil {
		return nil, err
	}
	var coins map[int]int
	if typ != TypeCourse {
		coins = s.coinMap(ctx, collectUserIDs(list))
	}
	for i := range list {
		rows := asRowMaps(list[i]["list"])
		if typ != TypeCourse {
			rows = fillCoins(rows, coins)
		}
		list[i]["list"] = markMe(rows, uid, typ)
	}
	return list, nil
}

func (s *Service) cachedBoards(ctx context.Context) (map[string]any, error) {
	key := boardsCacheKey(Site)
	var payload map[string]any
	if cache.GetJSON(ctx, s.rdb, key, &payload) {
		if _, ok := payload["duration"]; ok {
			return payload, nil
		}
	}
	payload, err := s.assembleFromSnapshots(ctx)
	if err != nil {
		return nil, err
	}
	cache.SetJSON(ctx, s.rdb, key, payload, CacheTTL)
	return payload, nil
}

func (s *Service) cachedHistory(ctx context.Context, typ, period string) ([]map[string]any, error) {
	key := historyCacheKey(Site, typ, period)
	var list []map[string]any
	if cache.GetJSON(ctx, s.rdb, key, &list) {
		return list, nil
	}
	currentKey := periodMeta(period, "current").Key
	rows, err := s.ranks.ListSnapshotsByPeriod(ctx, Site, typ, period)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	out := make([]map[string]any, 0)
	for _, row := range rows {
		if row.PeriodKey == currentKey {
			continue
		}
		if _, ok := seen[row.PeriodKey]; ok {
			continue
		}
		items := decodeList(row.Payload)
		if len(items) == 0 {
			continue
		}
		seen[row.PeriodKey] = struct{}{}
		out = append(out, map[string]any{
			"periodKey":   row.PeriodKey,
			"periodLabel": row.PeriodLabel,
			"updatedAt":   row.SnapshotDate.Format("2006-01-02"),
			"list":        items,
		})
	}
	cache.SetJSON(ctx, s.rdb, key, out, CacheTTL)
	return out, nil
}

func (s *Service) assembleFromSnapshots(ctx context.Context) (map[string]any, error) {
	week := periodMeta("week", "current")
	month := periodMeta("month", "current")
	rows, err := s.ranks.ListSnapshotsByPeriodKeys(ctx, Site, []string{week.Key, month.Key})
	if err != nil {
		return nil, err
	}
	best := map[string]repo.RankSnapRow{}
	for _, row := range rows {
		slot := ""
		switch row.PeriodKey {
		case week.Key:
			slot = "week"
		case month.Key:
			slot = "month"
		default:
			continue
		}
		k := row.Type + "|" + slot
		if _, ok := best[k]; ok {
			continue
		}
		best[k] = row
	}
	out := map[string]any{}
	for _, typ := range []string{TypeDuration, TypeCourse, TypeCoins} {
		out[typ] = map[string]any{
			"week":  s.boardFromSnap(ctx, best[typ+"|week"], typ, week),
			"month": s.boardFromSnap(ctx, best[typ+"|month"], typ, month),
		}
	}
	out["tips"] = map[string]any{
		"duration": "按全部学员本周/本月累计学习时长排名，完成后台定时更新。",
		"course":   "按课程本周/本月累计学习人数排名，每日 00:00 更新。",
		"coins":    "按学员本周/本月获得的积分排名，定时更新。",
	}
	return out, nil
}

func (s *Service) boardFromSnap(ctx context.Context, snap repo.RankSnapRow, typ string, meta periodInfo) map[string]any {
	hasSnap := snap.Type != "" || snap.Payload != "" || !snap.SnapshotDate.IsZero()
	list := []map[string]any{}
	if hasSnap {
		list = decodeList(snap.Payload)
	}
	if typ == TypeCoins && len(list) == 0 {
		list = s.computeCoins(ctx, meta.Start, meta.LiveEnd)
	}
	var updated any
	liveFallback := len(list) > 0 && typ == TypeCoins && (!hasSnap || snap.Payload == "" || snap.Payload == "[]")
	if liveFallback {
		updated = meta.LiveEnd
	} else if hasSnap {
		if snap.UpdatedAt != nil && !snap.UpdatedAt.IsZero() {
			updated = snap.UpdatedAt.Format("2006-01-02 15:04:05")
		} else if !snap.SnapshotDate.IsZero() {
			updated = snap.SnapshotDate.Format("2006-01-02") + " 00:00:00"
		}
	}
	return map[string]any{
		"list":        list,
		"updatedAt":   updated,
		"realtime":    typ == TypeDuration || typ == TypeCoins,
		"periodLabel": meta.Label,
	}
}

func (s *Service) computeCoins(ctx context.Context, start, end string) []map[string]any {
	rows, err := s.ranks.AggregateCoins(ctx, start, end, Limit)
	if err != nil || len(rows) == 0 {
		return []map[string]any{}
	}
	aggs := make([]coinAgg, 0, len(rows))
	for _, r := range rows {
		aggs = append(aggs, coinAgg{ID: r.ID, Value: r.Value})
	}
	return s.hydrateUsers(ctx, aggs, "积分")
}

func (s *Service) hydrateUsers(ctx context.Context, rows []coinAgg, unit string) []map[string]any {
	ids := make([]int, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	users := map[int]repo.RankUserBrief{}
	if len(ids) > 0 {
		list, _ := s.ranks.ListUsersBrief(ctx, ids)
		for _, u := range list {
			users[u.ID] = u
		}
	}
	out := make([]map[string]any, 0, len(rows))
	for i, r := range rows {
		name := "学员"
		avatar := ""
		if u, ok := users[r.ID]; ok {
			if u.NickName != nil && *u.NickName != "" {
				name = *u.NickName
			}
			avatar = u.Avatar
		}
		out = append(out, map[string]any{
			"rank": i + 1, "id": r.ID, "name": name, "avatar": avatar,
			"value": r.Value, "unit": unit,
		})
	}
	return out
}
