package v1

import (
	"net/http"
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, t usecase.UserInterface, l logger.Interface) {
	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, t, l)
	}
}
