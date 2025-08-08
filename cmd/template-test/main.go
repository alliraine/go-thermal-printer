package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/jonasclaes/go-thermal-printer/pkg/template"
)

func main() {
	// Example template content using new syntax
	templateContent := `{{bold "This is a test"}}
{{underline "Test underline"}}
{{italic "This is italic text"}}
Normal text here
{{bold (underline "Bold and underlined")}}`

	// Render the template to bytes (passing nil for data since no variables used)
	data, err := template.RenderToBytes(templateContent, nil)
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

	// Test with variables
	fmt.Println("\nTesting template with variables:")
	templateWithVars := `Store: {{bold .storeName}}
Total: {{bold (printf "$%.2f" .total)}}
Customer: {{underline .customerName}}`

	testData := map[string]interface{}{
		"storeName":    "Jonas' Store",
		"total":        29.99,
		"customerName": "John Doe",
	}

	varData, err := template.RenderToBytes(templateWithVars, testData)
	if err != nil {
		log.Fatalf("Failed to render template with variables: %v", err)
	}

	fmt.Printf("Generated %d bytes with variables:\n", len(varData))
	for _, b := range varData {
		if b >= 32 && b <= 126 {
			fmt.Printf("%c", b)
		} else {
			fmt.Printf("[%02X]", b)
		}
	}
	fmt.Println()

	fmt.Println("\nTemplate rendered successfully!")
	fmt.Println("You can now use this byte array with your PrintService.Print() method.")
}
