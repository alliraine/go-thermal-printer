package service

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
	"github.com/jonasclaes/go-thermal-printer/pkg/model"
	"go.bug.st/serial"
)

type stubSerialPort struct {
	mu      sync.Mutex
	writes  [][]byte
	closed  bool
	timeout time.Duration
}

func (s *stubSerialPort) SetMode(mode *serial.Mode) error { return nil }

func (s *stubSerialPort) Read(p []byte) (int, error) { return 0, io.EOF }

func (s *stubSerialPort) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writes = append(s.writes, append([]byte(nil), p...))
	return len(p), nil
}

func (s *stubSerialPort) Drain() error { return nil }

func (s *stubSerialPort) ResetInputBuffer() error { return nil }

func (s *stubSerialPort) ResetOutputBuffer() error { return nil }

func (s *stubSerialPort) SetDTR(bool) error { return nil }

func (s *stubSerialPort) SetRTS(bool) error { return nil }

func (s *stubSerialPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return &serial.ModemStatusBits{}, nil
}

func (s *stubSerialPort) SetReadTimeout(t time.Duration) error {
	s.timeout = t
	return nil
}

func (s *stubSerialPort) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errors.New("already closed")
	}
	s.closed = true
	return nil
}

func (s *stubSerialPort) Break(time.Duration) error { return nil }

type stubUSBTransport struct {
	mu     sync.Mutex
	writes [][]byte
	closed bool
}

func (s *stubUSBTransport) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writes = append(s.writes, append([]byte(nil), p...))
	return len(p), nil
}

func (s *stubUSBTransport) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errors.New("already closed")
	}
	s.closed = true
	return nil
}

func TestNewPrintServiceSerialModeInitializes(t *testing.T) {
	originalSerial := serialOpenFunc
	serialPort := &stubSerialPort{}
	serialOpenFunc = func(name string, mode *serial.Mode) (serial.Port, error) {
		if name != "/dev/ttyS0" {
			t.Fatalf("unexpected serial port path: %s", name)
		}
		if mode.BaudRate != 19200 {
			t.Fatalf("unexpected baud rate: %d", mode.BaudRate)
		}
		return serialPort, nil
	}
	t.Cleanup(func() { serialOpenFunc = originalSerial })

	cfg := &ConfigService{config: &model.AppConfig{
		Printer: model.PrinterConfig{
			Port:     "/dev/ttyS0",
			BaudRate: 19200,
			DataBits: 8,
			StopBits: 1,
			Parity:   0,
		},
	}}

	svc, err := NewPrintService(cfg)
	if err != nil {
		t.Fatalf("expected serial print service to initialize: %v", err)
	}
	t.Cleanup(func() { _ = svc.Close() })

	if !svc.statusSupported {
		t.Fatalf("expected serial transport to support status polling")
	}

	if port, ok := svc.port.(*stubSerialPort); !ok || port != serialPort {
		t.Fatalf("expected stub serial port to be used, got %#v", svc.port)
	}
}

func TestNewPrintServiceUSBModeInitializes(t *testing.T) {
	originalSerial := serialOpenFunc
	serialOpenFunc = func(name string, mode *serial.Mode) (serial.Port, error) {
		t.Fatalf("serialOpenFunc should not be called in USB mode")
		return nil, nil
	}
	t.Cleanup(func() { serialOpenFunc = originalSerial })

	transport := &stubUSBTransport{}
	var receivedPath string
	originalUSBFactory := usbTransportFactory
	usbTransportFactory = func(path string) (escpos.Transport, error) {
		receivedPath = path
		return transport, nil
	}
	t.Cleanup(func() { usbTransportFactory = originalUSBFactory })

	cfg := &ConfigService{config: &model.AppConfig{
		USBMode: true,
		Printer: model.PrinterConfig{
			Port: "/dev/usb/lp0",
		},
	}}

	svc, err := NewPrintService(cfg)
	if err != nil {
		t.Fatalf("expected usb print service to initialize: %v", err)
	}
	t.Cleanup(func() { _ = svc.Close() })

	if receivedPath != "/dev/usb/lp0" {
		t.Fatalf("usb transport factory received wrong path: %s", receivedPath)
	}

	if svc.statusSupported {
		t.Fatalf("expected usb transport to disable status polling")
	}

	if _, ok := svc.port.(*usbReadWriter); !ok {
		t.Fatalf("expected usbReadWriter to wrap usb transport, got %#v", svc.port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := svc.Status(ctx)
	if err == nil {
		t.Fatalf("expected status polling to return an error in USB mode")
	}
	if res.Error == nil {
		t.Fatalf("expected response error to be populated in USB mode")
	}
}
