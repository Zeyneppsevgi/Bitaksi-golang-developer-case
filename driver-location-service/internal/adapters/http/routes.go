package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Register(app *fiber.App, h *Handler, internalAuth fiber.Handler) {
	app.Get("/healthz", h.Healthz)
	app.Get("/readyz", h.Readyz)
	app.Get("/openapi.yaml", h.OpenAPIYAML)
	app.Get("/openapi.json", h.OpenAPIJSON)
	app.Get("/docs/*", adaptor.HTTPHandler(httpSwagger.Handler(httpSwagger.URL("/openapi.json"))))

	v1 := app.Group("/v1", internalAuth)
	v1.Post("/driver-locations/batch", h.Batch)
	v1.Post("/driver-locations/import", h.Import)
	v1.Get("/driver-locations/search", h.Search)
}
