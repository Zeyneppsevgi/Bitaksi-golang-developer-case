package auth

import (
	"context"
	"fmt"
	"strings"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type UserJWTAuthenticator struct {
	secret string
}

func NewUserJWTAuthenticator(secret string) *UserJWTAuthenticator {
	return &UserJWTAuthenticator{secret: secret}
}

func (a *UserJWTAuthenticator) IsAuthenticated(_ context.Context, authHeader string) (bool, error) {
	parts := strings.SplitN(strings.TrimSpace(authHeader), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return false, fmt.Errorf("invalid authorization header")
	}
	token := parts[1]
	parsed, err := jwtv5.Parse(token, func(t *jwtv5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(a.secret), nil
	})
	if err != nil || !parsed.Valid {
		return false, err
	}
	claims, ok := parsed.Claims.(jwtv5.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid claims")
	}
	flag, ok := claims["authenticated"].(bool)
	if !ok {
		return false, nil
	}
	return flag, nil
}
