package auth

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	dummyHash = "$2y$12$Drtj0fYeNNPEvDWy2dp5zu38eeUus1PI7HSYUWRMOoueJU1N/QlUq"
	cost      = 12
)

func PasswordOK(hash, password string) bool {
	if hash == "" {
		hash = dummyHash
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashPassword(password string) (string, error) {
	sum, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(sum), nil
}

// PasswordError is a password rejected by CheckPassword. Error() is safe to show to the user.
type PasswordError struct{ Message string }

func (e *PasswordError) Error() string {
	if e == nil || e.Message == "" {
		return "密码不符合要求"
	}
	return e.Message
}

var (
	asciiLetter = regexp.MustCompile(`[a-zA-Z]`)
	asciiDigit  = regexp.MustCompile(`\d`)
	weakExact   = map[string]struct{}{
		"abc12345": {}, "abcd1234": {}, "pass1234": {}, "password1": {},
		"qwer1234": {}, "admin123": {}, "admin888": {}, "test1234": {},
		"welcome1": {}, "iloveyou1": {}, "1qaz2wsx": {}, "1q2w3e4r": {},
		"qixue123": {}, "qixue1234": {},
	}
	weakParts = []string{"password", "qwerty", "123456", "abcdef", "111111", "000000"}
	keyRows   = []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
)

// CheckPassword rejects short, single-class, and commonly guessed passwords.
func CheckPassword(password string) error {
	n := utf8.RuneCountInString(password)
	if n < 8 || n > 64 {
		return &PasswordError{Message: "密码长度需要在8-64个字符之间"}
	}
	if !asciiLetter.MatchString(password) || !asciiDigit.MatchString(password) {
		return &PasswordError{Message: "密码必须同时包含字母和数字"}
	}
	if weakPassword(password) {
		return &PasswordError{Message: "密码过于简单，请更换更复杂的密码"}
	}
	return nil
}

func weakPassword(password string) bool {
	lower := strings.ToLower(password)
	if _, ok := weakExact[lower]; ok {
		return true
	}
	for _, part := range weakParts {
		if strings.Contains(lower, part) {
			return true
		}
	}
	if uniqueRunes(lower) <= 2 || isSequential(lower) {
		return true
	}
	letters, digits := splitClasses(lower)
	if len(letters) >= 3 && len(digits) >= 3 && (isSequential(letters) || onKeyboardRow(letters)) && isSequential(digits) {
		return true
	}
	return false
}

func uniqueRunes(s string) int {
	seen := map[rune]struct{}{}
	for _, r := range s {
		seen[r] = struct{}{}
	}
	return len(seen)
}

func isSequential(s string) bool {
	r := []rune(s)
	if len(r) < 3 {
		return false
	}
	up, down := true, true
	for i := 1; i < len(r); i++ {
		if r[i] != r[i-1]+1 {
			up = false
		}
		if r[i] != r[i-1]-1 {
			down = false
		}
	}
	return up || down
}

func onKeyboardRow(letters string) bool {
	if len(letters) < 3 {
		return false
	}
	rev := reverse(letters)
	for _, row := range keyRows {
		if strings.Contains(row, letters) || strings.Contains(row, rev) {
			return true
		}
	}
	return false
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func splitClasses(s string) (letters, digits string) {
	var lb, db strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			lb.WriteRune(r)
		case r >= '0' && r <= '9':
			db.WriteRune(r)
		}
	}
	return lb.String(), db.String()
}
