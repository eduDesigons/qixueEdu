package rank

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"
)

const (
	DefaultRefreshIntervalSec = 180
	dailyCheckInterval        = 60 * time.Second
	initialRebuildDelay       = 1 * time.Second
	lockTTL                   = 5 * time.Minute
)

// StartTicker runs in-process rank snapshot timers.
// Interval from RANK_SNAPSHOT_INTERVAL_SEC (default 180). Set to 0 to disable.
func (s *Service) StartTicker(ctx context.Context) {
	sec := envIntervalSec()
	if sec <= 0 {
		slog.Info("rank snapshot ticker disabled", "RANK_SNAPSHOT_INTERVAL_SEC", sec)
		return
	}
	interval := time.Duration(sec) * time.Second
	slog.Info("rank snapshot ticker started",
		"refresh_interval_sec", sec,
		"daily_check_sec", int(dailyCheckInterval.Seconds()),
	)

	go s.runTicker(ctx, interval)
}

func (s *Service) runTicker(ctx context.Context, refreshEvery time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("rank snapshot ticker panic recovered", "recover", r)
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-time.After(initialRebuildDelay):
		s.safeRun(ctx, "rebuild_all", func(c context.Context) error {
			return s.withSnapLock(c, "snap", lockTTL, func(c context.Context) error {
				return s.RebuildAll(c, true)
			})
		})
	}

	refresh := time.NewTicker(refreshEvery)
	daily := time.NewTicker(dailyCheckInterval)
	defer refresh.Stop()
	defer daily.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("rank snapshot ticker stopped")
			return
		case <-refresh.C:
			s.safeRun(ctx, "refresh_duration", func(c context.Context) error {
				return s.withSnapLock(c, "snap", lockTTL, func(c context.Context) error {
					return s.RefreshDurationAll(c)
				})
			})
		case <-daily.C:
			s.safeRun(ctx, "maybe_daily", func(c context.Context) error {
				return s.MaybeDailyRebuild(c)
			})
		}
	}
}

func (s *Service) safeRun(parent context.Context, name string, fn func(context.Context) error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("rank snapshot job panic recovered", "job", name, "recover", r)
		}
	}()
	ctx, cancel := context.WithTimeout(parent, lockTTL)
	defer cancel()
	slog.Info("rank snapshot job start", "job", name)
	if err := fn(ctx); err != nil {
		slog.Error("rank snapshot job failed", "job", name, "err", err)
		return
	}
	slog.Info("rank snapshot job success", "job", name)
}

func envIntervalSec() int {
	v := os.Getenv("RANK_SNAPSHOT_INTERVAL_SEC")
	if v == "" {
		return DefaultRefreshIntervalSec
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return DefaultRefreshIntervalSec
	}
	return n
}
