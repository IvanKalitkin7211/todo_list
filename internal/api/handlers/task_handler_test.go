package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-list/internal/domain/model"
	"todo-list/internal/testutils"
)

func TestHandler_Complete(t *testing.T) {
	e := echo.New()
	mockSvc := new(testutils.AllMocks)
	h := NewTaskHandler(mockSvc)
	uID := "test-user"

	t.Run("Create_Success", func(t *testing.T) {
		body := `{"title":"API","content":"Desc"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uID)

		mockSvc.On("CreateTask", mock.Anything, uID, "API", "Desc", "", "", mock.Anything).
			Return(model.Task{Title: "API"}, nil).Once()

		if assert.NoError(t, h.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("List_Tasks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uID)

		mockSvc.On("GetAllTasks", mock.Anything, uID).Return([]model.Task{{Title: "T1"}}, nil).Once()

		if assert.NoError(t, h.List(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "T1")
		}
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")
		c.Set("user_id", uID)

		mockSvc.On("DeleteTask", mock.Anything, "999", uID).Return(errors.New("not found")).Once()

		err := h.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Search_Tasks", func(t *testing.T) {
		// Эмулируем запрос /api/v1/tasks/search?q=milk
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/search?q=milk", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uID)

		mockSvc.On("SearchTasks", mock.Anything, "milk", uID).
			Return([]model.Task{{Title: "Buy milk"}}, nil).Once()

		if assert.NoError(t, h.Search(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Buy milk")
		}
	})

	t.Run("Stats_Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uID)

		mockSvc.On("Stats", mock.Anything, uID).
			Return(map[string]int64{"todo": 5, "done": 2}, nil).Once()

		if assert.NoError(t, h.Stats(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), `"todo":5`)
		}
	})

	t.Run("AddTag_Success", func(t *testing.T) {
		body := `{"tag":"urgent"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/tasks/:id/tags")
		c.SetParamNames("id")
		c.SetParamValues("123")
		c.Set("user_id", uID)

		mockSvc.On("AddTag", mock.Anything, "123", uID, "urgent").
			Return(model.Task{ID: [16]byte{1}}, nil).Once()

		if assert.NoError(t, h.AddTag(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("BulkUpdateStatus_Success", func(t *testing.T) {
		body := `{"ids":["1","2"],"status":"done"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", uID)

		mockSvc.On("BulkUpdateStatus", mock.Anything, []string{"1", "2"}, "done", uID).
			Return(nil).Once()

		if assert.NoError(t, h.BulkUpdateStatus(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Priority_Update_Handler", func(t *testing.T) {
		body := `{"priority":"high"}`
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/tasks/:id/priority")
		c.SetParamNames("id")
		c.SetParamValues("123")
		c.Set("user_id", uID)

		mockSvc.On("ChangePriority", mock.Anything, "123", uID, "high").
			Return(model.Task{Priority: "high"}, nil).Once()

		if assert.NoError(t, h.ChangePriority(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Archive_Handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/tasks/:id/archive")
		c.SetParamNames("id")
		c.SetParamValues("123")
		c.Set("user_id", uID)

		mockSvc.On("ArchiveTask", mock.Anything, "123", uID).
			Return(model.Task{Status: "archived"}, nil).Once()

		if assert.NoError(t, h.Archive(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("List_By_Filters_Handlers", func(t *testing.T) {
		// 1. ListByStatus
		reqStatus := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/status/todo", nil)
		recStatus := httptest.NewRecorder()
		cStatus := e.NewContext(reqStatus, recStatus)
		cStatus.SetPath("/api/v1/tasks/status/:status")
		cStatus.SetParamNames("status")
		cStatus.SetParamValues("todo")
		cStatus.Set("user_id", uID)

		mockSvc.On("GetTasksByStatus", mock.Anything, "todo", uID).Return([]model.Task{{Title: "S"}}, nil).Once()
		assert.NoError(t, h.ListByStatus(cStatus))

		// 2. ListByTag
		reqTag := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/tag/work", nil)
		recTag := httptest.NewRecorder()
		cTag := e.NewContext(reqTag, recTag)
		cTag.SetPath("/api/v1/tasks/tag/:tag")
		cTag.SetParamNames("tag")
		cTag.SetParamValues("work")
		cTag.Set("user_id", uID)

		mockSvc.On("GetTasksByTag", mock.Anything, "work", uID).Return([]model.Task{{Title: "T"}}, nil).Once()
		assert.NoError(t, h.ListByTag(cTag))
	})
}
