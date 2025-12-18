package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
	"todo-list/internal/domain/model"
)

type AuthHandler struct {
	DB     *gorm.DB
	Secret string
}

func NewAuthHandler(db *gorm.DB, secret string) *AuthHandler {
	return &AuthHandler{DB: db, Secret: secret}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "user already exists"})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "registration successful"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	var user model.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString([]byte(h.Secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t})
}
