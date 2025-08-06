package escpos

import (
	"bufio"
	"io"
	"log"
)

type ESCPOS struct {
	rw          io.ReadWriter
	printStream *bufio.Writer
}

// Create a new ESC/POS printer instance.
func NewESCPOS(rw io.ReadWriter) (escpos *ESCPOS) {

	escpos = &ESCPOS{
		rw:          rw,
		printStream: bufio.NewWriter(rw),
	}
	return
}

// Sends the buffered data to the printer.
func (escpos *ESCPOS) Print() error {
	return escpos.printStream.Flush()
}

// Write raw bytes to the printer.
func (escpos *ESCPOS) Write(data []byte) (int, error) {
	if len(data) > 0 {
		return escpos.printStream.Write(data)
	}
	return 0, nil
}

// Writes a string using the predefined options.
func (escpos *ESCPOS) Text(data string) (int, error) {
	return escpos.Write([]byte(data))
}

func (escpos *ESCPOS) clearPendingDataInReadBuffer() {
	// Clear any pending data in the read buffer
	buffer := make([]byte, 1024)
	totalCleared := 0
	for {
		n, err := escpos.rw.Read(buffer)
		if err != nil || n == 0 {
			break // No more data to read or an error occurred
		}
		totalCleared += n
		log.Printf("Cleared %d stale bytes from read buffer", n)
	}
}

func (escpos *ESCPOS) status(statusByte byte) (byte, error) {
	escpos.clearPendingDataInReadBuffer()

	_, err := escpos.rw.Write([]byte{0x10, 0x04, statusByte})
	if err != nil {
		log.Fatalf("failed to send status command: %v", err)
	}

	data := make([]byte, 1)
	_, err = escpos.rw.Read(data)
	if err != nil {
		log.Fatalf("failed to read status response: %v", err)
	}

	return data[0], nil
}

// Control Commands

func (escpos *ESCPOS) PrinterStatus() (byte, error) {
	return escpos.status(0x01)
}

func (escpos *ESCPOS) OfflineStatus() (byte, error) {
	return escpos.status(0x02)
}

func (escpos *ESCPOS) ErrorStatus() (byte, error) {
	return escpos.status(0x03)
}

func (escpos *ESCPOS) ContinuousPaperStatus() (byte, error) {
	return escpos.status(0x04)
}

func (escpos *ESCPOS) isPrinterStatus(mask byte) bool {
	status, err := escpos.PrinterStatus()
	if err != nil {
		log.Fatalf("failed to get printer status: %v", err)
	}

	return status&mask != 0
}

func (escpos *ESCPOS) isOfflineStatus(mask byte) bool {
	status, err := escpos.OfflineStatus()
	if err != nil {
		log.Fatalf("failed to get offline status: %v", err)
	}

	return status&mask != 0
}

func (escpos *ESCPOS) isErrorStatus(mask byte) bool {
	status, err := escpos.ErrorStatus()
	if err != nil {
		log.Fatalf("failed to get error status: %v", err)
	}

	return status&mask != 0
}

func (escpos *ESCPOS) isContinuousPaperStatus(mask byte) bool {
	status, err := escpos.ContinuousPaperStatus()
	if err != nil {
		log.Fatalf("failed to get continuous paper status: %v", err)
	}

	return status&mask != 0
}

func (escpos *ESCPOS) IsDrawerOpenCloseSignalHigh() bool {
	return escpos.isPrinterStatus(0x04)
}

func (escpos *ESCPOS) IsOffline() bool {
	return escpos.isPrinterStatus(0x08)
}

func (escpos *ESCPOS) IsCoverOpen() bool {
	return escpos.isOfflineStatus(0x04)
}

func (escpos *ESCPOS) IsPaperBeingFedByFeedButton() bool {
	return escpos.isOfflineStatus(0x08)
}

func (escpos *ESCPOS) IsPrintingBeingStopped() bool {
	return escpos.isOfflineStatus(0x20)
}

func (escpos *ESCPOS) IsAutocutterError() bool {
	return escpos.isErrorStatus(0x08)
}

func (escpos *ESCPOS) IsUnrecoverableError() bool {
	return escpos.isErrorStatus(0x20)
}

func (escpos *ESCPOS) IsAutoRecoverableError() bool {
	return escpos.isErrorStatus(0x40)
}

func (escpos *ESCPOS) IsPaperNearEnd() bool {
	return escpos.isContinuousPaperStatus(0x0C)
}

func (escpos *ESCPOS) IsPaperEnd() bool {
	return escpos.isContinuousPaperStatus(0x60)
}

// Font Commands

func (escpos *ESCPOS) Initialize() (int, error) {
	return escpos.Write([]byte{0x1B, 0x40})
}

func (escpos *ESCPOS) UnderlineMode(mode UnderlineMode) (int, error) {
	return escpos.Write([]byte{0x1B, 0x2D, byte(mode)})
}

func (escpos *ESCPOS) ItalicsMode(mode ItalicsMode) (int, error) {
	return escpos.Write([]byte{0x1B, 0x34, byte(mode)})
}

func (escpos *ESCPOS) EmphasisMode(mode EmphasisMode) (int, error) {
	return escpos.Write([]byte{0x1B, 0x45, byte(mode)})
}

func (escpos *ESCPOS) SelectCharacterFont(font CharacterFont) (int, error) {
	return escpos.Write([]byte{0x1B, 0x4D, byte(font)})
}

func (escpos *ESCPOS) SelectCharacterCodePage(codePage CharacterCodePage) (int, error) {
	return escpos.Write([]byte{0x1B, 0x74, byte(codePage)})
}

// Paper Movement Commands

func (escpos *ESCPOS) FullCut() (int, error) {
	return escpos.Write([]byte{0x1B, 0x6D})
}

func (escpos *ESCPOS) SelectCutModeAndCutPaper(cutMode CutMode) (int, error) {
	return escpos.Write([]byte{0x1D, 0x56, byte(cutMode)})
}

func (escpos *ESCPOS) PrintAndFeedPaperNLines(n int) (int, error) {
	return escpos.Write([]byte{0x1B, 0x64, byte(n)})
}

// Cursor Position Commands

func (escpos *ESCPOS) LineFeed() (int, error) {
	return escpos.Write([]byte{0x0A})
}

func (escpos *ESCPOS) FormFeed() (int, error) {
	return escpos.Write([]byte{0x0C})
}
