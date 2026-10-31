package auth

import "testing"

func TestPasswordOK(t *testing.T) {
	hash, err := HashPassword("secret-pass")
	if err != nil {
		t.Fatal(err)
	}
	if !PasswordOK(hash, "secret-pass") {
		t.Fatal("expected match")
	}
	if PasswordOK(hash, "other") {
		t.Fatal("expected mismatch")
	}
	if PasswordOK("", "secret-pass") {
		t.Fatal("empty hash must not match")
	}
}

func TestCheckPassword(t *testing.T) {
	ok := []string{"Qx7m2kp9", "Lesson2026", "k3m-River"}
	bad := []string{"", "abc123", "abcdefgh", "12345678", "password1", "abc12345", "qwer1234", "a1a1a1a1", "123456a"}
	for _, password := range ok {
		if err := CheckPassword(password); err != nil {
			t.Fatalf("%q rejected: %v", password, err)
		}
	}
	for _, password := range bad {
		if err := CheckPassword(password); err == nil {
			t.Fatalf("%q accepted", password)
		}
	}
}
