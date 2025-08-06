package service

import (
	"context"
	"fmt"
	"io"

	"github.com/jacobsa/go-serial/serial"
	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
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
	port        io.ReadWriteCloser
	printer     *escpos.ESCPOS
	printQueue  chan PrintJob
	statusQueue chan StatusRequest
	quit        chan struct{}
}

func NewPrintService() (*PrintService, error) {
	options := serial.OpenOptions{
		PortName:              "COM5",
		BaudRate:              19200,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	port, err := serial.Open(options)
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
func (pm *PrintService) worker() {
	for {
		select {
		case job := <-pm.printQueue:
			err := pm.print(job.Data)
			job.Response <- err

		case statusReq := <-pm.statusQueue:
			status := pm.status()
			statusReq.Response <- status

		case <-pm.quit:
			return
		}
	}
}

// print sends data to the printer
func (pm *PrintService) print(data []byte) error {
	_, err := pm.port.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to printer: %w", err)
	}

	readData := make([]byte, 1)
	_, err = pm.port.Read(readData)
	if err != nil {
		return fmt.Errorf("failed to read from printer: %w", err)
	}

	return nil
}

// status retrieves the printer status
func (pm *PrintService) status() StatusResponse {
	printerStatus, err := pm.printer.PrinterStatus()
	if err != nil {
		return StatusResponse{
			Error: fmt.Errorf("failed to get printer status: %w", err),
		}
	}

	offlineStatus, err := pm.printer.OfflineStatus()
	if err != nil {
		return StatusResponse{
			Error: fmt.Errorf("failed to get offline status: %w", err),
		}
	}

	errorStatus, err := pm.printer.ErrorStatus()
	if err != nil {
		return StatusResponse{
			Error: fmt.Errorf("failed to get error status: %w", err),
		}
	}

	continuousPaperStatus, err := pm.printer.ContinuousPaperStatus()
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
func (pm *PrintService) Print(ctx context.Context, data []byte) error {
	response := make(chan error, 1)
	job := PrintJob{
		Data:     data,
		Response: response,
	}

	select {
	case pm.printQueue <- job:
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
func (pm *PrintService) Status(ctx context.Context) (StatusResponse, error) {
	response := make(chan StatusResponse, 1)
	req := StatusRequest{
		Response: response,
	}

	select {
	case pm.statusQueue <- req:
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

// Close shuts down the printer manager
func (pm *PrintService) Close() error {
	close(pm.quit)
	return pm.port.Close()
}
