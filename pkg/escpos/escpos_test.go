package escpos

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type steppingWriter struct {
	limit   int
	written bytes.Buffer
}

func (s *steppingWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	step := s.limit
	if step <= 0 || step > len(p) {
		step = len(p)
	}
	s.written.Write(p[:step])
	return step, nil
}

func (s *steppingWriter) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func TestESCPOSWriteRetriesPartialWrites(t *testing.T) {
	writer := &steppingWriter{limit: 3}
	esc := NewESCPOS(writer)
	payload := []byte("partial writes need retries")

	written, err := esc.Write(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if written != len(payload) {
		t.Fatalf("expected %d bytes written, got %d", len(payload), written)
	}
	if got := writer.written.Bytes(); !bytes.Equal(got, payload) {
		t.Fatalf("writer saw %q, want %q", got, payload)
	}
}

type zeroProgressWriter struct{}

func (zeroProgressWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func (zeroProgressWriter) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func TestESCPOSWriteDetectsStalledWriter(t *testing.T) {
	esc := NewESCPOS(zeroProgressWriter{})

	_, err := esc.Write([]byte("abc"))
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("expected io.ErrShortWrite, got %v", err)
	}
}

type failingStatusRW struct {
	writeErr error
	readErr  error
	payload  []byte
}

func (f *failingStatusRW) Write(p []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return len(p), nil
}

func (f *failingStatusRW) Read(p []byte) (int, error) {
	if f.readErr != nil {
		return 0, f.readErr
	}
	if len(f.payload) == 0 {
		return 0, io.EOF
	}
	n := copy(p, f.payload)
	f.payload = f.payload[n:]
	return n, nil
}

func TestESCPOSStatusReturnsWriteError(t *testing.T) {
	want := errors.New("boom write")
	esc := NewESCPOS(&failingStatusRW{writeErr: want})

	if _, err := esc.status(0x01); !errors.Is(err, want) {
		t.Fatalf("expected wrapped write error, got %v", err)
	}
}

func TestESCPOSStatusReturnsReadError(t *testing.T) {
	want := errors.New("boom read")
	esc := NewESCPOS(&failingStatusRW{payload: []byte{0x00}, readErr: want})

	if _, err := esc.status(0x01); !errors.Is(err, want) {
		t.Fatalf("expected wrapped read error, got %v", err)
	}
}
