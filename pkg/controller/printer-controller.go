package controller

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jonasclaes/go-thermal-printer/pkg/dto"
	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
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
		printerGroup.POST("/print-template", controller.postPrinterPrintTemplateHandler)
		printerGroup.POST("/print-image", controller.PrintImage)

	}
}

// @Summary		Query printer status
// @Description	Query the printer status through the configured port.
// @Tags			Printer
// @Security ApiKeyAuth
// @Success		200	{object}	dto.PrinterStatusDto
// @Router			/api/v1/printer/status [get]
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

// @Summary		Print an array of bytes
// @Description	Print an array of bytes to the printer, with ESC/POS commands.
// @Tags			Printer
// @Security ApiKeyAuth
// @Param request body dto.PrinterPrintDto	true "Printer data"
// @Success		201
// @Router			/api/v1/printer/print [post]
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

// @Summary		Print a template
// @Description	Print a template with arbitrary data.
// @Tags			Printer
// @Security ApiKeyAuth
// @Param request body dto.PrinterPrintTemplateDto	true "Printer data"
// @Success		201
// @Router			/api/v1/printer/print-template [post]
func (pc *PrinterController) postPrinterPrintTemplateHandler(c *gin.Context) {
	var input dto.PrinterPrintTemplateDto
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(err)
		return
	}

	err := pc.printerService.PrintTemplate(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

// PrintImage: POST /api/v1/printer/print-image
// Body: { "imageBase64": "<...>", "maxWidthDots": 384 }
func (pc *PrinterController) PrintImage(c *gin.Context) {
	var req dto.PrintImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}
	bytes, err := escpos.EncodeImageToRasterBytes(req.ImageBase64, req.MaxWidthDots)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert image: " + err.Error()})
		return
	}

	payload := dto.PrinterPrintDto{Data: base64.StdEncoding.EncodeToString(bytes)}
	if err := pc.printerService.Print(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "print failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "bytesWritten": len(bytes)})
}
