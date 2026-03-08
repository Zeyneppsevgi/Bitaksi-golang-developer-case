package auth

import (
	"testing"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func TestUserJWTAuthenticator(t *testing.T) {
	secret := "user-secret"
	a := NewUserJWTAuthenticator(secret)
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtv5.MapClaims{"authenticated": true})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	ok, err := a.IsAuthenticated(nil, "Bearer "+signed)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected authenticated=true")
	}
}
