package v1

import (
	"net/http"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	t usecase.UserInterface
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface) {
	r := &userRoutes{t, l}

	// Задание 3: Rate Limiter — 10 запросов в минуту
	rateLimiter := utils.NewRateLimiter(10, time.Minute)

	h := handler.Group("/users")
	h.Use(rateLimiter.Middleware())
	{
		// Публичные маршруты
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		// Защищённые маршруты — требуют JWT
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		{
			protected.GET("/protected/hello", r.ProtectedFunc)

			// Задание 1: GetMe — только токен, никаких body-данных
			protected.GET("/me", r.GetMe)

			// Задание 2: PromoteUser — только для admin
			protected.PATCH("/promote/:id", utils.RoleMiddleware("admin"), r.PromoteUser)
		}
	}
}

// RegisterUser — регистрация нового пользователя
func (r *userRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := entity.User{
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully. Please check your email for verification code.",
		"session_id": sessionID,
		"user":       createdUser,
	})
}

// LoginUser — вход и получение JWT токена
func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.t.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ProtectedFunc — тестовый защищённый эндпоинт
func (r *userRoutes) ProtectedFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// GetMe — Задание 1
// Возвращает информацию о текущем пользователе
// Требования: только JWT токен в заголовке Authorization, никаких дополнительных данных
func (r *userRoutes) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	user, err := r.t.GetMe(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем всё, кроме пароля
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"verified": user.Verified,
	})
}

// PromoteUser — Задание 2
// Меняет роль пользователя на "admin"
// Доступно ТОЛЬКО пользователям с ролью admin (через RoleMiddleware)
func (r *userRoutes) PromoteUser(c *gin.Context) {
	idParam := c.Param("id")
	targetID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updatedUser, err := r.t.PromoteUser(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user": gin.H{
			"id":       updatedUser.ID,
			"username": updatedUser.Username,
			"email":    updatedUser.Email,
			"role":     updatedUser.Role,
		},
	})
}
