package main

import (
	"bytes"
	"encoding/base64"
	"log"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
)

func main() {
	var buffer bytes.Buffer

	// Create a new ESC/POS printer instance.
	e := escpos.NewESCPOS(&buffer)

	e.SelectCharacterCodePage(0x00)
	e.Initialize()

	e.Text("Hello, ESC/POS Printer!\n")

	e.UnderlineMode(escpos.Underline1DotThick)
	e.Text("Hello, ESC/POS Printer!\n")
	e.UnderlineMode(escpos.UnderlineOff)

	e.UnderlineMode(escpos.Underline2DotThick)
	e.Text("Hello, ESC/POS Printer!\n")
	e.UnderlineMode(escpos.UnderlineOff)

	e.ItalicsMode(escpos.ItalicsOn)
	e.Text("Hello, ESC/POS Printer!\n")
	e.ItalicsMode(escpos.ItalicsOff)

	e.EmphasisMode(escpos.EmphasisOn)
	e.Text("Hello, ESC/POS Printer!\n")
	e.EmphasisMode(escpos.EmphasisOff)

	e.PrintAndFeedPaperNLines(7)
	// e.SelectCutModeAndCutPaper(escpos.CutModeFull)

	e.Print()

	// print buffer as base64
	log.Printf("Print buffer (base64): %s", base64.StdEncoding.EncodeToString(buffer.Bytes()))
}
