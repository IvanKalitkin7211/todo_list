package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"todo-list/internal/api/dto"
	"todo-list/internal/domain/service"
)

type TaskHandler interface {
	Create(c echo.Context) error
	List(c echo.Context) error
	Get(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error

	ChangeStatus(c echo.Context) error
	ListByStatus(c echo.Context) error
	Search(c echo.Context) error
	GetToday(c echo.Context) error
	GetOverdue(c echo.Context) error
	Archive(c echo.Context) error
	Unarchive(c echo.Context) error
	ChangePriority(c echo.Context) error
	ListByPriority(c echo.Context) error
	AddTag(c echo.Context) error
	RemoveTag(c echo.Context) error
	ListByTag(c echo.Context) error
	BulkDelete(c echo.Context) error
	BulkUpdateStatus(c echo.Context) error
	Stats(c echo.Context) error
}

type taskHandlerImpl struct {
	service service.TaskService
}

func NewTaskHandler(s service.TaskService) TaskHandler {
	return &taskHandlerImpl{service: s}
}

func (h *taskHandlerImpl) Create(c echo.Context) error {
	var req dto.TaskRequestDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	var due *time.Time
	if req.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, req.DueDate); err == nil {
			due = &t
		}
	}
	ctx := c.Request().Context()
	task, err := h.service.CreateTask(ctx, req.Title, req.Content, req.Status, req.Priority, due)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) List(c echo.Context) error {
	ctx := c.Request().Context()
	tasks, err := h.service.GetAllTasks(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]dto.TaskResponseDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.ToTaskResponseDTO(t))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *taskHandlerImpl) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id"})
	}
	ctx := c.Request().Context()
	task, err := h.service.GetTaskByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id"})
	}
	var req dto.TaskRequestDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	var due *time.Time
	if req.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, req.DueDate); err == nil {
			due = &t
		}
	}
	ctx := c.Request().Context()
	task, err := h.service.UpdateTask(ctx, id, req.Title, req.Content, req.Status, req.Priority, due)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing id"})
	}
	ctx := c.Request().Context()
	if err := h.service.DeleteTask(ctx, id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *taskHandlerImpl) ChangeStatus(c echo.Context) error {
	id := c.Param("id")
	var body struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	ctx := c.Request().Context()
	task, err := h.service.ChangeStatus(ctx, id, body.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ListByStatus(c echo.Context) error {
	status := c.Param("status")
	ctx := c.Request().Context()
	tasks, err := h.service.GetTasksByStatus(ctx, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]dto.TaskResponseDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.ToTaskResponseDTO(t))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *taskHandlerImpl) Search(c echo.Context) error {
	q := c.QueryParam("q")
	ctx := c.Request().Context()
	tasks, err := h.service.SearchTasks(ctx, q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]dto.TaskResponseDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.ToTaskResponseDTO(t))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *taskHandlerImpl) GetToday(c echo.Context) error {
	ctx := c.Request().Context()
	tasks, err := h.service.GetTodayTasks(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]dto.TaskResponseDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.ToTaskResponseDTO(t))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *taskHandlerImpl) GetOverdue(c echo.Context) error {
	ctx := c.Request().Context()
	tasks, err := h.service.GetOverdueTasks(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]dto.TaskResponseDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.ToTaskResponseDTO(t))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *taskHandlerImpl) Archive(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	task, err := h.service.ArchiveTask(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) Unarchive(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	task, err := h.service.UnarchiveTask(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ChangePriority(c echo.Context) error {
	id := c.Param("id")
	var body struct {
		Priority string `json:"priority"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	ctx := c.Request().Context()
	task, err := h.service.ChangePriority(ctx, id, body.Priority)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ListByPriority(c echo.Context) error {
	p := c.Param("priority")
	ctx := c.Request().Context()
	tasks, err := h.service.GetTasksByPriority(ctx, p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]dto.TaskResponseDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.ToTaskResponseDTO(t))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *taskHandlerImpl) AddTag(c echo.Context) error {
	id := c.Param("id")
	var body struct {
		Tag string `json:"tag"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	tag := strings.TrimSpace(body.Tag)
	if tag == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tag empty"})
	}
	ctx := c.Request().Context()
	task, err := h.service.AddTag(ctx, id, tag)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) RemoveTag(c echo.Context) error {
	id := c.Param("id")
	tag := c.Param("tag")
	if tag == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing tag"})
	}
	ctx := c.Request().Context()
	task, err := h.service.RemoveTag(ctx, id, tag)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ListByTag(c echo.Context) error {
	tag := c.Param("tag")
	ctx := c.Request().Context()
	tasks, err := h.service.GetTasksByTag(ctx, tag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]dto.TaskResponseDTO, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.ToTaskResponseDTO(t))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *taskHandlerImpl) BulkDelete(c echo.Context) error {
	var body struct {
		IDs []string `json:"ids"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	ctx := c.Request().Context()
	if err := h.service.BulkDelete(ctx, body.IDs); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *taskHandlerImpl) BulkUpdateStatus(c echo.Context) error {
	var body struct {
		IDs    []string `json:"ids"`
		Status string   `json:"status"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	ctx := c.Request().Context()
	if err := h.service.BulkUpdateStatus(ctx, body.IDs, body.Status); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *taskHandlerImpl) Stats(c echo.Context) error {
	ctx := c.Request().Context()
	stats, err := h.service.Stats(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, stats)
}
