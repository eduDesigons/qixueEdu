package auth

import "golang.org/x/crypto/bcrypt"

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
