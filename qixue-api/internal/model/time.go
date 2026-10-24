package model

import "time"

func FormatTime(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.Format("2006-01-02 15:04:05")
	return &s
}
