package template

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Template represents a thermal printer template
type Template struct {
	content string
}

// NewTemplate creates a new template
func NewTemplate(content string) *Template {
	return &Template{
		content: content,
	}
}

// Renderer handles the conversion of template to ESCPOS commands
type Renderer struct {
	buffer *bytes.Buffer
}

// NewRenderer creates a new template renderer
func NewRenderer() *Renderer {
	return &Renderer{
		buffer: &bytes.Buffer{},
	}
}

// Render processes the template and returns the ESCPOS byte array
func (r *Renderer) Render(template *Template) ([]byte, error) {
	r.buffer.Reset()

	// Initialize printer
	r.writeESCPOSCommand([]byte{0x1B, 0x40}) // Initialize

	r.writeESCPOSCommand([]byte{0x1B, 0x74, 0x0F})

	r.processContent(template.content)

	r.writeESCPOSCommand([]byte{0x1B, 0x64, byte(9)})

	return r.buffer.Bytes(), nil
}

// processContent recursively processes template content, handling nested tags
func (r *Renderer) processContent(content string) ([]byte, error) {
	for {
		// Find the next opening tag
		openTagRegex := regexp.MustCompile(`<(\w+)>`)
		openMatch := openTagRegex.FindStringSubmatch(content)

		if openMatch == nil {
			// No more tags, write remaining content
			r.buffer.WriteString(content)
			break
		}

		tagName := strings.ToLower(openMatch[1])
		openTag := openMatch[0]
		openIndex := strings.Index(content, openTag)

		// Find the corresponding closing tag
		closeTag := fmt.Sprintf("</%s>", openMatch[1])
		closeIndex := r.findMatchingCloseTag(content, openTag, closeTag, openIndex)

		if closeIndex == -1 {
			return nil, fmt.Errorf("missing closing tag for <%s>", openMatch[1])
		}

		// Write content before the tag
		r.buffer.WriteString(content[:openIndex])

		// Apply the formatting command
		if err := r.applyOpenTag(tagName); err != nil {
			return nil, err
		}

		// Process the content inside the tag (this handles nested tags)
		tagContent := content[openIndex+len(openTag) : closeIndex]
		innerRenderer := NewRenderer()
		_, err := innerRenderer.processContent(tagContent)
		if err != nil {
			return nil, err
		}
		r.buffer.Write(innerRenderer.buffer.Bytes())

		// Apply the closing command
		if err := r.applyCloseTag(tagName); err != nil {
			return nil, err
		}

		// Continue with the rest of the content
		content = content[closeIndex+len(closeTag):]
	}

	return r.buffer.Bytes(), nil
}

// findMatchingCloseTag finds the matching closing tag, accounting for nested tags
func (r *Renderer) findMatchingCloseTag(content, openTag, closeTag string, startIndex int) int {
	openCount := 1
	searchStart := startIndex + len(openTag)

	for openCount > 0 {
		nextOpen := strings.Index(content[searchStart:], openTag)
		nextClose := strings.Index(content[searchStart:], closeTag)

		if nextClose == -1 {
			return -1 // No matching close tag
		}

		// Adjust indices to be relative to the original content
		if nextOpen != -1 {
			nextOpen += searchStart
		}
		nextClose += searchStart

		if nextOpen != -1 && nextOpen < nextClose {
			// Found another opening tag before the closing tag
			openCount++
			searchStart = nextOpen + len(openTag)
		} else {
			// Found a closing tag
			openCount--
			if openCount == 0 {
				return nextClose
			}
			searchStart = nextClose + len(closeTag)
		}
	}

	return -1
}

// applyOpenTag applies the ESCPOS command for opening a tag
func (r *Renderer) applyOpenTag(tagName string) error {
	switch tagName {
	case "bold":
		// Emphasis mode on
		r.writeESCPOSCommand([]byte{0x1B, 0x45, 0x01})
	case "underline":
		// Underline mode on (1 dot thick)
		r.writeESCPOSCommand([]byte{0x1B, 0x2D, 0x01})
	case "italic", "italics":
		// Italics mode on
		r.writeESCPOSCommand([]byte{0x1B, 0x34, 0x01})
	case "fontb":
		// Select character font B
		r.writeESCPOSCommand([]byte{0x1B, 0x4D, 0x01})
	default:
		return fmt.Errorf("unsupported tag: %s", tagName)
	}
	return nil
}

// applyCloseTag applies the ESCPOS command for closing a tag
func (r *Renderer) applyCloseTag(tagName string) error {
	switch tagName {
	case "bold":
		// Emphasis mode off
		r.writeESCPOSCommand([]byte{0x1B, 0x45, 0x00})
	case "underline":
		// Underline mode off
		r.writeESCPOSCommand([]byte{0x1B, 0x2D, 0x00})
	case "italic", "italics":
		// Italics mode off
		r.writeESCPOSCommand([]byte{0x1B, 0x34, 0x00})
	case "fontb":
		// Select character font A (default)
		r.writeESCPOSCommand([]byte{0x1B, 0x4D, 0x00})
	default:
		return fmt.Errorf("unsupported tag: %s", tagName)
	}
	return nil
}

// writeESCPOSCommand writes ESCPOS command bytes to the buffer
func (r *Renderer) writeESCPOSCommand(command []byte) {
	r.buffer.Write(command)
}

// RenderToBytes is a convenience function that creates a renderer and renders the template
func RenderToBytes(templateContent string) ([]byte, error) {
	template := NewTemplate(templateContent)
	renderer := NewRenderer()
	return renderer.Render(template)
}

// RenderTemplateFileWithVariables reads a template file and renders it with variable substitution
func RenderTemplateFileWithVariables(filePath string, variables map[string]string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	return RenderToBytesWithVariables(string(content), variables)
}

// RenderToBytesWithVariables renders template content with variable substitution
func RenderToBytesWithVariables(templateContent string, variables map[string]string) ([]byte, error) {
	// Substitute variables in the template content
	// Variables are expected to be in the format {{variableName}}
	content := substituteVariables(templateContent, variables)

	template := NewTemplate(content)
	renderer := NewRenderer()
	return renderer.Render(template)
}

// substituteVariables replaces {{variableName}} with the corresponding value from the variables map
func substituteVariables(content string, variables map[string]string) string {
	// Use regex to find all variables in the format {{variableName}}
	variableRegex := regexp.MustCompile(`\{\{(\w+)\}\}`)

	return variableRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract the variable name (remove {{ and }})
		varName := match[2 : len(match)-2]

		// Look up the variable in the map
		if value, exists := variables[varName]; exists {
			return value
		}

		// If variable not found, return the original placeholder
		return match
	})
}
