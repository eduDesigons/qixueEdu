package textutil

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var stripNoise = regexp.MustCompile(`[\s\p{P}\p{S}]+`)

// AverageSimilarity averages cosine, Jaccard, and Levenshtein similarity.
func AverageSimilarity(standard, user any) float64 {
	return (CosineSimilarity(standard, user) + JaccardSimilarity(standard, user) + LevenshteinSimilarity(standard, user)) / 3
}

func CosineSimilarity(standard, user any) float64 {
	sv := countValues(tokens(standard))
	uv := countValues(tokens(user))
	if len(sv) == 0 || len(uv) == 0 {
		return 0
	}
	intersection := 0.0
	for w, c := range sv {
		if uc, ok := uv[w]; ok {
			intersection += float64(c * uc)
		}
	}
	ns := norm(sv)
	nu := norm(uv)
	if ns == 0 || nu == 0 {
		return 0
	}
	return (intersection / (ns * nu)) * 100
}

func JaccardSimilarity(standard, user any) float64 {
	sw := unique(tokens(standard))
	uw := unique(tokens(user))
	unionSet := map[string]struct{}{}
	for _, w := range sw {
		unionSet[w] = struct{}{}
	}
	for _, w := range uw {
		unionSet[w] = struct{}{}
	}
	if len(unionSet) == 0 {
		return 0
	}
	uwSet := map[string]struct{}{}
	for _, w := range uw {
		uwSet[w] = struct{}{}
	}
	inter := 0
	for _, w := range sw {
		if _, ok := uwSet[w]; ok {
			inter++
		}
	}
	return (float64(inter) / float64(len(unionSet))) * 100
}

func LevenshteinSimilarity(standard, user any) float64 {
	a := plain(standard)
	b := plain(user)
	maxLen := len([]rune(a))
	if lb := len([]rune(b)); lb > maxLen {
		maxLen = lb
	}
	if maxLen == 0 {
		return 0
	}
	dist := charDistance(a, b)
	return (1 - float64(dist)/float64(maxLen)) * 100
}

func plain(text any) string {
	s := flatten(text)
	s = strings.ToLower(s)
	return stripNoise.ReplaceAllString(s, "")
}

func flatten(text any) string {
	switch v := text.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		if v {
			return "1"
		}
		return ""
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, flatten(item))
		}
		return strings.Join(parts, "")
	case []string:
		return strings.Join(v, "")
	case map[string]any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, flatten(item))
		}
		return strings.Join(parts, "")
	default:
		return ""
	}
}

func tokens(text any) []string {
	chars := []rune(plain(text))
	n := len(chars)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []string{string(chars[0])}
	}
	grams := make([]string, 0, n-1)
	for i := 0; i < n-1; i++ {
		grams = append(grams, string(chars[i])+string(chars[i+1]))
	}
	return grams
}

func countValues(items []string) map[string]int {
	out := map[string]int{}
	for _, item := range items {
		out[item]++
	}
	return out
}

func unique(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func norm(m map[string]int) float64 {
	sum := 0.0
	for _, c := range m {
		sum += float64(c * c)
	}
	return math.Sqrt(sum)
}

func charDistance(left, right string) int {
	a := []rune(left)
	b := []rune(right)
	n, m := len(a), len(b)
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}
	if n < m {
		a, b, n, m = b, a, m, n
	}
	prev := make([]int, m+1)
	for j := 0; j <= m; j++ {
		prev[j] = j
	}
	for i := 1; i <= n; i++ {
		curr := make([]int, m+1)
		curr[0] = i
		for j := 1; j <= m; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
		}
		prev = curr
	}
	return prev[m]
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}
