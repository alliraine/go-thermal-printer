package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jonasclaes/go-thermal-printer/pkg/common"
	"github.com/jonasclaes/go-thermal-printer/pkg/service"
)

type ApiKeyMiddleware struct {
	configService *service.ConfigService
}

func NewApiKeyMiddleware(configService *service.ConfigService) *ApiKeyMiddleware {
	return &ApiKeyMiddleware{
		configService: configService,
	}
}

func (m *ApiKeyMiddleware) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := m.Verify(c)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}

		c.Next()
	}
}

func (m *ApiKeyMiddleware) Verify(c *gin.Context) error {
	apiKey := c.GetHeader("X-Api-Key")

	if apiKey != m.configService.GetServerConfig().ApiKey {
		return &common.InvalidAPIKeyError{}
	}

	return nil
}
