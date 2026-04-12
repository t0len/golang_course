package utils

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// HashPassword хэширует пароль с помощью bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword сравнивает хэшированный пароль с введённым
func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// GenerateJWT создаёт JWT токен с user_id, role и exp
func GenerateJWT(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWTAuthMiddleware проверяет JWT токен из заголовка Authorization
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

// RoleMiddleware проверяет роль пользователя из JWT клеймов
// Задание 2: только пользователи с нужной ролью проходят дальше
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			return
		}

		if role.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied: insufficient permissions",
			})
			return
		}

		c.Next()
	}
}

// ---- Rate Limiter (Задание 3) ----

type clientEntry struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientEntry
	limit   int
	window  time.Duration
}

// NewRateLimiter создаёт rate limiter: limit запросов за window времени
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientEntry),
		limit:   limit,
		window:  window,
	}

	// Фоновая горутина очищает устаревшие записи
	go func() {
		for {
			time.Sleep(window)
			rl.mu.Lock()
			for key, c := range rl.clients {
				if time.Since(c.lastSeen) > window {
					delete(rl.clients, key)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

// Middleware возвращает gin middleware для ограничения запросов
// Аутентифицированный: ключ = userID из JWT
// Анонимный: ключ = ClientIP
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Определяем ключ идентификации
		key := c.ClientIP()
		if userID, exists := c.Get("userID"); exists {
			key = userID.(string)
		}

		rl.mu.Lock()
		entry, exists := rl.clients[key]
		if !exists || time.Since(entry.lastSeen) > rl.window {
			// Новый клиент или окно истекло — сбрасываем счётчик
			rl.clients[key] = &clientEntry{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			c.Next()
			return
		}

		entry.count++
		entry.lastSeen = time.Now()

		if entry.count > rl.limit {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}
		rl.mu.Unlock()

		c.Next()
	}
}
