package cert

import (
	"testing"
	"time"
)

func TestResolveBoxType(t *testing.T) {
	cases := []struct {
		box  certBox
		want string
	}{
		{certBox{Type: "field"}, "field"},
		{certBox{Type: "text"}, "text"},
		{certBox{Type: "stamp"}, "stamp"},
		{certBox{Type: "image"}, "stamp"},
		{certBox{Key: "text_1"}, "text"},
		{certBox{Key: "stamp_seal"}, "stamp"},
		{certBox{Key: "image_logo"}, "stamp"},
		{certBox{Key: "name"}, "field"},
	}
	for _, tc := range cases {
		if got := resolveBoxType(tc.box); got != tc.want {
			t.Fatalf("resolveBoxType(%+v)=%q want %q", tc.box, got, tc.want)
		}
	}
}

func TestInterpolate(t *testing.T) {
	vals := map[string]string{"name": "张三", "course": "Go入门", "time": "2026年3月1日", "no": "QX1"}
	got := interpolate("学员{name}完成{course}，编号｛no｝", vals)
	want := "学员张三完成Go入门，编号QX1"
	if got != want {
		t.Fatalf("interpolate got %q want %q", got, want)
	}
}

func TestNormalizeCertSpaces(t *testing.T) {
	in := "a\u3000b\u00a0c"
	got := normalizeCertSpaces(in)
	want := "a   b c"
	if got != want {
		t.Fatalf("normalizeCertSpaces got %q want %q", got, want)
	}
}

func TestFormatCertDate(t *testing.T) {
	tm := time.Date(2026, 3, 5, 12, 0, 0, 0, time.UTC)
	got := formatCertDate(tm)
	// Asia/Shanghai = UTC+8 → still 2026-03-05
	if got != "2026年3月5日" {
		t.Fatalf("formatCertDate got %q", got)
	}
}

func TestFindFont(t *testing.T) {
	if path := findFont(); path == "" {
		t.Log("no font found (CI may lack fonts); render will skip text")
	} else {
		t.Log("font:", path)
	}
}
