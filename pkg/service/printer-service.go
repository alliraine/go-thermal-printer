package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
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

	data, err := decodePrintPayload(input.Data)
	if err != nil {
		return err
	}

	return ps.printService.Print(ctx, data)
}

func (ps *PrinterService) PrintTemplate(c context.Context, input dto.PrinterPrintTemplateDto) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	err := ps.printService.PrintTemplateWithVariables(ctx, input.TemplateFile, input.Variables)

	return err
}

func decodePrintPayload(encoded string) ([]byte, error) {
	trimmed := strings.TrimSpace(encoded)
	if trimmed == "" {
		return []byte{}, nil
	}

	buf := make([]byte, base64.StdEncoding.DecodedLen(len(trimmed)))
	n, err := base64.StdEncoding.Decode(buf, []byte(trimmed))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return buf[:n], nil
}
