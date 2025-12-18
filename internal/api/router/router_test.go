package router

import (
	"testing"
	"todo-list/internal/api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// 1. Создаем мок для интерфейса TaskHandler
type mockTaskHandler struct{}

// Реализуем ВСЕ методы, которые используются в NewRouter
func (m *mockTaskHandler) Create(c echo.Context) error           { return nil }
func (m *mockTaskHandler) List(c echo.Context) error             { return nil }
func (m *mockTaskHandler) Get(c echo.Context) error              { return nil }
func (m *mockTaskHandler) Update(c echo.Context) error           { return nil }
func (m *mockTaskHandler) Delete(c echo.Context) error           { return nil }
func (m *mockTaskHandler) ChangeStatus(c echo.Context) error     { return nil }
func (m *mockTaskHandler) ChangePriority(c echo.Context) error   { return nil }
func (m *mockTaskHandler) Archive(c echo.Context) error          { return nil }
func (m *mockTaskHandler) Unarchive(c echo.Context) error        { return nil }
func (m *mockTaskHandler) ListByStatus(c echo.Context) error     { return nil }
func (m *mockTaskHandler) ListByPriority(c echo.Context) error   { return nil }
func (m *mockTaskHandler) ListByTag(c echo.Context) error        { return nil }
func (m *mockTaskHandler) Search(c echo.Context) error           { return nil }
func (m *mockTaskHandler) GetToday(c echo.Context) error         { return nil }
func (m *mockTaskHandler) GetOverdue(c echo.Context) error       { return nil }
func (m *mockTaskHandler) AddTag(c echo.Context) error           { return nil }
func (m *mockTaskHandler) RemoveTag(c echo.Context) error        { return nil }
func (m *mockTaskHandler) BulkDelete(c echo.Context) error       { return nil }
func (m *mockTaskHandler) BulkUpdateStatus(c echo.Context) error { return nil }
func (m *mockTaskHandler) Stats(c echo.Context) error            { return nil }

func TestNewRouter(t *testing.T) {
	e := echo.New()

	taskH := &mockTaskHandler{}
	authH := &handlers.AuthHandler{}
	secret := "test-secret"

	NewRouter(e, taskH, authH, secret)

	assert.Greater(t, len(e.Routes()), 0)

	found := false
	for _, r := range e.Routes() {
		if r.Path == "/auth/login" {
			found = true
			break
		}
	}
	assert.True(t, found, "Маршрут /auth/login должен быть зарегистрирован")
}
