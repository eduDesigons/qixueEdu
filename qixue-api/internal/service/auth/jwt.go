package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrExpired = errors.New("token expired")

type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

type Issuer struct {
	key    []byte
	expire time.Duration
}

func NewIssuer(key string, expire time.Duration) Issuer {
	return Issuer{key: []byte(key), expire: expire}
}

func (i Issuer) Sign(id int, now time.Time) (string, int64, error) {
	exp := now.Add(i.expire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	signed, err := token.SignedString(i.key)
	if err != nil {
		return "", 0, err
	}
	return signed, exp.Unix(), nil
}

func (i Issuer) Parse(raw string) (id int, exp int64, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return i.key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, 0, ErrExpired
		}
		return 0, 0, err
	}
	if !token.Valid || claims.ID <= 0 {
		return 0, 0, errors.New("invalid token")
	}
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Unix()
	}
	return claims.ID, exp, nil
}
