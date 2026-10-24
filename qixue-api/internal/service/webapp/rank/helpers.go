package rank

import (
	"context"
	"encoding/json"

	"qixue-api/internal/http/bind"
)

func decodeList(payload string) []map[string]any {
	if payload == "" || payload == "null" {
		return []map[string]any{}
	}
	var raw any
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return []map[string]any{}
	}
	return asRowMaps(raw)
}

func asRowMaps(v any) []map[string]any {
	switch rows := v.(type) {
	case []map[string]any:
		return rows
	case []any:
		out := make([]map[string]any, 0, len(rows))
		for _, item := range rows {
			if m, ok := item.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	default:
		return []map[string]any{}
	}
}

func collectUserIDs(blocks []map[string]any) []int {
	ids := make([]int, 0)
	seen := map[int]struct{}{}
	for _, block := range blocks {
		for _, row := range asRowMaps(block["list"]) {
			id := asInt(row["id"])
			if id <= 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids
}

func (s *Service) coinMap(ctx context.Context, ids []int) map[int]int {
	out := map[int]int{}
	if len(ids) == 0 {
		return out
	}
	rows, _ := s.ranks.ListUserCoins(ctx, ids)
	for _, r := range rows {
		out[r.ID] = r.Coins
	}
	return out
}

func fillCoins(list []map[string]any, coins map[int]int) []map[string]any {
	for i := range list {
		id := asInt(list[i]["id"])
		if v, ok := coins[id]; ok {
			list[i]["coins"] = v
		} else if list[i]["coins"] == nil {
			list[i]["coins"] = 0
		}
	}
	return list
}

func (s *Service) attachBoardCoins(ctx context.Context, payload map[string]any) map[string]any {
	blocks := make([]map[string]any, 0, 4)
	for _, typ := range []string{TypeDuration, TypeCoins} {
		block, _ := payload[typ].(map[string]any)
		if block == nil {
			continue
		}
		for _, period := range []string{"week", "month"} {
			board, _ := block[period].(map[string]any)
			if board == nil {
				continue
			}
			blocks = append(blocks, map[string]any{"list": board["list"]})
		}
	}
	coins := s.coinMap(ctx, collectUserIDs(blocks))
	for _, typ := range []string{TypeDuration, TypeCoins} {
		block, _ := payload[typ].(map[string]any)
		if block == nil {
			continue
		}
		for _, period := range []string{"week", "month"} {
			board, _ := block[period].(map[string]any)
			if board == nil {
				continue
			}
			board["list"] = fillCoins(asRowMaps(board["list"]), coins)
			block[period] = board
		}
		payload[typ] = block
	}
	return payload
}

func (s *Service) withMe(payload map[string]any, uid int) map[string]any {
	for _, typ := range []string{TypeDuration, TypeCourse, TypeCoins} {
		block, _ := payload[typ].(map[string]any)
		if block == nil {
			continue
		}
		for _, period := range []string{"week", "month"} {
			board, _ := block[period].(map[string]any)
			if board == nil {
				continue
			}
			board["list"] = markMe(asRowMaps(board["list"]), uid, typ)
			block[period] = board
		}
		payload[typ] = block
	}
	return payload
}

func markMe(list []map[string]any, uid int, typ string) []map[string]any {
	out := make([]map[string]any, 0, len(list))
	for _, item := range list {
		copy := map[string]any{}
		for k, v := range item {
			copy[k] = v
		}
		copy["isMe"] = typ != TypeCourse && uid > 0 && asInt(item["id"]) == uid
		out = append(out, copy)
	}
	return out
}

func asInt(v any) int {
	return bind.AsInt(v)
}
