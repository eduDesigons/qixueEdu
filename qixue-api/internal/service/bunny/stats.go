package bunny

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"qixue-api/internal/platform/async"
	"qixue-api/internal/platform/cache"
)

const (
	statsAPIBase  = "https://api.bunny.net"
	statsCacheTTL = 5 * time.Minute
	statsCacheVer = "v11"
	statsTZ       = "Asia/Shanghai"
)

type Stats struct {
	cfg  Config
	rdb  *redis.Client
	http *http.Client
}

func NewStats(cfg Config, rdb *redis.Client) *Stats {
	return &Stats{
		cfg:  cfg,
		rdb:  rdb,
		http: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *Stats) DashboardPayload(ctx context.Context) (map[string]any, error) {
	if !s.cfg.StatsReady() {
		return map[string]any{
			"enabled": false,
			"message": "未配置 BUNNY_ACCOUNT_API_KEY",
		}, nil
	}
	cacheKey := "bunny_cdn_dash_" + statsCacheVer + "_account"
	var cached map[string]any
	if cache.GetJSON(ctx, s.rdb, cacheKey, &cached) {
		if month, ok := cached["month"].(map[string]any); ok {
			if dates, ok := month["dates"].([]any); ok && len(dates) > 0 {
				return cached, nil
			}
		}
	}

	loc, err := time.LoadLocation(statsTZ)
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	monthFrom := today.AddDate(0, 0, -29)
	todayKey := today.Format("2006-01-02")
	yesterdayKey := today.AddDate(0, 0, -1).Format("2006-01-02")
	dayBeforeKey := today.AddDate(0, 0, -2).Format("2006-01-02")

	var (
		monthRaw map[string]any
		errMsg   string
		hourly   hourlyCompare
	)
	_ = async.Run(
		func() error {
			monthRaw, errMsg = s.fetchStatistics(monthFrom, now.Add(24*time.Hour-time.Second), false)
			return nil
		},
		func() error {
			hourly = s.fetchCompareHourly(today, now, loc)
			return nil
		},
	)
	if monthRaw == nil {
		msg := errMsg
		if msg == "" {
			msg = "Bunny 统计接口请求失败"
		}
		return map[string]any{"enabled": false, "message": msg}, nil
	}
	bwByDay := chartByShanghaiDay(monthRaw["BandwidthUsedChart"], loc)
	reqByDay := chartByShanghaiDay(monthRaw["RequestsServedChart"], loc)

	dates := make([]string, 0, 30)
	bandwidth := make([]int, 0, 30)
	requests := make([]int, 0, 30)
	for d := monthFrom; !d.After(today); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		dates = append(dates, key)
		bandwidth = append(bandwidth, bwByDay[key])
		requests = append(requests, reqByDay[key])
	}
	sumInt := func(vals []any) int {
		total := 0
		for _, v := range vals {
			switch n := v.(type) {
			case int:
				total += n
			case float64:
				total += int(n)
			}
		}
		return total
	}
	todayBwList := make([]any, len(hourly.todayBandwidth))
	todayReqList := make([]any, len(hourly.todayRequests))
	for i := range hourly.todayBandwidth {
		if hourly.todayBandwidth[i] == nil {
			todayBwList[i] = nil
			todayReqList[i] = nil
		} else {
			todayBwList[i] = *hourly.todayBandwidth[i]
			todayReqList[i] = *hourly.todayRequests[i]
		}
	}
	todayBandwidth := sumInt(todayBwList)
	if todayBandwidth == 0 {
		todayBandwidth = bwByDay[todayKey]
	}
	todayRequests := sumInt(todayReqList)
	if todayRequests == 0 {
		todayRequests = reqByDay[todayKey]
	}

	payload := map[string]any{
		"enabled":  true,
		"scope":    "account",
		"timezone": statsTZ,
		"today": map[string]any{
			"date":                        todayKey,
			"yesterday_date":              yesterdayKey,
			"day_before_date":             dayBeforeKey,
			"bandwidth":                   todayBandwidth,
			"requests":                    todayRequests,
			"hours":                       hourly.hours,
			"bandwidth_hourly":            todayBwList,
			"requests_hourly":             todayReqList,
			"bandwidth_yesterday_hourly":  hourly.yesterdayBandwidth,
			"requests_yesterday_hourly":   hourly.yesterdayRequests,
			"bandwidth_day_before_hourly": hourly.dayBeforeBandwidth,
			"requests_day_before_hourly":  hourly.dayBeforeRequests,
		},
		"month": map[string]any{
			"from":            monthFrom.Format("2006-01-02"),
			"to":              today.Format("2006-01-02"),
			"total_bandwidth": sumInts(bandwidth),
			"total_requests":  sumInts(requests),
			"dates":           dates,
			"bandwidth":       bandwidth,
			"requests":        requests,
		},
		"fetched_at": now.Format("2006-01-02 15:04:05"),
	}
	cache.SetJSON(ctx, s.rdb, cacheKey, payload, statsCacheTTL)
	return payload, nil
}

type hourlyCompare struct {
	hours              []string
	todayBandwidth     []*int
	todayRequests      []*int
	yesterdayBandwidth []int
	yesterdayRequests  []int
	dayBeforeBandwidth []int
	dayBeforeRequests  []int
}

