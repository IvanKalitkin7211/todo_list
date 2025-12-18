package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAuthHandler_Unit(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{Conn: dbMock})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	e := echo.New()
	h := &AuthHandler{DB: db, Secret: "test"}

	t.Run("Login_UserNotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"no@test.com","password":"123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
