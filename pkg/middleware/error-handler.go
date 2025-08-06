package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorHandlerMiddleware struct{}

func NewErrorHandlerMiddleware() *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{}
}

func (m *ErrorHandlerMiddleware) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}
