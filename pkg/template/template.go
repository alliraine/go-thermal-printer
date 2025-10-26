package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/unicode/norm"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
)

const (
	defaultCharactersPerLine = 32
	defaultCharacterCodePage = escpos.CharacterCodePageDefault
)

var textReplacer = strings.NewReplacer(
	"\u2018", "'",
	"\u2019", "'",
	"\u201C", "\"",
	"\u201D", "\"",
	"\u2013", "-",
	"\u2014", "-",
	"\u2026", "...",
	"\u2022", "*",
	"\u00B7", "*",
	"\u2122", "TM",
	"\u00AE", "(R)",
	"\u00A9", "(C)",
	"\u00F1", "n",
	"\u00D1", "N",
	"\u00E9", "e",
	"\u00E8", "e",
	"\u00EA", "e",
	"\u00EB", "e",
	"\u00E1", "a",
	"\u00E0", "a",
	"\u00E2", "a",
	"\u00E3", "a",
	"\u00F3", "o",
	"\u00F2", "o",
	"\u00F4", "o",
	"\u00F5", "o",
	"\u00FA", "u",
	"\u00F9", "u",
	"\u00FB", "u",
	"\u00FC", "u",
	"\u00FF", "y",
	"\u00FD", "y",
	"\u266A", "*",
	"\u266B", "*",
	"\u2605", "*",
	"\u2606", "*",
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
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x45\x01%s\x1B\x45\x00", text)
		},
		"underline": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x2D\x01%s\x1B\x2D\x00", text)
		},
		"italic": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x34\x01%s\x1B\x34\x00", text)
		},
		"italics": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x34\x01%s\x1B\x34\x00", text)
		},
		"fontb": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x4D\x01%s\x1B\x4D\x00", text)
		},
		"center": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x61\x01%s\x1B\x61\x00", text)
		},
		"left": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x61\x00%s", text)
		},
		"right": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x61\x02%s\x1B\x61\x00", text)
		},
		"rotate90": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1B\x56\x01%s\x1B\x56\x00", text)
		},
		"wrap": func(text string, maxWidth ...int) string {
			width := defaultCharactersPerLine
			if len(maxWidth) > 0 && maxWidth[0] > 0 {
				width = maxWidth[0]
			}
			return wrapText(text, width)
		},
		"fontOptions": func(args ...any) (string, error) {
			return buildFontOptions(args...)
		},
		"image": func(data any, maxWidth ...int) (string, error) {
			var encoded string
			switch v := data.(type) {
			case string:
				encoded = v
			case []byte:
				encoded = string(v)
			default:
				return "", fmt.Errorf("image expects string or []byte, got %T", data)
			}
			width := 0
			if len(maxWidth) > 0 {
				width = maxWidth[0]
			}
			bytes, err := escpos.EncodeImageToRasterBytes(encoded, width)
			if err != nil {
				return "", fmt.Errorf("image render failed: %w", err)
			}
			bytes = append(bytes, 0x1B, 0x74, byte(defaultCharacterCodePage))
			return string(bytes), nil
		},
		"qr": func(data string, args ...any) (string, error) {
			qrBytes, err := buildQRCode(data, args...)
			if err != nil {
				return "", err
			}
			qrBytes = append(qrBytes, 0x1B, 0x74, byte(defaultCharacterCodePage))
			return string(qrBytes), nil
		},
		"align": func(position string) (string, error) {
			switch strings.ToLower(strings.TrimSpace(position)) {
			case "left":
				return string([]byte{0x1B, 0x61, 0x00}), nil
			case "center":
				return string([]byte{0x1B, 0x61, 0x01}), nil
			case "right":
				return string([]byte{0x1B, 0x61, 0x02}), nil
			default:
				return "", fmt.Errorf("align expects left, center, or right; got %s", position)
			}
		},
		"feed": func(lines int) (string, error) {
			if lines < 0 || lines > 255 {
				return "", fmt.Errorf("feed expects 0-255 lines; got %d", lines)
			}
			return string([]byte{0x1B, 0x64, byte(lines)}), nil
		},
		"feedDots": func(dots int) (string, error) {
			if dots < 0 || dots > 255 {
				return "", fmt.Errorf("feedDots expects 0-255 dots; got %d", dots)
			}
			return string([]byte{0x1B, 0x4A, byte(dots)}), nil
		},
		"icon": func(name string, opts ...any) (string, error) {
			return iconTemplateFunc(name, opts...)
		},
		"cut": func(mode ...string) (string, error) {
			cutMode := byte(0x00)
			if len(mode) > 0 {
				switch strings.ToLower(strings.TrimSpace(mode[0])) {
				case "partial":
					cutMode = 0x01
				case "full", "":
					cutMode = 0x00
				default:
					return "", fmt.Errorf("cut expects 'full' or 'partial'; got %s", mode[0])
				}
			}
			return string([]byte{0x1D, 0x56, cutMode}), nil
		},
		"doubleWidth": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1D\x21\x10%s\x1D\x21\x00", text)
		},
		"doubleHeight": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1D\x21\x01%s\x1D\x21\x00", text)
		},
		"doubleSize": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1D\x21\x11%s\x1D\x21\x00", text)
		},
		"invert": func(text string) string {
			text = encodeToCodePage(text)
			return fmt.Sprintf("\x1D\x42\x01%s\x1D\x42\x00", text)
		},
		"lineSpacing": func(dots int) (string, error) {
			if dots < 0 || dots > 255 {
				return "", fmt.Errorf("lineSpacing expects 0-255 dots; got %d", dots)
			}
			return string([]byte{0x1B, 0x33, byte(dots)}), nil
		},
		"reset": func() string {
			return string([]byte{0x1B, 0x40, 0x1B, 0x74, byte(defaultCharacterCodePage)})
		},
	}
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return encodeToCodePage(text)
	}

	lines := strings.Split(text, "\n")
	var builder strings.Builder

	for i, line := range lines {
		wrappedLines := wrapLine(line, width)
		for j, wrapped := range wrappedLines {
			builder.WriteString(wrapped)
			if j < len(wrappedLines)-1 {
				builder.WriteByte('\n')
			}
		}
		if i < len(lines)-1 {
			builder.WriteByte('\n')
		}
	}

	return encodeToCodePage(builder.String())
}

