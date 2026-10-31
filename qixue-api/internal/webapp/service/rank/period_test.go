package rank

import (
	"fmt"
	"testing"
	"time"
)

func TestStartOfWeekMondayISO(t *testing.T) {
	// 2026-03-18 is Wednesday; week starts Monday 2026-03-16 → ISO 2026-W12.
	now := time.Date(2026, 3, 18, 15, 30, 0, 0, time.Local)
	start := startOfWeekMonday(now)
	if start.Year() != 2026 || start.Month() != time.March || start.Day() != 16 {
		t.Fatalf("week start = %v", start)
	}
	y, w := start.ISOWeek()
	key := fmt.Sprintf("%d-W%02d", y, w)
	if key != "2026-W12" {
		t.Fatalf("key = %q", key)
	}
}

func TestNormalizeType(t *testing.T) {
	if normalizeType("coins") != TypeCoins {
		t.Fatal("coins")
	}
	if normalizeType("nope") != TypeDuration {
		t.Fatal("default")
	}
}

func TestEnvIntervalDefault(t *testing.T) {
	t.Setenv("RANK_SNAPSHOT_INTERVAL_SEC", "")
	if envIntervalSec() != DefaultRefreshIntervalSec {
		t.Fatalf("default = %d", envIntervalSec())
	}
	t.Setenv("RANK_SNAPSHOT_INTERVAL_SEC", "0")
	if envIntervalSec() != 0 {
		t.Fatal("want 0 to disable")
	}
	t.Setenv("RANK_SNAPSHOT_INTERVAL_SEC", "90")
	if envIntervalSec() != 90 {
		t.Fatal("want 90")
	}
}
