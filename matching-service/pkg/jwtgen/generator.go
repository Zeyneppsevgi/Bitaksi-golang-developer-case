package jwtgen

import (
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func Generator() (string, error) {

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtv5.MapClaims{
		"authenticated": true,
		"exp":           time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte("user-secret"))
	if err != nil {
		return "", err
	}

	return signed, nil
}