func wrapLine(line string, width int) []string {
	if width <= 0 {
		return []string{line}
	}

	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	var (
		result     []string
		current    strings.Builder
		currentLen int
	)

	appendCurrent := func() {
		if currentLen == 0 {
			return
		}
		result = append(result, current.String())
		current.Reset()
		currentLen = 0
	}

	for _, word := range words {
		wordLen := utf8.RuneCountInString(word)

		if currentLen == 0 {
			current.WriteString(word)
			currentLen = wordLen
			if wordLen > width {
				appendCurrent()
			}
			continue
		}

		if currentLen+1+wordLen <= width {
			current.WriteByte(' ')
			current.WriteString(word)
			currentLen += 1 + wordLen
			continue
		}

		appendCurrent()

		if wordLen > width {
			result = append(result, word)
			continue
		}

		current.WriteString(word)
		currentLen = wordLen
	}

	appendCurrent()

	if len(result) == 0 {
		return []string{""}
	}

	return result
}

func reverseForRotate(text string) string {
	if text == "" {
		return ""
	}

	bytes := []byte(text)
	result := make([]byte, 0, len(bytes))
	start := 0

	for i := 0; i <= len(bytes); i++ {
		if i == len(bytes) || bytes[i] == '\n' {
			line := bytes[start:i]
			var suffix byte
			if len(line) > 0 && line[len(line)-1] == '\r' {
				suffix = '\r'
				line = line[:len(line)-1]
			}

			reverseBytes(line)
			result = append(result, line...)
			if suffix != 0 {
				result = append(result, suffix)
			}
			if i < len(bytes) {
				result = append(result, '\n')
			}
			start = i + 1
		}
	}

	return string(result)
}

func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func asciiFallback(r rune) string {
	switch r {
	case '\r', '\n', '\t':
		return string(r)
	}

	decomposed := norm.NFD.String(string(r))
	var builder strings.Builder
	for _, dr := range decomposed {
		if unicode.Is(unicode.Mn, dr) {
			continue
		}
		if dr < 0x80 {
			builder.WriteRune(dr)
		}
	}
	if builder.Len() > 0 {
		return builder.String()
	}
	return "?"
}

