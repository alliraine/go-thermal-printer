package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"testing"

	"github.com/jonasclaes/go-thermal-printer/pkg/dto"
	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
)

func newBenchmarkPrinterService(b testing.TB) *PrinterService {
	b.Helper()

	buffer := &bytes.Buffer{}
	ps := &PrintService{
		port:        buffer,
		printer:     escpos.NewESCPOS(buffer),
		printQueue:  make(chan PrintJob, 16),
		statusQueue: make(chan StatusRequest, 1),
		quit:        make(chan struct{}),
	}

	go ps.worker()

	printerService, err := NewPrinterService(ps)
	if err != nil {
		b.Fatalf("failed to create printer service: %v", err)
	}

	previousWriter := log.Writer()
	log.SetOutput(io.Discard)

	b.Cleanup(func() {
		log.SetOutput(previousWriter)
		_ = ps.Close()
	})

	return printerService
}

// BenchmarkPrintImageLegacy measures the previous /print-image pipeline that
// performed a base64 round-trip before queuing the print job.
func BenchmarkPrintImageLegacy(b *testing.B) {
	printerService := newBenchmarkPrinterService(b)
	payload := make([]byte, 4096)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoded := base64.StdEncoding.EncodeToString(payload)
		if err := printerService.Print(ctx, dto.PrinterPrintDto{Data: encoded}); err != nil {
			b.Fatalf("print failed: %v", err)
		}
	}
}

// BenchmarkPrintImageDirect exercises the new PrintBytes helper which bypasses
// the base64 encode/decode cycle by sending the raster bytes directly.
func BenchmarkPrintImageDirect(b *testing.B) {
	printerService := newBenchmarkPrinterService(b)
	payload := make([]byte, 4096)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := printerService.PrintBytes(ctx, payload); err != nil {
			b.Fatalf("print failed: %v", err)
		}
	}
}
