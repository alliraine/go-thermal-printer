package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
	"github.com/jonasclaes/go-thermal-printer/pkg/template"
	"go.bug.st/serial"
)

var (
	serialOpenFunc      = serial.Open
	usbTransportFactory = func(path string) (escpos.Transport, error) {
		return escpos.NewUSBTransport(path)
	}
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
	port            io.ReadWriter
	printer         *escpos.ESCPOS
	printQueue      chan PrintJob
	statusQueue     chan StatusRequest
	quit            chan struct{}
	statusSupported bool
}

type usbReadWriter struct {
	transport escpos.Transport
}

func (u *usbReadWriter) Write(p []byte) (int, error) {
	return u.transport.Write(p)
}

func (u *usbReadWriter) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("printer status is not supported in USB mode")
}

func (u *usbReadWriter) Close() error {
	return u.transport.Close()
}

func NewPrintService(configService *ConfigService) (*PrintService, error) {
	printerConfig := configService.GetPrinterConfig()
	appConfig := configService.GetConfig()

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

	var (
		port            io.ReadWriter
		statusSupported bool
	)

	if appConfig.TestMode {
		var buff bytes.Buffer
		port = &buff
	} else {
		path := printerConfig.Port
		if appConfig.USBMode {
			transport, err := usbTransportFactory(path)
			if err != nil {
				return nil, fmt.Errorf("failed to open usb printer: %w", err)
			}
			port = &usbReadWriter{transport: transport}
		} else {
			if strings.HasPrefix(path, "/dev/usb") || strings.HasPrefix(path, "/dev/lp") {
				file, err := os.OpenFile(path, os.O_RDWR, 0)
				if err != nil {
					return nil, fmt.Errorf("failed to open printer device file: %w", err)
				}
				port = file
			} else {
				_port, err := serialOpenFunc(path, mode)
				if err != nil {
					return nil, fmt.Errorf("failed to open serial port: %w", err)
				}
				port = _port
			}
			statusSupported = true
		}
	}

	printer := escpos.NewESCPOS(port)

	pm := &PrintService{
		port:            port,
		printer:         printer,
		printQueue:      make(chan PrintJob, 100),
		statusQueue:     make(chan StatusRequest, 10),
		quit:            make(chan struct{}),
		statusSupported: statusSupported,
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
	if len(data) == 0 {
		log.Printf("print-service: received empty print job")
		return nil
	}

	previewLen := len(data)
	if previewLen > 64 {
		previewLen = 64
	}
	tailLen := len(data)
	if tailLen > 64 {
		tailLen = 64
	}
	log.Printf(
		"print-service: writing %d bytes (head=%s tail=%s)",
		len(data),
		hex.EncodeToString(data[:previewLen]),
		hex.EncodeToString(data[len(data)-tailLen:]),
	)

	written, err := ps.printer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to printer: %w", err)
	}
	if written != len(data) {
		return fmt.Errorf("failed to write to printer: short write %d/%d", written, len(data))
	}

	log.Printf("print-service: write complete")

	return nil
}

// status retrieves the printer status
func (ps *PrintService) status() StatusResponse {
	if !ps.statusSupported {
		return StatusResponse{
			Error: fmt.Errorf("printer status is not supported for the configured transport"),
		}
	}

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
func (ps *PrintService) PrintTemplate(ctx context.Context, templateContent string, data any) error {
	renderedData, err := template.RenderToBytes(templateContent, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return ps.Print(ctx, renderedData)
}

// PrintTemplateWithVariables renders a template file with variables and prints it to the thermal printer
func (ps *PrintService) PrintTemplateWithVariables(ctx context.Context, templateFile string, variables map[string]any) error {
	renderedData, err := template.RenderTemplateFileWithVariables(templateFile, variables)
	if err != nil {
		return fmt.Errorf("failed to render template with variables: %w", err)
	}

	return ps.Print(ctx, renderedData)
}

func (ps *PrintService) Close() error {
	close(ps.quit)
	if c, ok := ps.port.(interface{ Close() error }); ok {
		return c.Close()
	}
	return nil
}
