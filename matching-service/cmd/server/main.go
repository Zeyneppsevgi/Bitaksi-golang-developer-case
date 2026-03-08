package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/matching-service/internal/adapters/auth"
	"github.com/matching-service/internal/adapters/config"
	driverlocationadapter "github.com/matching-service/internal/adapters/driverlocation"
	httpadapter "github.com/matching-service/internal/adapters/http"
	"github.com/matching-service/internal/adapters/observability"
	"github.com/matching-service/internal/core/usecase"
)

func main() {
	cfg := config.Load()
	logger := observability.NewLogger("matching-service")

	client := driverlocationadapter.NewClient(cfg.DriverLocationBaseURL, cfg.InternalAPIKey)
	findUC := usecase.NewFindDriver(client)
	yml, err := os.ReadFile("openapi/openapi.yaml")
	if err != nil {
		panic(err)
	}
	jsn, err := os.ReadFile("openapi/openapi.json")
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{AppName: "matching-service"})
	app.Use(observability.RequestIDMiddleware())
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		logger.Info("request",
			slog.String("request_id", observability.RequestID(c)),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return err
	})

	authenticator := auth.NewUserJWTAuthenticator(cfg.UserJWTSecret)
	h := httpadapter.NewHandler(findUC, yml, jsn)
	httpadapter.Register(app, h, httpadapter.UserAuthMiddleware(authenticator))

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			logger.Error("listen failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	_ = app.Shutdown()
}
