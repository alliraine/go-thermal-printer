package service

import (
	"context"
	"fmt"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
	"github.com/jonasclaes/go-thermal-printer/pkg/template"
	"go.bug.st/serial"
)

type PrintJob struct {
	Data     []byte
	Response chan error
}

type StatusResponse struct {
	PrinterStatus         byte
	OfflineStatus         byte
	ErrorStatus           byte
	ContinuousPaperStatus byte
	Error                 error
}

type StatusRequest struct {
	Response chan StatusResponse
}

type PrintService struct {
	port        serial.Port
	printer     *escpos.ESCPOS
	printQueue  chan PrintJob
	statusQueue chan StatusRequest
	quit        chan struct{}
}

func NewPrintService(configService *ConfigService) (*PrintService, error) {
	printerConfig := configService.GetPrinterConfig()

	mode := &serial.Mode{
		BaudRate: printerConfig.BaudRate,
		DataBits: printerConfig.DataBits,
	}

	// Convert StopBits from int to serial.StopBits
	switch printerConfig.StopBits {
	case 1:
		mode.StopBits = serial.OneStopBit
	case 2:
		mode.StopBits = serial.TwoStopBits
	default:
		mode.StopBits = serial.OneStopBit
	}

	// Convert Parity from int to serial.Parity
	switch printerConfig.Parity {
	case 0:
		mode.Parity = serial.NoParity
	case 1:
		mode.Parity = serial.OddParity
	case 2:
		mode.Parity = serial.EvenParity
	case 3:
		mode.Parity = serial.MarkParity
	case 4:
		mode.Parity = serial.SpaceParity
	default:
		mode.Parity = serial.NoParity
	}

	port, err := serial.Open(printerConfig.Port, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}

	printer := escpos.NewESCPOS(port)

	pm := &PrintService{
		port:        port,
		printer:     printer,
		printQueue:  make(chan PrintJob, 100),
		statusQueue: make(chan StatusRequest, 10),
		quit:        make(chan struct{}),
	}

	// Start the worker goroutine
	go pm.worker()

	return pm, nil
}

// worker processes all serial communication sequentially
func (ps *PrintService) worker() {
	for {
		select {
		case job := <-ps.printQueue:
			err := ps.print(job.Data)
			job.Response <- err

		case statusReq := <-ps.statusQueue:
			status := ps.status()
			statusReq.Response <- status

		case <-ps.quit:
			return
		}
	}
}

// print sends data to the printer
func (ps *PrintService) print(data []byte) error {
	_, err := ps.printer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to printer: %w", err)
	}

	return nil
}

// status retrieves the printer status
func (ps *PrintService) status() StatusResponse {
	printerStatus, err := ps.printer.PrinterStatus()
	if err != nil {
		return StatusResponse{
			Error: fmt.Errorf("failed to get printer status: %w", err),
		}
	}

	offlineStatus, err := ps.printer.OfflineStatus()
	if err != nil {
		return StatusResponse{
			Error: fmt.Errorf("failed to get offline status: %w", err),
		}
	}

	errorStatus, err := ps.printer.ErrorStatus()
	if err != nil {
		return StatusResponse{
			Error: fmt.Errorf("failed to get error status: %w", err),
		}
	}

	continuousPaperStatus, err := ps.printer.ContinuousPaperStatus()
	if err != nil {
		return StatusResponse{
			Error: fmt.Errorf("failed to get continuous paper status: %w", err),
		}
	}

	return StatusResponse{
		PrinterStatus:         printerStatus,
		OfflineStatus:         offlineStatus,
		ErrorStatus:           errorStatus,
		ContinuousPaperStatus: continuousPaperStatus,
		Error:                 nil,
	}
}

// Print queues a print job and waits for the response
func (ps *PrintService) Print(ctx context.Context, data []byte) error {
	response := make(chan error, 1)
	job := PrintJob{
		Data:     data,
		Response: response,
	}

	select {
	case ps.printQueue <- job:
		// Job queued successfully, wait for response
		select {
		case err := <-response:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Status retrieves the printer status and waits for the response
func (ps *PrintService) Status(ctx context.Context) (StatusResponse, error) {
	response := make(chan StatusResponse, 1)
	req := StatusRequest{
		Response: response,
	}

	select {
	case ps.statusQueue <- req:
		// Status request queued successfully, wait for response
		select {
		case status := <-response:
			return status, status.Error
		case <-ctx.Done():
			return StatusResponse{}, ctx.Err()
		}
	case <-ctx.Done():
		return StatusResponse{}, ctx.Err()
	}
}

// PrintTemplate renders a template and prints it to the thermal printer
func (ps *PrintService) PrintTemplate(ctx context.Context, templateContent string) error {
	data, err := template.RenderToBytes(templateContent)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return ps.Print(ctx, data)
}

// PrintTemplateWithVariables renders a template file with variables and prints it to the thermal printer
func (ps *PrintService) PrintTemplateWithVariables(ctx context.Context, templateFile string, variables map[string]string) error {
	data, err := template.RenderTemplateFileWithVariables(templateFile, variables)
	if err != nil {
		return fmt.Errorf("failed to render template with variables: %w", err)
	}

	return ps.Print(ctx, data)
}

func (ps *PrintService) Close() error {
	close(ps.quit)
	return ps.port.Close()
}
