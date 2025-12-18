package router

import (
	"github.com/labstack/echo/v4"
	"todo-list/internal/api/handlers"
	"todo-list/internal/api/middleware"
)

func NewRouter(e *echo.Echo, h handlers.TaskHandler, ah *handlers.AuthHandler, secret string) {
	// Открытые маршруты
	e.POST("/auth/register", ah.Register)
	e.POST("/auth/login", ah.Login)

	// Защищенные маршруты (только с JWT)
	api := e.Group("/api/v1/tasks")
	api.Use(middleware.AuthMiddleware(secret))

	api.POST("", h.Create)
	api.GET("", h.List)
	api.GET("/:id", h.Get)
	api.PUT("/:id", h.Update)
	api.DELETE("/:id", h.Delete)

	api.PATCH("/:id/status", h.ChangeStatus)
	api.PATCH("/:id/priority", h.ChangePriority)
	api.PATCH("/:id/archive", h.Archive)
	api.PATCH("/:id/unarchive", h.Unarchive)

	api.GET("/status/:status", h.ListByStatus)
	api.GET("/priority/:priority", h.ListByPriority)
	api.GET("/tag/:tag", h.ListByTag)
	api.GET("/search", h.Search)
	api.GET("/today", h.GetToday)
	api.GET("/overdue", h.GetOverdue)

	api.POST("/:id/tags", h.AddTag)
	api.DELETE("/:id/tags/:tag", h.RemoveTag)

	api.POST("/bulk-delete", h.BulkDelete)
	api.POST("/bulk-status", h.BulkUpdateStatus)
	api.GET("/stats", h.Stats)
}
