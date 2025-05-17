package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"product-service/app/handler"
	"product-service/app/middleware"
	"product-service/app/repository/db"
	"product-service/app/repository/warehouse"
	"product-service/app/usecase"
	"product-service/config"
	"product-service/pkg/logger"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
)

func main() {
	// init logger
	logger.InitLogger()

	// init config
	cfg, err := config.InitConfig(context.Background())
	if err != nil {
		slog.Error("failed to init config", "error", err)
		return
	}

	// init database
	dbConn, err := db.NewPostgres(cfg.Db)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer dbConn.Close()

	// init redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Error("failed to close redis client", "error", err)
		}
	}()

	reqValidator := validator.New()
	productReadRepo := db.NewProductReadRepository(dbConn)
	warehouseRepo := warehouse.NewWarehouseRepository(redisClient, cfg.WarehouseService.Host, time.Duration(0))

	productReadUsecase := usecase.NewUserUsecase(productReadRepo, warehouseRepo, cfg)

	productReadHandler := handler.NewProductReadHandler(productReadUsecase, reqValidator)

	// Initialize HTTP web framework
	app := fiber.New()
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(middleware.RequestIDMiddleware())

	handler.SetupRouter(app, productReadHandler)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			slog.Error("Failed to listen", "port", cfg.Port)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Gracefully shutdown")
	err = app.Shutdown()
	if err != nil {
		slog.Warn("Unfortunately the shutdown wasn't smooth", "err", err)
	}
}
