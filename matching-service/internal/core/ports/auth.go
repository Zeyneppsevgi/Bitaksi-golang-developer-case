package ports

import "context"

type UserAuthenticator interface {
	IsAuthenticated(ctx context.Context, authHeader string) (bool, error)
}