func sumRuneCounts(m map[rune]int) int {
	total := 0
	for _, count := range m {
		total += count
	}
	return total
}

func formatReplacementSummary(order []rune, counts map[rune]int) string {
	if len(order) == 0 {
		return ""
	}

	var builder strings.Builder
	max := len(order)
	if max > 5 {
		max = 5
	}
	for i := 0; i < max; i++ {
		if i > 0 {
			builder.WriteString(", ")
		}
		r := order[i]
		builder.WriteString(fmt.Sprintf("%q→%q(x%d)", r, asciiFallback(r), counts[r]))
	}
	if len(order) > max {
		builder.WriteString(", …")
	}
	return builder.String()
}

func encodeToCodePage(text string) string {
	if text == "" {
		return ""
	}

	if !utf8.ValidString(text) {
		log.Printf("template encode: invalid UTF-8 encountered (len=%d)", len(text))
		return text
	}

	normalized := textReplacer.Replace(text)
	encoder := charmap.CodePage437.NewEncoder()
	var buf bytes.Buffer
	buf.Grow(len(normalized))
	replacements := make(map[rune]int)
	replacementOrder := make([]rune, 0, 8)

	var scratch [4]byte

	for _, r := range normalized {
		if r == '\r' || r == '\n' || r == '\t' {
			buf.WriteRune(r)
			continue
		}

		runeBytes := []byte(string(r))
		encoder.Reset()
		nDst, _, err := encoder.Transform(scratch[:], runeBytes, true)
		if err != nil {
			replacements[r]++
			if replacements[r] == 1 {
				replacementOrder = append(replacementOrder, r)
			}
			buf.WriteString(asciiFallback(r))
			continue
		}
		buf.Write(scratch[:nDst])
	}

	if len(replacements) > 0 {
		log.Printf("template encode: replaced %d runes (%d unique). %s",
			sumRuneCounts(replacements), len(replacements), formatReplacementSummary(replacementOrder, replacements))
	}

	return buf.String()
}

func buildFontOptions(args ...any) (string, error) {
	if len(args)%2 != 0 {
		return "", fmt.Errorf("fontOptions expects key/value pairs")
	}

	type fontSize struct {
		width     int
		height    int
		widthSet  bool
		heightSet bool
	}

	settings := fontSize{width: 1, height: 1}
	var builder strings.Builder

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			return "", fmt.Errorf("fontOptions expects string keys, got %T", args[i])
		}

		value := args[i+1]
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "font":
			fontValue, err := parseFont(value)
			if err != nil {
				return "", err
			}
			builder.Write([]byte{0x1B, 0x4D, fontValue})

		case "width", "charwidth", "scalewidth":
			width, err := toInt(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions width: %w", err)
			}
			if width < 1 || width > 8 {
				return "", fmt.Errorf("fontOptions width must be between 1 and 8; got %d", width)
			}
			settings.width = width
			settings.widthSet = true

		case "height", "charheight", "scaleheight":
			height, err := toInt(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions height: %w", err)
			}
			if height < 1 || height > 8 {
				return "", fmt.Errorf("fontOptions height must be between 1 and 8; got %d", height)
			}
			settings.height = height
			settings.heightSet = true

		case "linespacing":
			spacing, err := toInt(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions lineSpacing: %w", err)
			}
			if spacing < 0 || spacing > 255 {
				return "", fmt.Errorf("fontOptions lineSpacing must be between 0 and 255; got %d", spacing)
			}
			builder.Write([]byte{0x1B, 0x33, byte(spacing)})

		case "charspacing", "characterspacing", "spacing":
			spacing, err := toInt(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions charSpacing: %w", err)
			}
			if spacing < 0 || spacing > 255 {
				return "", fmt.Errorf("fontOptions charSpacing must be between 0 and 255; got %d", spacing)
			}
			builder.Write([]byte{0x1B, 0x20, byte(spacing)})

		case "bold", "emphasized":
			enabled, err := toBool(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions bold: %w", err)
			}
			var flag byte
			if enabled {
				flag = 0x01
			}
			builder.Write([]byte{0x1B, 0x45, flag})

		case "underline":
			under, err := toBoolOrInt(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions underline: %w", err)
			}
			builder.Write([]byte{0x1B, 0x2D, under})

		case "underlinelevel":
			level, err := toInt(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions underlineLevel: %w", err)
			}
			if level < 0 || level > 2 {
				return "", fmt.Errorf("fontOptions underlineLevel must be between 0 and 2; got %d", level)
			}
			builder.Write([]byte{0x1B, 0x2D, byte(level)})

		case "invert", "reverse", "negative":
			enabled, err := toBool(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions invert: %w", err)
			}
			flag := byte(0x00)
			if enabled {
				flag = 0x01
			}
			builder.Write([]byte{0x1D, 0x42, flag})

		case "doublestrike":
			enabled, err := toBool(value)
			if err != nil {
				return "", fmt.Errorf("fontOptions doubleStrike: %w", err)
			}
			flag := byte(0x00)
			if enabled {
				flag = 0x01
			}
			builder.Write([]byte{0x1B, 0x47, flag})

		default:
			return "", fmt.Errorf("fontOptions: unknown option %q", key)
		}
	}

	if settings.widthSet || settings.heightSet {
		width := settings.width
		height := settings.height
		charSize := byte(((width - 1) << 4) | (height - 1))
		builder.Write([]byte{0x1D, 0x21, charSize})
	}

	return builder.String(), nil
}

