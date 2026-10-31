package paper

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"qixue-api/internal/model"
)

var chineseRe = regexp.MustCompile(`[\x{4e00}-\x{9fa5}]`)

func completionOptions(answer any, title string) ([]string, []int) {
	answers := toStringSlice(answer)
	keyWords := make([]string, 0)
	answerCount := make([]int, 0, len(answers))
	for _, item := range answers {
		answerCount = append(answerCount, utf8.RuneCountInString(item))
		for _, r := range item {
			keyWords = append(keyWords, string(r))
		}
	}
	if len(keyWords) < 15 {
		chars := chineseRe.FindAllString(title, -1)
		seen := map[string]struct{}{}
		unique := make([]string, 0, len(chars))
		for _, ch := range chars {
			if _, ok := seen[ch]; ok {
				continue
			}
			seen[ch] = struct{}{}
			unique = append(unique, ch)
		}
		rand.Shuffle(len(unique), func(i, j int) { unique[i], unique[j] = unique[j], unique[i] })
		for len(keyWords) < 15 && len(unique) > 0 {
			keyWords = append(keyWords, unique[0])
			unique = unique[1:]
		}
	}
	rand.Shuffle(len(keyWords), func(i, j int) { keyWords[i], keyWords[j] = keyWords[j], keyWords[i] })
	return keyWords, answerCount
}

func answersEqual(standard, user any) bool {
	a := toComparableSlice(standard)
	b := toComparableSlice(user)
	if len(a) != len(b) {
		return false
	}
	counts := map[string]int{}
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	for _, c := range counts {
		if c != 0 {
			return false
		}
	}
	return true
}

func toComparableSlice(v any) []string {
	switch n := v.(type) {
	case nil:
		return nil
	case string:
		if n == "" {
			return nil
		}
		return []string{n}
	case float64:
		return []string{strconv.FormatFloat(n, 'f', -1, 64)}
	case int:
		return []string{strconv.Itoa(n)}
	case []any:
		out := make([]string, 0, len(n))
		for _, item := range n {
			out = append(out, toComparableSlice(item)...)
		}
		return out
	case []string:
		return n
	default:
		return nil
	}
}

func toStringSlice(v any) []string {
	switch n := v.(type) {
	case []any:
		out := make([]string, 0, len(n))
		for _, item := range n {
			out = append(out, flattenScalar(item))
		}
		return out
	case []string:
		return n
	case string:
		if n == "" {
			return []string{}
		}
		return []string{n}
	default:
		return []string{}
	}
}

func flattenScalar(v any) string {
	switch n := v.(type) {
	case string:
		return n
	case float64:
		return strconv.FormatFloat(n, 'f', -1, 64)
	case int:
		return strconv.Itoa(n)
	default:
		return ""
	}
}

func isEmptyAnswer(v any) bool {
	if v == nil {
		return true
	}
	switch n := v.(type) {
	case string:
		return strings.TrimSpace(n) == ""
	case []any:
		return len(n) == 0
	case []string:
		return len(n) == 0
	default:
		return false
	}
}

func decodeJSON(raw json.RawMessage) any {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return string(raw)
	}
	return value
}

func formatLogTimes(row map[string]any) {
	for _, key := range []string{"created_at", "updated_at"} {
		switch v := row[key].(type) {
		case time.Time:
			row[key] = v.Format("2006-01-02 15:04:05")
		case *time.Time:
			row[key] = model.FormatTime(v)
		}
	}
}
