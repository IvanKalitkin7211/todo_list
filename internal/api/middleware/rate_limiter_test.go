package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"todo-list/config"
)

func TestRateLimiterMiddleware(t *testing.T) {
	e := echo.New()

	// Настройки для тестов
	cfg := &config.RateLimiterConfig{
		Enabled:      true,
		Limit:        2,
		Window:       time.Minute,
		WindowSec:    60,
		ErrorMessage: "Too many requests",
	}

	t.Run("Disabled_Config", func(t *testing.T) {
		disabledCfg := &config.RateLimiterConfig{Enabled: false}
		mw := RateLimiterMiddleware(nil, disabledCfg)
		handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Allow_Request", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		mw := RateLimiterMiddleware(db, cfg)
		handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		redisKey := "rate_limit_192.168.1.1"
		mock.ExpectTxPipeline()
		mock.ExpectIncr(redisKey).SetVal(1)
		mock.ExpectExpire(redisKey, cfg.Window).SetVal(true)
		mock.ExpectTxPipelineExec()

		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "2", rec.Header().Get("X-RateLimit-Limit"))
		assert.Equal(t, "1", rec.Header().Get("X-RateLimit-Remaining"))
	})

	t.Run("Limit_Exceeded", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		mw := RateLimiterMiddleware(db, cfg)
		handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.2:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		redisKey := "rate_limit_192.168.1.2"
		mock.ExpectTxPipeline()
		mock.ExpectIncr(redisKey).SetVal(3) // Превышаем лимит (3 > 2)
		mock.ExpectExpire(redisKey, cfg.Window).SetVal(true)
		mock.ExpectTxPipelineExec()

		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusTooManyRequests, rec.Code)
		assert.Equal(t, "0", rec.Header().Get("X-RateLimit-Remaining"))
	})

	t.Run("Redis_Error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		mw := RateLimiterMiddleware(db, cfg)
		handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.3:1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock.ExpectTxPipeline()
		mock.ExpectIncr("rate_limit_192.168.1.3").SetErr(fmt.Errorf("redis down"))
		mock.ExpectTxPipelineExec()

		assert.NoError(t, handler(c))
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})
}
