package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
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

func (h *taskHandlerImpl) getUserID(c echo.Context) string {
	return c.Get("user_id").(string)
}

func (h *taskHandlerImpl) Create(c echo.Context) error {
	var req dto.TaskRequestDTO
	c.Bind(&req)
	var due *time.Time
	if req.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, req.DueDate); err == nil {
			due = &t
		}
	}
	task, err := h.service.CreateTask(c.Request().Context(), h.getUserID(c), req.Title, req.Content, req.Status, req.Priority, due)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) List(c echo.Context) error {
	tasks, err := h.service.GetAllTasks(c.Request().Context(), h.getUserID(c))
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
	task, err := h.service.GetTaskByID(c.Request().Context(), c.Param("id"), h.getUserID(c))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) Update(c echo.Context) error {
	var req dto.TaskRequestDTO
	c.Bind(&req)
	var due *time.Time
	if req.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, req.DueDate); err == nil {
			due = &t
		}
	}
	task, err := h.service.UpdateTask(c.Request().Context(), c.Param("id"), h.getUserID(c), req.Title, req.Content, req.Status, req.Priority, due)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) Delete(c echo.Context) error {
	err := h.service.DeleteTask(c.Request().Context(), c.Param("id"), h.getUserID(c))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *taskHandlerImpl) ChangeStatus(c echo.Context) error {
	var body struct {
		Status string `json:"status"`
	}
	c.Bind(&body)
	task, err := h.service.ChangeStatus(c.Request().Context(), c.Param("id"), h.getUserID(c), body.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ListByStatus(c echo.Context) error {
	tasks, err := h.service.GetTasksByStatus(c.Request().Context(), c.Param("status"), h.getUserID(c))
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
	tasks, err := h.service.SearchTasks(c.Request().Context(), c.QueryParam("q"), h.getUserID(c))
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
	tasks, err := h.service.GetTodayTasks(c.Request().Context(), h.getUserID(c))
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
	tasks, err := h.service.GetOverdueTasks(c.Request().Context(), h.getUserID(c))
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
	task, err := h.service.ArchiveTask(c.Request().Context(), c.Param("id"), h.getUserID(c))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) Unarchive(c echo.Context) error {
	task, err := h.service.UnarchiveTask(c.Request().Context(), c.Param("id"), h.getUserID(c))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ChangePriority(c echo.Context) error {
	var body struct {
		Priority string `json:"priority"`
	}
	c.Bind(&body)
	task, err := h.service.ChangePriority(c.Request().Context(), c.Param("id"), h.getUserID(c), body.Priority)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ListByPriority(c echo.Context) error {
	tasks, err := h.service.GetTasksByPriority(c.Request().Context(), c.Param("priority"), h.getUserID(c))
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
	var body struct {
		Tag string `json:"tag"`
	}
	c.Bind(&body)
	task, err := h.service.AddTag(c.Request().Context(), c.Param("id"), h.getUserID(c), body.Tag)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) RemoveTag(c echo.Context) error {
	task, err := h.service.RemoveTag(c.Request().Context(), c.Param("id"), h.getUserID(c), c.Param("tag"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ToTaskResponseDTO(task))
}

func (h *taskHandlerImpl) ListByTag(c echo.Context) error {
	tasks, err := h.service.GetTasksByTag(c.Request().Context(), c.Param("tag"), h.getUserID(c))
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
	c.Bind(&body)
	err := h.service.BulkDelete(c.Request().Context(), body.IDs, h.getUserID(c))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *taskHandlerImpl) BulkUpdateStatus(c echo.Context) error {
	var body struct {
		IDs    []string `json:"ids"`
		Status string   `json:"status"`
	}

	c.Bind(&body)
	err := h.service.BulkUpdateStatus(c.Request().Context(), body.IDs, body.Status, h.getUserID(c))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *taskHandlerImpl) Stats(c echo.Context) error {
	s, err := h.service.Stats(c.Request().Context(), h.getUserID(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}