func (s *Stats) fetchCompareHourly(todayStart, now time.Time, loc *time.Location) hourlyCompare {
	hours := make([]string, 24)
	for h := 0; h <= 23; h++ {
		hours[h] = fmt.Sprintf("%02d:00", h)
	}
	empty := make([]int, 24)
	endHour := now.Hour()
	dayBefore := todayStart.AddDate(0, 0, -2)
	yesterdayKey := todayStart.AddDate(0, 0, -1).Format("2006-01-02")
	dayBeforeKey := dayBefore.Format("2006-01-02")
	todayKey := todayStart.Format("2006-01-02")

	raw, _ := s.fetchStatistics(dayBefore, now.Truncate(time.Hour).Add(time.Hour-time.Second), true)
	out := hourlyCompare{
		hours:              hours,
		todayBandwidth:     make([]*int, 24),
		todayRequests:      make([]*int, 24),
		yesterdayBandwidth: empty,
		yesterdayRequests:  append([]int(nil), empty...),
		dayBeforeBandwidth: append([]int(nil), empty...),
		dayBeforeRequests:  append([]int(nil), empty...),
	}
	if raw == nil {
		for h := endHour + 1; h <= 23; h++ {
			out.todayBandwidth[h] = nil
			out.todayRequests[h] = nil
		}
		for h := 0; h <= endHour; h++ {
			z := 0
			out.todayBandwidth[h] = &z
			out.todayRequests[h] = &z
		}
		return out
	}
	bwToday := chartByShanghaiHour(raw["BandwidthUsedChart"], todayKey, loc)
	reqToday := chartByShanghaiHour(raw["RequestsServedChart"], todayKey, loc)
	bwY := chartByShanghaiHour(raw["BandwidthUsedChart"], yesterdayKey, loc)
	reqY := chartByShanghaiHour(raw["RequestsServedChart"], yesterdayKey, loc)
	bwB := chartByShanghaiHour(raw["BandwidthUsedChart"], dayBeforeKey, loc)
	reqB := chartByShanghaiHour(raw["RequestsServedChart"], dayBeforeKey, loc)

	yBw := make([]int, 24)
	yReq := make([]int, 24)
	bBw := make([]int, 24)
	bReq := make([]int, 24)
	for h := 0; h <= 23; h++ {
		yBw[h] = bwY[h]
		yReq[h] = reqY[h]
		bBw[h] = bwB[h]
		bReq[h] = reqB[h]
		if h <= endHour {
			v := bwToday[h]
			r := reqToday[h]
			out.todayBandwidth[h] = &v
			out.todayRequests[h] = &r
		} else {
			out.todayBandwidth[h] = nil
			out.todayRequests[h] = nil
		}
	}
	out.yesterdayBandwidth = yBw
	out.yesterdayRequests = yReq
	out.dayBeforeBandwidth = bBw
	out.dayBeforeRequests = bReq
	return out
}

func (s *Stats) fetchStatistics(from, to time.Time, hourly bool) (map[string]any, string) {
	fromUTC := from.UTC()
	toUTC := to.UTC()
	q := url.Values{}
	q.Set("dateFrom", fromUTC.Format("2006-01-02T15:04:05Z"))
	q.Set("dateTo", toUTC.Format("2006-01-02T15:04:05Z"))
	q.Set("loadBandwidthUsed", "true")
	q.Set("loadRequestsServed", "true")
	if hourly {
		q.Set("hourly", "true")
		q.Set("exactRange", "true")
	}
	req, err := http.NewRequest(http.MethodGet, statsAPIBase+"/statistics?"+q.Encode(), nil)
	if err != nil {
		return nil, "Bunny 统计请求异常：" + err.Error()
	}
	req.Header.Set("AccessKey", s.cfg.AccountAPIKey)
	req.Header.Set("Accept", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, "Bunny 统计请求异常：" + err.Error()
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	var jsonMap map[string]any
	if err := json.Unmarshal(body, &jsonMap); err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := ""
		if jsonMap != nil {
			if m, ok := jsonMap["Message"].(string); ok {
				msg = m
			}
		}
		if msg != "" {
			return nil, "Bunny 统计失败：" + msg
		}
		return nil, fmt.Sprintf("Bunny 统计失败：HTTP %d", resp.StatusCode)
	}
	return jsonMap, ""
}

func chartByShanghaiDay(chart any, loc *time.Location) map[string]int {
	out := map[string]int{}
	m, ok := chart.(map[string]any)
	if !ok {
		return out
	}
	for ts, val := range m {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			t, err = time.ParseInLocation("2006-01-02T15:04:05Z", ts, time.UTC)
			if err != nil {
				continue
			}
		}
		day := t.In(loc).Format("2006-01-02")
		out[day] += toInt(val)
	}
	return out
}

func chartByShanghaiHour(chart any, dayKey string, loc *time.Location) map[int]int {
	out := map[int]int{}
	m, ok := chart.(map[string]any)
	if !ok {
		return out
	}
	for ts, val := range m {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			t, err = time.ParseInLocation("2006-01-02T15:04:05Z", ts, time.UTC)
			if err != nil {
				continue
			}
		}
		local := t.In(loc)
		if local.Format("2006-01-02") != dayKey {
			continue
		}
		out[local.Hour()] += toInt(val)
	}
	return out
}

func toInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	case string:
		i, _ := strconv.Atoi(strings.TrimSpace(n))
		return i
	default:
		return 0
	}
}

func sumInts(vals []int) int {
	n := 0
	for _, v := range vals {
		n += v
	}
	return n
}
