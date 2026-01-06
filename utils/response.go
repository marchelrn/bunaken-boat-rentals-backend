package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` 
}

func APIResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func APIError(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}