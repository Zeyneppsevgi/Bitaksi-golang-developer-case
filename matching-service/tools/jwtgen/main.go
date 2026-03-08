package main

import (
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func main() {
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtv5.MapClaims{
		"authenticated": true,
		"exp":           time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte("user-secret"))
	if err != nil {
		panic(err)
	}
	fmt.Println(signed)
}
