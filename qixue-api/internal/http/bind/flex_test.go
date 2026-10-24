package bind

import (
	"encoding/json"
	"testing"
)

func TestIntUnmarshalJSON(t *testing.T) {
	cases := []struct {
		raw  string
		want int
	}{
		{`73`, 73},
		{`"73"`, 73},
		{`" 73 "`, 73},
		{`73.0`, 73},
		{`null`, 0},
		{`""`, 0},
		{`"x"`, 0},
		{`true`, 1},
		{`false`, 0},
		{`"true"`, 1},
		{`"false"`, 0},
	}
	for _, tc := range cases {
		var v Int
		if err := json.Unmarshal([]byte(tc.raw), &v); err != nil {
			t.Fatalf("unmarshal %s: %v", tc.raw, err)
		}
		if v.Int() != tc.want {
			t.Fatalf("%s => %d, want %d", tc.raw, v.Int(), tc.want)
		}
	}
}

func TestAsInt(t *testing.T) {
	if AsInt("42") != 42 || AsInt(float64(7)) != 7 || AsInt(nil) != 0 {
		t.Fatal("AsInt mismatch")
	}
	if AsInt64("9007199254740993") != 9007199254740993 {
		t.Fatal("AsInt64 mismatch")
	}
}

func TestIntUnmarshalText(t *testing.T) {
	var v Int
	if err := v.UnmarshalText([]byte("73")); err != nil || v.Int() != 73 {
		t.Fatalf("UnmarshalText: %v %d", err, v.Int())
	}
}
