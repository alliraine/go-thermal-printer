package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/jonasclaes/go-thermal-printer/pkg/controller"
	"github.com/jonasclaes/go-thermal-printer/pkg/middleware"

	// gin-swagger middleware
	ginSwagger "github.com/swaggo/gin-swagger"

	// swagger embed files
	swaggerFiles "github.com/swaggo/files"

	docs "github.com/jonasclaes/go-thermal-printer/pkg/docs"
)

func initRouter(svc *services) (*gin.Engine, error) {
	docs.SwaggerInfo.Host = svc.configService.GetServerConfig().SwaggerHost

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
