package observability

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-Id"

func RequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		reqID := c.Get(HeaderRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set(HeaderRequestID, reqID)
		c.Locals("request_id", reqID)
		return c.Next()
	}
}

func RequestID(c fiber.Ctx) string {
	v := c.Locals("request_id")
	s, _ := v.(string)
	return s
}
