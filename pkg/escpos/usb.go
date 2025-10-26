package escpos

import (
	"fmt"
	"os"
)

// USBTransport writes raw ESC/POS bytes to a USB printer-class node (e.g. /dev/usb/lp0).
type USBTransport struct {
	f *os.File
}

func NewUSBTransport(path string) (*USBTransport, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_SYNC, 0)
	if err != nil {
		return nil, fmt.Errorf("open usb device %s: %w", path, err)
	}
	t := &USBTransport{f: f}
	// ESC @ (initialize). Ignore error to stay non-fatal if the device buffers.
	_, _ = t.Write([]byte{0x1B, 0x40})
	return t, nil
}

func (t *USBTransport) Write(p []byte) (int, error) { return t.f.Write(p) }
func (t *USBTransport) Close() error                { return t.f.Close() }
