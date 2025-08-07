package escpos

import (
	"io"
	"log"

	"go.bug.st/serial"
)

type ESCPOS struct {
	rw io.ReadWriter
}

func NewESCPOS(rw io.ReadWriter) (escpos *ESCPOS) {
	escpos = &ESCPOS{
		rw: rw,
	}
	return
}

// Write raw bytes to the printer.
func (p *ESCPOS) Write(data []byte) (int, error) {
	if len(data) > 0 {
		return p.rw.Write(data)
	}
	return 0, nil
}

// Reads raw bytes from the printer.
func (p *ESCPOS) Read(length int) ([]byte, error) {
	buf := make([]byte, length)
	n, err := p.rw.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// Reads a single byte from the printer.
func (p *ESCPOS) ReadByte() (byte, error) {
	data, err := p.Read(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (p *ESCPOS) Text(data string) (int, error) {
	return p.Write([]byte(data))
}

func (p *ESCPOS) resetInputBuffer() {
	if serialPort, ok := p.rw.(serial.Port); ok {
		err := serialPort.ResetInputBuffer()
		if err != nil {
			log.Fatalf("failed to reset input buffer: %v", err)
		}
	}
}

func (escpos *ESCPOS) status(statusByte byte) (byte, error) {
	escpos.resetInputBuffer()

	_, err := escpos.rw.Write([]byte{0x10, 0x04, statusByte})
	if err != nil {
		log.Fatalf("failed to send status command: %v", err)
	}

	data, err := escpos.ReadByte()
	if err != nil {
		log.Fatalf("failed to read status response: %v", err)
	}

	return data, nil
}

// Control Commands

func (p *ESCPOS) PrinterStatus() (byte, error) {
	return p.status(0x01)
}

func (p *ESCPOS) OfflineStatus() (byte, error) {
	return p.status(0x02)
}

func (p *ESCPOS) ErrorStatus() (byte, error) {
	return p.status(0x03)
}

func (p *ESCPOS) ContinuousPaperStatus() (byte, error) {
	return p.status(0x04)
}

func (p *ESCPOS) isPrinterStatus(mask byte) bool {
	status, err := p.PrinterStatus()
	if err != nil {
		log.Fatalf("failed to get printer status: %v", err)
	}

	return status&mask != 0
}

func (p *ESCPOS) isOfflineStatus(mask byte) bool {
	status, err := p.OfflineStatus()
	if err != nil {
		log.Fatalf("failed to get offline status: %v", err)
	}

	return status&mask != 0
}

func (p *ESCPOS) isErrorStatus(mask byte) bool {
	status, err := p.ErrorStatus()
	if err != nil {
		log.Fatalf("failed to get error status: %v", err)
	}

	return status&mask != 0
}

func (p *ESCPOS) isContinuousPaperStatus(mask byte) bool {
	status, err := p.ContinuousPaperStatus()
	if err != nil {
		log.Fatalf("failed to get continuous paper status: %v", err)
	}

	return status&mask != 0
}

func (p *ESCPOS) IsDrawerOpenCloseSignalHigh() bool {
	return p.isPrinterStatus(0x04)
}

func (p *ESCPOS) IsOffline() bool {
	return p.isPrinterStatus(0x08)
}

func (p *ESCPOS) IsCoverOpen() bool {
	return p.isOfflineStatus(0x04)
}

func (p *ESCPOS) IsPaperBeingFedByFeedButton() bool {
	return p.isOfflineStatus(0x08)
}

func (p *ESCPOS) IsPrintingBeingStopped() bool {
	return p.isOfflineStatus(0x20)
}

func (p *ESCPOS) IsAutocutterError() bool {
	return p.isErrorStatus(0x08)
}

func (p *ESCPOS) IsUnrecoverableError() bool {
	return p.isErrorStatus(0x20)
}

func (p *ESCPOS) IsAutoRecoverableError() bool {
	return p.isErrorStatus(0x40)
}

func (p *ESCPOS) IsPaperNearEnd() bool {
	return p.isContinuousPaperStatus(0x0C)
}

func (p *ESCPOS) IsPaperEnd() bool {
	return p.isContinuousPaperStatus(0x60)
}

// Font Commands

func (p *ESCPOS) Initialize() (int, error) {
	return p.Write([]byte{0x1B, 0x40})
}

func (p *ESCPOS) UnderlineMode(mode UnderlineMode) (int, error) {
	return p.Write([]byte{0x1B, 0x2D, byte(mode)})
}

func (p *ESCPOS) ItalicsMode(mode ItalicsMode) (int, error) {
	return p.Write([]byte{0x1B, 0x34, byte(mode)})
}

func (p *ESCPOS) EmphasisMode(mode EmphasisMode) (int, error) {
	return p.Write([]byte{0x1B, 0x45, byte(mode)})
}

func (p *ESCPOS) SelectCharacterFont(font CharacterFont) (int, error) {
	return p.Write([]byte{0x1B, 0x4D, byte(font)})
}

func (p *ESCPOS) SelectCharacterCodePage(codePage CharacterCodePage) (int, error) {
	return p.Write([]byte{0x1B, 0x74, byte(codePage)})
}

// Paper Movement Commands

func (p *ESCPOS) FullCut() (int, error) {
	return p.Write([]byte{0x1B, 0x6D})
}

func (p *ESCPOS) SelectCutModeAndCutPaper(cutMode CutMode) (int, error) {
	return p.Write([]byte{0x1D, 0x56, byte(cutMode)})
}

func (p *ESCPOS) PrintAndFeedPaperNLines(n int) (int, error) {
	return p.Write([]byte{0x1B, 0x64, byte(n)})
}

// Cursor Position Commands

func (p *ESCPOS) LineFeed() (int, error) {
	return p.Write([]byte{0x0A})
}

func (p *ESCPOS) FormFeed() (int, error) {
	return p.Write([]byte{0x0C})
}
