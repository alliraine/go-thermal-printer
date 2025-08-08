package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
)

// Template represents a thermal printer template
type Template struct {
	tmpl *template.Template
}

// NewTemplate creates a new template
func NewTemplate(content string) (*Template, error) {
	tmpl, err := template.New("thermal").Funcs(getTemplateFuncs()).Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Template{
		tmpl: tmpl,
	}, nil
}

// Renderer handles the conversion of template to ESCPOS commands
type Renderer struct {
	escpos *escpos.ESCPOS
	buffer *bytes.Buffer
}

func NewRenderer() *Renderer {
	buffer := &bytes.Buffer{}
	escpos := escpos.NewESCPOS(buffer)

	return &Renderer{
		escpos: escpos,
		buffer: buffer,
	}
}

// getTemplateFuncs returns the custom functions available in templates
func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"bold": func(text string) string {
			return fmt.Sprintf("\x1B\x45\x01%s\x1B\x45\x00", text)
		},
		"underline": func(text string) string {
			return fmt.Sprintf("\x1B\x2D\x01%s\x1B\x2D\x00", text)
		},
		"italic": func(text string) string {
			return fmt.Sprintf("\x1B\x34\x01%s\x1B\x34\x00", text)
		},
		"italics": func(text string) string {
			return fmt.Sprintf("\x1B\x34\x01%s\x1B\x34\x00", text)
		},
		"fontb": func(text string) string {
			return fmt.Sprintf("\x1B\x4D\x01%s\x1B\x4D\x00", text)
		},
	}
}

// Render processes the template and returns the ESCPOS byte array
func (r *Renderer) Render(template *Template, data any) ([]byte, error) {
	r.buffer.Reset()

	r.escpos.Initialize()
	r.escpos.SelectCharacterCodePage(escpos.CharacterCodePagePC858)

	// Execute template to a temporary buffer
	tempBuffer := &bytes.Buffer{}
	if err := template.tmpl.Execute(tempBuffer, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Process the rendered content for ESCPOS formatting
	if err := r.processRenderedContent(tempBuffer.String()); err != nil {
		return nil, err
	}

	r.escpos.PrintAndFeedPaperNLines(9)

	return r.buffer.Bytes(), nil
}

// processRenderedContent handles the rendered template content with embedded ESCPOS commands
func (r *Renderer) processRenderedContent(content string) error {
	// The template functions have already embedded ESCPOS commands
	// We just need to write the content and handle the raw ESCPOS bytes
	r.buffer.WriteString(content)
	return nil
}

// RenderToBytes is a convenience function that creates a renderer and renders the template
func RenderToBytes(templateContent string, data any) ([]byte, error) {
	template, err := NewTemplate(templateContent)
	if err != nil {
		return nil, err
	}

	renderer := NewRenderer()
	return renderer.Render(template, data)
}

// RenderTemplateFileWithVariables reads a template file and renders it with variable substitution
func RenderTemplateFileWithVariables(filePath string, variables map[string]any) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	return RenderToBytesWithVariables(string(content), variables)
}

// RenderToBytesWithVariables renders template content with variable substitution
func RenderToBytesWithVariables(templateContent string, variables map[string]any) ([]byte, error) {
	template, err := NewTemplate(templateContent)
	if err != nil {
		return nil, err
	}

	renderer := NewRenderer()
	return renderer.Render(template, variables)
}
