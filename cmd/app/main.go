package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tender-bridge/cmd/app/server"
	"tender-bridge/config"
	"tender-bridge/internal/cache"
	"tender-bridge/internal/handler"
	"tender-bridge/internal/repository"
	"tender-bridge/internal/service"
	"tender-bridge/pkg/logger"
	"tender-bridge/pkg/setup"

	"github.com/go-redis/redis/v8"
)

// @title Tender Management System API
// @version 1.0
// @description API Server for Application
// @host localhost:8080
// @BasePath
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logger.GetLogger()

	db, err := setup.SetupPostgresConnection(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	redisCache := cache.NewRedisCache(redisClient)

	repos := repository.NewRepository(db, logger)
	services := service.NewService(repos, redisCache, cfg, logger)
	handlers := handler.NewHandler(services, logger)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.HTTPHost, cfg.HTTPPort, handlers.InitRoutes(cfg)); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logger.Info("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Warn("App shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Errorf("error occured on db connection close: %s", err.Error())
	}
}
