package handlers

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	t.Run("Register_Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "users"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		body := `{"email":"new@test.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.Register(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Contains(t, rec.Body.String(), "registration successful")
		}
	})

	t.Run("Register_Conflict", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "users"`).
			WillReturnError(errors.New("duplicate key error"))
		mock.ExpectRollback()

		body := `{"email":"exists@test.com","password":"123"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.NoError(t, h.Register(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("Login_Success", func(t *testing.T) {
		password := "secret123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		uID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
			AddRow(uID, "test@test.com", string(hash))
		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

		body := `{"email":"test@test.com","password":"secret123"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "token")
		}
	})

	t.Run("Register_InvalidInput", func(t *testing.T) {
		body := `{"email":"","password":""}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.NoError(t, h.Register(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
