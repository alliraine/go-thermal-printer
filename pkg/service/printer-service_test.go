package service

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestDecodePrintPayloadPreservesDecodedBytes(t *testing.T) {
	input := []byte{0x1B, 0x40, 0x0A}
	encoded := base64.StdEncoding.EncodeToString(input)

	decoded, err := decodePrintPayload(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(decoded, input) {
		t.Fatalf("decoded bytes mismatch: got %v want %v", decoded, input)
	}
}

func TestDecodePrintPayloadTrimsWhitespace(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte{0x00, 0x01, 0x02})
	padded := "  \n" + encoded + "\n\t"

	decoded, err := decodePrintPayload(padded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(decoded) != 3 {
		t.Fatalf("expected 3 decoded bytes, got %d", len(decoded))
	}
}

func TestDecodePrintPayloadEmptyInput(t *testing.T) {
	decoded, err := decodePrintPayload("   ")
	if err != nil {
		t.Fatalf("unexpected error for empty input: %v", err)
	}
	if len(decoded) != 0 {
		t.Fatalf("expected zero-length slice, got %d", len(decoded))
	}
}

func TestDecodePrintPayloadInvalidBase64(t *testing.T) {
	if _, err := decodePrintPayload("!!"); err == nil {
		t.Fatal("expected error for invalid base64 input, got nil")
	}
}

