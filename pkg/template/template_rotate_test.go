package template

import (
    "testing"
)

func TestReverseForRotateSimple(t *testing.T) {
    input := encodeToCodePage("ABC")
    got := reverseForRotate(input)
    if got != "CBA" {
        t.Fatalf("reverseForRotate simple: expected CBA, got %q", got)
    }
}

func TestReverseForRotateWithNewlines(t *testing.T) {
    input := encodeToCodePage("AB\nCD")
    got := reverseForRotate(input)
    if got != "BA\nDC" {
        t.Fatalf("reverseForRotate newline: expected BA\nDC, got %q", got)
    }
}

func TestReverseForRotateCarriageReturn(t *testing.T) {
    input := encodeToCodePage("AB\r\nCD")
    got := reverseForRotate(input)
    if got != "BA\r\nDC" {
        t.Fatalf("reverseForRotate carriage return: expected BA\r\nDC, got %q", got)
    }
}
