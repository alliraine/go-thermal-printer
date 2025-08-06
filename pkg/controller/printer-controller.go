package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jonasclaes/go-thermal-printer/pkg/dto"
	"github.com/jonasclaes/go-thermal-printer/pkg/service"
)

type PrinterController struct {
	printerService *service.PrinterService
}

func NewPrinterController(group *gin.RouterGroup, printerService *service.PrinterService) {
	controller := &PrinterController{
		printerService: printerService,
	}

	{
		printerGroup := group.Group("/printer")
		printerGroup.GET("/status", controller.getPrinterStatusHandler)
		printerGroup.POST("/print", controller.postPrinterPrintHandler)
	}
}

func (pc *PrinterController) getPrinterStatusHandler(c *gin.Context) {
	status, err := pc.printerService.GetPrinterStatus(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.PrinterStatusDto{
		PrinterStatus:         status.PrinterStatus,
		OfflineStatus:         status.OfflineStatus,
		ErrorStatus:           status.ErrorStatus,
		ContinuousPaperStatus: status.ContinuousPaperStatus,
	})
}

func (pc *PrinterController) postPrinterPrintHandler(c *gin.Context) {
	var input dto.PrinterPrintDto
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(err)
		return
	}

	err := pc.printerService.Print(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}
