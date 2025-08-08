package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/jonasclaes/go-thermal-printer/pkg/template"
)

func main() {
	// Test template with variables
	templateContent := `<Bold>{{storeName}}</Bold>
Date: {{date}}
<Underline>Order #: {{orderNumber}}</Underline>

<Bold>Total: {{total}}</Bold>

Thank you, {{customerName}}!`

	// Create variables map
	variables := map[string]string{
		"storeName":   "My Coffee Shop",
		"date":        "2025-08-07",
		"orderNumber": "12345",
		"total":       "$15.50",
		"customerName": "John Doe",
	}

	// Render the template with variables
	data, err := template.RenderToBytesWithVariables(templateContent, variables)
	if err != nil {
		log.Fatalf("Failed to render template with variables: %v", err)
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

	fmt.Println("\nTemplate with variables rendered successfully!")

	// Test template file rendering
	fmt.Println("\nTesting template file rendering...")
	
	fileVariables := map[string]string{
		"storeName":    "Coffee & More",
		"date":         "2025-08-07",
		"time":         "14:30:25",
		"orderNumber":  "67890",
		"items":        "Coffee        $3.50\nSandwich      $8.75",
		"subtotal":     "$12.25",
		"tax":          "$0.98",
		"total":        "$13.23",
		"customerName": "Jane Smith",
	}

	fileData, err := template.RenderTemplateFileWithVariables("templates/receipt.tmpl", fileVariables)
	if err != nil {
		log.Printf("Failed to render template file (this is expected if file doesn't exist): %v", err)
	} else {
		fmt.Printf("Template file rendered successfully! Generated %d bytes\n", len(fileData))
		
		// Print readable representation
		fmt.Printf("\nTemplate file readable representation:\n")
		for _, b := range fileData {
			if b >= 32 && b <= 126 {
				fmt.Printf("%c", b)
			} else {
				fmt.Printf("[%02X]", b)
			}
		}
		fmt.Println()
	}
}