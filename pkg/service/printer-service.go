package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jonasclaes/go-thermal-printer/pkg/dto"
)

type PrinterService struct {
	printService *PrintService
}

func NewPrinterService(printService *PrintService) (*PrinterService, error) {
	return &PrinterService{
		printService: printService,
	}, nil
}

func (ps *PrinterService) GetPrinterStatus(c context.Context) (StatusResponse, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	status, err := ps.printService.Status(ctx)

	return status, err
}

func (ps *PrinterService) Print(c context.Context, input dto.PrinterPrintDto) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	data := make([]byte, base64.StdEncoding.DecodedLen(len(input.Data)))
	_, err := base64.StdEncoding.Decode(data, []byte(input.Data))
	if err != nil {
		return fmt.Errorf("failed to decode base64 data: %w", err)
	}

	err = ps.printService.Print(ctx, data)

	return err
}

func (ps *PrinterService) PrintTemplate(c context.Context, input dto.PrinterPrintTemplateDto) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	err := ps.printService.PrintTemplateWithVariables(ctx, input.TemplateFile, input.Variables)

	return err
}
