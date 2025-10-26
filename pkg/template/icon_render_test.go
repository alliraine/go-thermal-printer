package template

import (
	"bytes"
	"testing"
)

func TestSanitizeIconKey(t *testing.T) {
	got := sanitizeIconKey("Action-Face")
	if got != "actionface" {
		t.Fatalf("sanitizeIconKey: got %q", got)
	}
}

func TestLookupIconData(t *testing.T) {
	data, canonical, err := lookupIconData("ActionFace")
	if err != nil {
		t.Fatalf("lookupIconData failed: %v", err)
	}
	if canonical != "ActionFace" {
		t.Fatalf("expected canonical name ActionFace, got %s", canonical)
	}
	if len(data) == 0 {
		t.Fatal("expected icon data")
	}

	sanitizedData, _, err := lookupIconData("action-face")
	if err != nil {
		t.Fatalf("lookupIconData sanitised: %v", err)
	}
	if !bytes.Equal(data, sanitizedData) {
		t.Fatal("sanitized lookup returned different data")
	}
}

func TestIconTemplateFunc(t *testing.T) {
	result, err := iconTemplateFunc("ActionFace")
	if err != nil {
		t.Fatalf("iconTemplateFunc error: %v", err)
	}
	payload := []byte(result)
	if len(payload) < 5 {
		t.Fatalf("icon payload too small: %d", len(payload))
	}
	if payload[0] != 0x1B || payload[1] != 0x40 {
		t.Fatalf("expected payload to start with ESC @, got % x", payload[:2])
	}
	if payload[len(payload)-3] != 0x1B || payload[len(payload)-2] != 0x74 {
		t.Fatalf("expected payload to end with ESC t, got % x", payload[len(payload)-3:])
	}
	feed := payload[len(payload)-4]
	if feed != 0x00 {
		t.Fatalf("expected default feed 0, got %d", feed)
	}

	custom, err := iconTemplateFunc("ActionFace", 64, 3)
	if err != nil {
		t.Fatalf("iconTemplateFunc custom error: %v", err)
	}
	customPayload := []byte(custom)
	if customPayload[len(customPayload)-4] != 0x03 {
		t.Fatalf("expected custom feed 3, got %d", customPayload[len(customPayload)-4])
	}
}