func buildQRCode(data string, args ...any) ([]byte, error) {
	if data == "" {
		return nil, fmt.Errorf("qr: data is required")
	}

	options := qrOptions{
		scale:         8,
		errorLevel:    qrcode.Medium,
		disableBorder: false,
		maxWidth:      0,
	}

	if len(args) > 0 {
		if len(args)%2 != 0 {
			return nil, fmt.Errorf("qr expects key/value option pairs")
		}
		for i := 0; i < len(args); i += 2 {
			key, ok := args[i].(string)
			if !ok {
				return nil, fmt.Errorf("qr option keys must be strings, got %T", args[i])
			}
			value := args[i+1]
			switch strings.ToLower(strings.TrimSpace(key)) {
			case "size", "module", "modulesize", "scale":
				scale, err := toInt(value)
				if err != nil {
					return nil, fmt.Errorf("qr size: %w", err)
				}
				if scale < 1 {
					return nil, fmt.Errorf("qr size must be >= 1; got %d", scale)
				}
				options.scale = scale

			case "error", "errorlevel", "correction", "ecc":
				errLevel, err := parseQRErrorLevel(value)
				if err != nil {
					return nil, err
				}
				options.errorLevel = errLevel

			case "border", "margin", "quietzone":
				border, err := toInt(value)
				if err != nil {
					return nil, fmt.Errorf("qr border: %w", err)
				}
				options.disableBorder = border <= 0

			case "maxwidth", "width":
				maxWidth, err := toInt(value)
				if err != nil {
					return nil, fmt.Errorf("qr maxWidth: %w", err)
				}
				if maxWidth < 0 {
					return nil, fmt.Errorf("qr maxWidth must be >= 0; got %d", maxWidth)
				}
				options.maxWidth = maxWidth

			default:
				return nil, fmt.Errorf("qr: unknown option %q", key)
			}
		}
	}

	qrCode, err := qrcode.New(data, options.errorLevel)
	if err != nil {
		return nil, fmt.Errorf("qr: failed to encode data: %w", err)
	}
	qrCode.DisableBorder = options.disableBorder

	if options.scale <= 0 {
		options.scale = 8
	}

	pngBytes, err := qrCode.PNG(-options.scale)
	if err != nil {
		return nil, fmt.Errorf("qr: failed to render image: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(pngBytes)
	imgBytes, err := escpos.EncodeImageToRasterBytes(encoded, options.maxWidth)
	if err != nil {
		return nil, fmt.Errorf("qr: failed to convert image: %w", err)
	}

	return imgBytes, nil
}

type qrOptions struct {
	scale         int
	errorLevel    qrcode.RecoveryLevel
	disableBorder bool
	maxWidth      int
}

func parseFont(value any) (byte, error) {
	switch v := value.(type) {
	case string:
		switch strings.ToUpper(strings.TrimSpace(v)) {
		case "A":
			return 0x00, nil
		case "B":
			return 0x01, nil
		case "C":
			return 0x02, nil
		default:
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return 0x00, fmt.Errorf("fontOptions font: unsupported value %q", v)
			}
			return parseFont(parsed)
		}
	case int:
		if v < 0 || v > 2 {
			return 0x00, fmt.Errorf("fontOptions font must be 0, 1, or 2; got %d", v)
		}
		return byte(v), nil
	case int64:
		return parseFont(int(v))
	case float64:
		if v != float64(int(v)) {
			return 0x00, fmt.Errorf("fontOptions font requires an integer, got %v", v)
		}
		return parseFont(int(v))
	default:
		return 0x00, fmt.Errorf("fontOptions font: value type %T not supported", value)
	}
}

func parseQRErrorLevel(value any) (qrcode.RecoveryLevel, error) {
	switch v := value.(type) {
	case string:
		switch strings.ToUpper(strings.TrimSpace(v)) {
		case "L", "0":
			return qrcode.Low, nil
		case "M", "1":
			return qrcode.Medium, nil
		case "Q", "2":
			return qrcode.High, nil
		case "H", "3":
			return qrcode.Highest, nil
		default:
			return qrcode.Medium, fmt.Errorf("qr errorLevel: unsupported value %q", v)
		}
	case int:
		if v < 0 || v > 3 {
			return qrcode.Medium, fmt.Errorf("qr errorLevel must be 0-3; got %d", v)
		}
		return []qrcode.RecoveryLevel{qrcode.Low, qrcode.Medium, qrcode.High, qrcode.Highest}[v], nil
	case int64:
		return parseQRErrorLevel(int(v))
	case float64:
		if v != float64(int(v)) {
			return qrcode.Medium, fmt.Errorf("qr errorLevel must be integer, got %v", v)
		}
		return parseQRErrorLevel(int(v))
	default:
		return qrcode.Medium, fmt.Errorf("qr errorLevel: unsupported type %T", value)
	}
}

func toInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		if float32(int(v)) != v {
			return 0, fmt.Errorf("value %v is not an integer", v)
		}
		return int(v), nil
	case float64:
		if float64(int(v)) != v {
			return 0, fmt.Errorf("value %v is not an integer", v)
		}
		return int(v), nil
	case string:
		trimmed := strings.TrimSpace(v)
		parsed, err := strconv.Atoi(trimmed)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

func toBool(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "t", "1", "yes", "y", "on":
			return true, nil
		case "false", "f", "0", "no", "n", "off":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean value %q", v)
		}
	case int:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case float64:
		if float64(int(v)) != v {
			return false, fmt.Errorf("boolean value must be integer, got %v", v)
		}
		return int(v) != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

func toBoolOrInt(value any) (byte, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return 0x01, nil
		}
		return 0x00, nil
	case int:
		if v < 0 || v > 2 {
			return 0, fmt.Errorf("value must be between 0 and 2; got %d", v)
		}
		return byte(v), nil
	case int64:
		return toBoolOrInt(int(v))
	case float64:
		if float64(int(v)) != v {
			return 0, fmt.Errorf("value must be integer, got %v", v)
		}
		return toBoolOrInt(int(v))
	case string:
		trimmed := strings.TrimSpace(v)
		switch strings.ToLower(trimmed) {
		case "true", "t", "yes", "y", "on":
			return 0x01, nil
		case "false", "f", "no", "n", "off":
			return 0x00, nil
		default:
			parsed, err := strconv.Atoi(trimmed)
			if err != nil {
				return 0, fmt.Errorf("cannot convert %q to underline value", v)
			}
			return toBoolOrInt(parsed)
		}
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
}

// Render processes the template and returns the ESCPOS byte array
func (r *Renderer) Render(template *Template, data any) ([]byte, error) {
	r.buffer.Reset()

	r.escpos.Initialize()
	r.escpos.SelectCharacterCodePage(defaultCharacterCodePage)

	// Execute template to a temporary buffer
	tempBuffer := &bytes.Buffer{}
	if err := template.tmpl.Execute(tempBuffer, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	r.escpos.Write(tempBuffer.Bytes())

	r.escpos.PrintAndFeedPaperNLines(9)
	r.escpos.FullCut()

	return r.buffer.Bytes(), nil
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
