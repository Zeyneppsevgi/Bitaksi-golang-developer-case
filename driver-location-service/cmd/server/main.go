package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/driver-location-service/internal/adapters/auth"
	"github.com/driver-location-service/internal/adapters/config"
	httpadapter "github.com/driver-location-service/internal/adapters/http"
	mongoadapter "github.com/driver-location-service/internal/adapters/mongo"
	"github.com/driver-location-service/internal/adapters/observability"
	"github.com/driver-location-service/internal/core/usecase"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()
	logger := observability.NewLogger("driver-location-service")

	mongoClient, err := mongoadapter.NewClient(context.Background(), cfg.MongoURI)
	if err != nil {
		logger.Error("mongo connect failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	repo := mongoadapter.NewRepository(mongoClient, cfg.MongoDB, cfg.MongoCollection)
	if err := repo.CreateIndexes(context.Background()); err != nil {
		logger.Error("index creation failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	upsertUC := usecase.NewUpsertLocations(repo, time.Now)
	importUC := usecase.NewImportCSV(repo, 1000, time.Now)
	searchUC := usecase.NewSearchNearest(repo)

	if cfg.SeedOnStart {
		f, openErr := os.Open(cfg.SeedFile)
		if openErr != nil {
			logger.Error("seed open failed", slog.String("error", openErr.Error()), slog.String("seed_file", cfg.SeedFile))
		} else {
			imported, failed, seedErr := importUC.Execute(context.Background(), f)
			_ = f.Close()
			if seedErr != nil {
				logger.Error("seed failed", slog.String("error", seedErr.Error()))
			} else {
				logger.Info("seed completed", slog.Int("imported", imported), slog.Int("failed", failed))
			}
		}
	}

	yml, err := os.ReadFile("openapi/openapi.yaml")
	if err != nil {
		panic(err)
	}
	jsn, err := os.ReadFile("openapi/openapi.json")
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{AppName: "driver-location-service"})
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

	h := httpadapter.NewHandler(
		upsertUC,
		importUC,
		searchUC,
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return repo.Ping(ctx)
		},
		yml,
		jsn,
		logger,
	)
	httpadapter.Register(app, h, httpadapter.InternalAuthMiddleware(auth.NewAPIKeyVerifier(cfg.InternalAPIKey)))

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			logger.Error("listen failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(ctx)
	_ = mongoClient.Disconnect(ctx)
}
