package applog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DailyWriter writes to dir/prefix-YYYY-MM-DD.log and drops files older than keepDays.
type DailyWriter struct {
	dir      string
	prefix   string
	keepDays int

	mu   sync.Mutex
	day  string
	file *os.File
}

func NewDailyWriter(dir, prefix string, keepDays int) (*DailyWriter, error) {
	if dir == "" {
		dir = "logs"
	}
	if prefix == "" {
		prefix = "qixue-api"
	}
	if keepDays < 1 {
		keepDays = 7
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	w := &DailyWriter{dir: dir, prefix: prefix, keepDays: keepDays}
	if err := w.rotateLocked(time.Now()); err != nil {
		return nil, err
	}
	w.cleanupLocked(time.Now())
	return w, nil
}

func (w *DailyWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	day := now.Format("2006-01-02")
	if day != w.day {
		if err := w.rotateLocked(now); err != nil {
			return 0, err
		}
		w.cleanupLocked(now)
	}
	return w.file.Write(p)
}

func (w *DailyWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

func (w *DailyWriter) Path() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return ""
	}
	return w.file.Name()
}

func (w *DailyWriter) rotateLocked(now time.Time) error {
	day := now.Format("2006-01-02")
	path := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, day))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if w.file != nil {
		_ = w.file.Close()
	}
	w.file = f
	w.day = day
	return nil
}

func (w *DailyWriter) cleanupLocked(now time.Time) {
	// 保留含今天在内的 keepDays 个自然日，更早的删掉。
	oldest := now.AddDate(0, 0, -(w.keepDays - 1)).Format("2006-01-02")
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}
	prefix := w.prefix + "-"
	suffix := ".log"
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, suffix) {
			continue
		}
		day := strings.TrimSuffix(strings.TrimPrefix(name, prefix), suffix)
		if len(day) != 10 || day >= oldest {
			continue
		}
		_ = os.Remove(filepath.Join(w.dir, name))
	}
}
