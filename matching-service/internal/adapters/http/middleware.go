package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/matching-service/internal/core/ports"
	"github.com/matching-service/pkg/response"
)

func UserAuthMiddleware(auth ports.UserAuthenticator) fiber.Handler {
	return func(c fiber.Ctx) error {
		ok, err := auth.IsAuthenticated(c.Context(), c.Get(fiber.HeaderAuthorization))
		if err != nil || !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Failure("UNAUTHORIZED", "user authentication failed", nil))
		}
		return c.Next()
	}
}
