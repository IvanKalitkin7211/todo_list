package router

import (
	"github.com/labstack/echo/v4"
	"todo-list/internal/api/handlers"
)

func NewRouter(e *echo.Echo, h handlers.TaskHandler) {
	api := e.Group("/api/task")

	// CRUD
	api.GET("/", h.List)               // GET /api/task/
	api.GET("/list", h.List)           // GET /api/task/list
	api.GET("/:id", h.Get)             // GET /api/task/:id
	api.POST("/create", h.Create)      // POST /api/task/create
	api.PUT("/update/:id", h.Update)   // PUT update
	api.PATCH("/update/:id", h.Update) // PATCH update
	api.DELETE("/delete/:id", h.Delete)

	// status / priority
	api.PATCH("/:id/status", h.ChangeStatus)
	api.GET("/status/:status", h.ListByStatus)
	api.PATCH("/:id/priority", h.ChangePriority)
	api.GET("/priority/:priority", h.ListByPriority)

	// tags
	api.POST("/:id/tag", h.AddTag)
	api.DELETE("/:id/tag/:tag", h.RemoveTag)
	api.GET("/tag/:tag", h.ListByTag)

	// archive
	api.PATCH("/:id/archive", h.Archive)
	api.PATCH("/:id/unarchive", h.Unarchive)
	api.GET("/archived", h.List) // client can filter archived if needed (service supports)

	// search / filters
	api.GET("/search", h.Search)  // /api/task/search?q=...
	api.GET("/today", h.GetToday) // tasks due today
	api.GET("/overdue", h.GetOverdue)

	// bulk operations
	api.DELETE("/bulk", h.BulkDelete)             // body: { ids: [] }
	api.PATCH("/bulk/status", h.BulkUpdateStatus) // body: { ids: [], status: "done" }

	// stats
	api.GET("/stats", h.Stats)
}
