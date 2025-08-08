package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/jonasclaes/go-thermal-printer/pkg/template"
)

func main() {
	// Example template content
	templateContent := `<Bold>This is a test</Bold>
<Underline>Test underline</Underline>
<Italic>This is italic text</Italic>
Normal text here
<Bold><Underline>Bold and underlined</Underline></Bold>
<FontB>This uses font B</FontB>`

	// Render the template to bytes
	data, err := template.RenderToBytes(templateContent)
	if err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}

	// Print the raw bytes (for debugging)
	fmt.Printf("Generated %d bytes:\n", len(data))
	for i, b := range data {
		if i > 0 && i%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02X ", b)
	}
	fmt.Println()

	// Print a readable representation
	fmt.Printf("\nReadable representation:\n")
	for _, b := range data {
		if b >= 32 && b <= 126 {
			fmt.Printf("%c", b)
		} else {
			fmt.Printf("[%02X]", b)
		}
	}
	fmt.Println()

	fmt.Printf("\nBase64 encoded output:\n%s\n", base64.StdEncoding.EncodeToString(data))

	fmt.Println()

	fmt.Println("\nTemplate rendered successfully!")
	fmt.Println("You can now use this byte array with your PrintService.Print() method.")
}
