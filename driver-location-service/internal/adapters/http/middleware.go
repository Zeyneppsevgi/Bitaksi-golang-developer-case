package http

import (
	"github.com/driver-location-service/internal/adapters/auth"
	"github.com/driver-location-service/pkg/response"
	"github.com/gofiber/fiber/v3"
)

func InternalAuthMiddleware(verifier *auth.APIKeyVerifier) fiber.Handler {
	return func(c fiber.Ctx) error {
		if !verifier.Verify(c.Get("X-Internal-Api-Key")) {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Failure("UNAUTHORIZED", "internal authentication failed", nil))
		}
		return c.Next()
	}
}
