package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"todo-list/config"
	"todo-list/internal/api/handlers"
	"todo-list/internal/api/router"
	"todo-list/internal/domain/model"
	"todo-list/internal/domain/service"
	"todo-list/internal/infrastructure/database/postgres"
	"todo-list/internal/infrastructure/repository"
)

func Start() {
	cfg := config.NewConfig()

	// 1. Инициализация БД
	dbConn, err := postgres.ProvideDBClient(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// 2. Автомиграция (создает таблицы пользователей и задач)
	db := dbConn.GetDB()
	if err := db.AutoMigrate(&model.User{}, &model.Task{}, &model.Tag{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Сборка слоев (Dependency Injection)
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)

	// 4. Настройка Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 5. Инициализация роутера
	router.NewRouter(e, taskHandler, authHandler, cfg.JWTSecret)

	// 6. Запуск сервера
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
