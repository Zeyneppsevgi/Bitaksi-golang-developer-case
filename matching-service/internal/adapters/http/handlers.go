package http

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/matching-service/internal/core/domain"
	"github.com/matching-service/internal/core/usecase"
	"github.com/matching-service/pkg/jwtgen"
	"github.com/matching-service/pkg/response"
)

type Handler struct {
	findUC     *usecase.FindDriver
	openapiYML []byte
	openapiJSN []byte
}

func NewHandler(findUC *usecase.FindDriver, openapiYML []byte, openapiJSN []byte) *Handler {
	return &Handler{findUC: findUC, openapiYML: openapiYML, openapiJSN: openapiJSN}
}

func (h *Handler) Healthz(c fiber.Ctx) error {
	return c.JSON(response.Success(map[string]any{"status": "ok"}))
}

func (h *Handler) OpenAPIYAML(c fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "application/yaml")
	return c.Send(h.openapiYML)
}

func (h *Handler) OpenAPIJSON(c fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Send(h.openapiJSN)
}

func (h *Handler) Match(c fiber.Ctx) error {
	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", "invalid lon", nil))
	}
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", "invalid lat", nil))
	}
	radius, err := strconv.ParseInt(c.Query("radius_m", "3000"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", "invalid radius_m", nil))
	}
	match, err := h.findUC.Execute(c.Context(), lon, lat, radius, c.Get("X-Request-Id"))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", err.Error(), nil))
		case errors.Is(err, domain.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(response.Failure("NOT_FOUND", "no driver found", nil))
		case errors.Is(err, domain.ErrUpstreamUnavailable):
			return c.Status(fiber.StatusServiceUnavailable).JSON(response.Failure("UPSTREAM_UNAVAILABLE", "driver-location unavailable", nil))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(response.Failure("INTERNAL", fmt.Sprintf("%v", err), nil))
		}
	}
	return c.JSON(response.Success(match))
}

func (h *Handler) GenerateToken(c fiber.Ctx) error {

	token, err := jwtgen.Generator()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(response.Failure("INTERNAL", err.Error(), nil))
	}

	return c.JSON(response.Success(map[string]string{
		"token": token,
	}))
}
