package escpos

import (
	"fmt"
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

// Write raw bytes to the printer ensuring the full payload is delivered.
func (p *ESCPOS) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	total := 0
	for total < len(data) {
		n, err := p.rw.Write(data[total:])
		if n > 0 {
			total += n
		}
		if err != nil {
			return total, err
		}
		if n == 0 {
			return total, io.ErrShortWrite
		}
	}

	return total, nil
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
			log.Printf("warning: failed to reset input buffer: %v", err)
		}
	}
}

func (escpos *ESCPOS) status(statusByte byte) (byte, error) {
	escpos.resetInputBuffer()

	if _, err := escpos.rw.Write([]byte{0x10, 0x04, statusByte}); err != nil {
		return 0, fmt.Errorf("failed to send status command: %w", err)
	}

	data, err := escpos.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("failed to read status response: %w", err)
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

func (p *ESCPOS) isPrinterStatus(mask byte) (bool, error) {
	status, err := p.PrinterStatus()
	if err != nil {
		return false, err
	}

	return status&mask != 0, nil
}

func (p *ESCPOS) isOfflineStatus(mask byte) (bool, error) {
	status, err := p.OfflineStatus()
	if err != nil {
		return false, err
	}

	return status&mask != 0, nil
}

func (p *ESCPOS) isErrorStatus(mask byte) (bool, error) {
	status, err := p.ErrorStatus()
	if err != nil {
		return false, err
	}

	return status&mask != 0, nil
}

func (p *ESCPOS) isContinuousPaperStatus(mask byte) (bool, error) {
	status, err := p.ContinuousPaperStatus()
	if err != nil {
		return false, err
	}

	return status&mask != 0, nil
}

func (p *ESCPOS) IsDrawerOpenCloseSignalHigh() (bool, error) {
	return p.isPrinterStatus(0x04)
}

func (p *ESCPOS) IsOffline() (bool, error) {
	return p.isPrinterStatus(0x08)
}

func (p *ESCPOS) IsCoverOpen() (bool, error) {
	return p.isOfflineStatus(0x04)
}

func (p *ESCPOS) IsPaperBeingFedByFeedButton() (bool, error) {
	return p.isOfflineStatus(0x08)
}

func (p *ESCPOS) IsPrintingBeingStopped() (bool, error) {
	return p.isOfflineStatus(0x20)
}

func (p *ESCPOS) IsAutocutterError() (bool, error) {
	return p.isErrorStatus(0x08)
}

func (p *ESCPOS) IsUnrecoverableError() (bool, error) {
	return p.isErrorStatus(0x20)
}

func (p *ESCPOS) IsAutoRecoverableError() (bool, error) {
	return p.isErrorStatus(0x40)
}

func (p *ESCPOS) IsPaperNearEnd() (bool, error) {
	return p.isContinuousPaperStatus(0x0C)
}

func (p *ESCPOS) IsPaperEnd() (bool, error) {
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
