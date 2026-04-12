package v1

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, errorResponse{msg})
}
