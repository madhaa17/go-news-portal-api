package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"news-app/config"
	"news-app/internal/adapter/handler"
	"news-app/internal/adapter/repository"
	"news-app/internal/core/service"
	"news-app/lib/auth"
	"news-app/lib/middleware"
	"news-app/lib/pagination"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnPostgres()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to connect to database: %v", err)
		return
	}

	// cloudflareR2
	cdfR2 := cfg.LoadAwsConfig()
	_ = s3.NewFromConfig(cdfR2)
	_ = auth.NewJwt(cfg)
	_ = middleware.NewMiddleware(cfg)
	_ = pagination.NewPagination()

	// repository
	authRepo := repository.NewAuthRepository(db.DB)

	// service
	authService := service.NewAuthService(authRepo, cfg, auth.NewJwt(cfg))

	// handler
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(
		logger.Config{
			Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
		},
	))

	api := app.Group("/api")
	api.Post("/auth/login", authHandler.Login)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Logger.Println("Server shuttdown of 5s")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
