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
