package bind

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Int unmarshals JSON numbers or numeric strings (route/query 常把 id 当字符串传来).
type Int int

func (v *Int) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*v = 0
		return nil
	}
	if bytes.Equal(b, []byte("true")) {
		*v = 1
		return nil
	}
	if bytes.Equal(b, []byte("false")) {
		*v = 0
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		s = strings.TrimSpace(s)
		if s == "true" || s == "1" {
			*v = 1
			return nil
		}
		if s == "false" || s == "0" || s == "" {
			*v = 0
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			*v = 0
			return nil
		}
		*v = Int(n)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		i, err := n.Int64()
		if err != nil {
			f, ferr := n.Float64()
			if ferr != nil {
				*v = 0
				return nil
			}
			*v = Int(f)
			return nil
		}
		*v = Int(i)
		return nil
	}
	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		*v = 0
		return nil
	}
	*v = Int(f)
	return nil
}

func (v Int) Int() int { return int(v) }

// Ints converts a slice of FlexInt values to plain ints.
func Ints(xs []Int) []int {
	out := make([]int, len(xs))
	for i, x := range xs {
		out[i] = x.Int()
	}
	return out
}

// UnmarshalText lets gin form/query binding accept numeric strings.
func (v *Int) UnmarshalText(text []byte) error {
	*v = Int(Atoi(string(text)))
	return nil
}

// Atoi parses a trimmed decimal string; empty/invalid → 0.
func Atoi(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

// AsInt coerces common JSON/map/SQL values to int (string/float/[]byte/int*).
func AsInt(v any) int {
	return int(AsInt64(v))
}

// AsInt64 coerces common JSON/map/SQL values to int64 (string/float/[]byte/int*).
// Prefer this for bigint columns (e.g. media_file.id) when reading request ids.
func AsInt64(v any) int64 {
	switch n := v.(type) {
	case nil:
		return 0
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return n
	case uint:
		return int64(n)
	case uint64:
		return int64(n)
	case float32:
		return int64(n)
	case float64:
		return int64(n)
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			f, ferr := n.Float64()
			if ferr != nil {
				return 0
			}
			return int64(f)
		}
		return i
	case bool:
		if n {
			return 1
		}
		return 0
	case []byte:
		return int64(Atoi(string(n)))
	case string:
		s := strings.TrimSpace(n)
		if s == "" {
			return 0
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}
		return i
	case Int:
		return int64(n.Int())
	case *Int:
		if n == nil {
			return 0
		}
		return int64(n.Int())
	default:
		return 0
	}
}

// Present reports whether v is a non-nil, non-empty optional filter value.
func Present(v any) bool {
	if v == nil {
		return false
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) != ""
	}
	return true
}

// Query looks up the first non-empty query/post-form value among keys and parses int.
func Query(c *gin.Context, keys ...string) int {
	for _, key := range keys {
		if s := strings.TrimSpace(c.Query(key)); s != "" {
			return Atoi(s)
		}
		if s := strings.TrimSpace(c.PostForm(key)); s != "" {
			return Atoi(s)
		}
	}
	return 0
}

// First returns the first non-empty trimmed string among values.
func First(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// RequestInt reads id-like params from query/form, then optional JSON object fields.
// body is typically map[string]any from a prior ShouldBindJSON; may be nil.
func RequestInt(c *gin.Context, body map[string]any, keys ...string) int {
	for _, key := range keys {
		if s := strings.TrimSpace(c.Query(key)); s != "" {
			return Atoi(s)
		}
		if s := strings.TrimSpace(c.PostForm(key)); s != "" {
			return Atoi(s)
		}
	}
	for _, key := range keys {
		if body == nil {
			break
		}
		if v, ok := body[key]; ok {
			return AsInt(v)
		}
	}
	return 0
}
