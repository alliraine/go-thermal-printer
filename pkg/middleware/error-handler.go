package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jonasclaes/go-thermal-printer/pkg/common"
)

type ErrorHandlerMiddleware struct{}

func NewErrorHandlerMiddleware() *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{}
}

func (m *ErrorHandlerMiddleware) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			var appErr common.AppError
			if errors.As(err, &appErr) {
				errorResponse(c, appErr.HttpStatusCode(), appErr.Error())
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	if len(message) > 0 {
		message = strings.ToUpper(message[:1]) + message[1:]
	}
	c.JSON(statusCode, gin.H{"error": message})
}
