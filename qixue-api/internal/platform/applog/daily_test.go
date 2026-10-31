package applog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCleanupKeepsSevenCalendarDays(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 10, 22, 12, 0, 0, 0, time.Local)
	for i := 0; i < 12; i++ {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		path := filepath.Join(dir, "qixue-api-"+day+".log")
		if err := os.WriteFile(path, []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	w := &DailyWriter{dir: dir, prefix: "qixue-api", keepDays: 7}
	w.cleanupLocked(now)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 7 {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("want 7 log files, got %d: %v", len(entries), names)
	}
	oldest := now.AddDate(0, 0, -6).Format("2006-01-02")
	for _, e := range entries {
		day := e.Name()[len("qixue-api-") : len(e.Name())-len(".log")]
		if day < oldest {
			t.Fatalf("stale file kept: %s", e.Name())
		}
	}
}

func TestDailyWriterWrite(t *testing.T) {
	dir := t.TempDir()
	w, err := NewDailyWriter(dir, "qixue-api", 7)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	if _, err := w.Write([]byte("hello\n")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(w.Path())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello\n" {
		t.Fatalf("got %q", data)
	}
}
