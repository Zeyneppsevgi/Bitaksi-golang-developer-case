package http

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/usecase"
	"github.com/driver-location-service/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	upsertUC   *usecase.UpsertLocations
	importUC   *usecase.ImportCSV
	searchUC   *usecase.SearchNearest
	readyCheck func() error
	openapiYML []byte
	openapiJSN []byte
	logger     *slog.Logger
}

func NewHandler(upsertUC *usecase.UpsertLocations, importUC *usecase.ImportCSV, searchUC *usecase.SearchNearest, readyCheck func() error, openapiYML []byte, openapiJSN []byte, logger *slog.Logger) *Handler {
	return &Handler{upsertUC: upsertUC, importUC: importUC, searchUC: searchUC, readyCheck: readyCheck, openapiYML: openapiYML, openapiJSN: openapiJSN, logger: logger}
}

func (h *Handler) Healthz(c fiber.Ctx) error {
	return c.JSON(response.Success(map[string]any{"status": "ok"}))
}

func (h *Handler) Readyz(c fiber.Ctx) error {
	if err := h.readyCheck(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(response.Failure("INTERNAL", "mongo not ready", nil))
	}
	return c.JSON(response.Success(map[string]any{"status": "ready"}))
}

func (h *Handler) OpenAPIYAML(c fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "application/yaml")
	return c.Send(h.openapiYML)
}

func (h *Handler) OpenAPIJSON(c fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Send(h.openapiJSN)
}

func (h *Handler) Batch(c fiber.Ctx) error {
	var req batchRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", "invalid request body", nil))
	}
	items := make([]domain.DriverLocation, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.DriverLocation{DriverID: item.DriverID, Location: domain.Point{Type: item.Location.Type, Coordinates: item.Location.Coordinates}})
	}
	res, err := h.upsertUC.Execute(c.Context(), items)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(response.Success(map[string]int64{"upserted": res.Upserted, "updated": res.Updated}))
}

func (h *Handler) Import(c fiber.Ctx) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", "file is required", nil))
	}
	f, err := fh.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", "cannot open file", nil))
	}
	defer f.Close()

	imported, failed, err := h.importUC.Execute(c.Context(), f)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(response.Success(map[string]int{"imported": imported, "failed": failed}))
}

func (h *Handler) Search(c fiber.Ctx) error {
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
	items, err := h.searchUC.Execute(c.Context(), lon, lat, radius, 100)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(response.Success(map[string]any{
		"center":  map[string]float64{"lon": lon, "lat": lat},
		"radiusM": radius,
		"items":   items,
	}))
}

func mapError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return c.Status(fiber.StatusBadRequest).JSON(response.Failure("VALIDATION_ERROR", err.Error(), nil))
	case errors.Is(err, domain.ErrUnauthorized):
		return c.Status(fiber.StatusUnauthorized).JSON(response.Failure("UNAUTHORIZED", "unauthorized", nil))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(response.Failure("INTERNAL", fmt.Sprintf("%v", err), nil))
	}
}
