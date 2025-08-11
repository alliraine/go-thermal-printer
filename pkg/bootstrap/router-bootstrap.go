package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/jonasclaes/go-thermal-printer/pkg/controller"
	"github.com/jonasclaes/go-thermal-printer/pkg/middleware"
)

func initRouter(svc *services) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(middleware.NewErrorHandlerMiddleware().Add())

	apiKeyMiddleware := middleware.NewApiKeyMiddleware(svc.configService).Add()

	rootGroup := router.Group("/")

	{
		root := router.Group("/")
		controller.NewHealthController(rootGroup)

		api := root.Group("/api", apiKeyMiddleware)
		v1 := api.Group("/v1")
		controller.NewPrinterController(v1, svc.printerService)
	}

	return router, nil
}
