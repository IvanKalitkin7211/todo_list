package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"todo-list/config"
	"todo-list/internal/api/handlers"
	md "todo-list/internal/api/middleware"
	"todo-list/internal/api/router"
	"todo-list/internal/domain/service"
	"todo-list/internal/infrastructure/cache/redis"
	"todo-list/internal/infrastructure/database/postgres"
	"todo-list/internal/infrastructure/repository"
)

func Start() {
	cfg := config.NewConfig()

	dbConn, err := postgres.ProvideDBClient(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	redisClient, err := redis.ProvideRedisClient(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		defer redisClient.Close()
	}

	db := dbConn.GetDB()

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if redisClient != nil {
		e.Use(md.RateLimiterMiddleware(redisClient, &cfg.RateLimiter))
	}

	router.NewRouter(e, taskHandler, authHandler, cfg.JWTSecret)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
