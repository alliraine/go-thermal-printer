package service

import (
	"errors"
	"testing"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
)

type failingStatusReadWriter struct {
	writeErr error
	readErr  error
}

func (f failingStatusReadWriter) Write(p []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return len(p), nil
}

func (f failingStatusReadWriter) Read(p []byte) (int, error) {
	if f.readErr != nil {
		return 0, f.readErr
	}
	return 0, nil
}

func TestPrintServiceStatusReturnsStructuredWriteError(t *testing.T) {
	want := errors.New("write failure")
	ps := &PrintService{printer: escpos.NewESCPOS(failingStatusReadWriter{writeErr: want})}

	resp := ps.status()
	if resp.Error == nil {
		t.Fatalf("expected error response, got nil")
	}
	if !errors.Is(resp.Error, want) {
		t.Fatalf("expected error to wrap %v, got %v", want, resp.Error)
	}
	if resp.PrinterStatus != 0 || resp.OfflineStatus != 0 || resp.ErrorStatus != 0 || resp.ContinuousPaperStatus != 0 {
		t.Fatalf("expected zero status bytes when error occurs, got %+v", resp)
	}
}

func TestPrintServiceStatusReturnsStructuredReadError(t *testing.T) {
	want := errors.New("read failure")
	ps := &PrintService{printer: escpos.NewESCPOS(failingStatusReadWriter{readErr: want})}

	resp := ps.status()
	if resp.Error == nil {
		t.Fatalf("expected error response, got nil")
	}
	if !errors.Is(resp.Error, want) {
		t.Fatalf("expected error to wrap %v, got %v", want, resp.Error)
	}
}
