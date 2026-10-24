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
