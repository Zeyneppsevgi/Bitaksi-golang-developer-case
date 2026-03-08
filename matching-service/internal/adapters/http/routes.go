package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Register(app *fiber.App, h *Handler, auth fiber.Handler) {
	app.Get("/healthz", h.Healthz)
	app.Get("/openapi.yaml", h.OpenAPIYAML)
	app.Get("/openapi.json", h.OpenAPIJSON)
	app.Get("/docs/*", adaptor.HTTPHandler(httpSwagger.Handler(httpSwagger.URL("/openapi.json"))))

	v1 := app.Group("/v1", auth)
	v1.Get("/match", h.Match)
	v1.Get("/token", h.GenerateToken)

}
