package escpos

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	imagedraw "image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"strings"
	"unicode"

	xdraw "golang.org/x/image/draw"
)

// EncodeImageToRasterBytes decodes a base64 image, scales to maxWidthDots (58mm ≈ 384),
// dithers to 1-bit, and returns ESC/POS GS v 0 raster bytes.
func EncodeImageToRasterBytes(imgB64 string, maxWidthDots int) ([]byte, error) {
	cleaned := normalizeBase64(imgB64)
	if cleaned == "" {
		return nil, errors.New("empty image data")
	}
	raw, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	return ImageToRasterBytes(img, maxWidthDots)
}

// ImageToRasterBytes converts an image to ESC/POS raster bytes using GS v 0
// commands. The image is scaled to fit maxWidthDots (58mm ≈ 384) while
// preserving aspect ratio.
func ImageToRasterBytes(img image.Image, maxWidthDots int) ([]byte, error) {
	if img == nil {
		return nil, errors.New("nil image")
	}
	if maxWidthDots <= 0 {
		maxWidthDots = 384
	}
	alignedMaxWidth := (maxWidthDots + 7) &^ 7
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w == 0 || h == 0 {
		return nil, errors.New("empty image")
	}
	scale := float64(maxWidthDots) / float64(w)
	if scale > 1.0 {
		scale = 1.0
	}
	newW := int(math.Floor(float64(w)*scale + 0.5))
	if newW <= 0 {
		newW = 1
	}
	newW = (newW + 7) &^ 7
	if newW > alignedMaxWidth {
		alignedMaxWidth = newW
	}
	newH := int(math.Floor(float64(h)*scale + 0.5))
	if newH <= 0 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, b, xdraw.Src, nil)

	gray := image.NewGray(dst.Bounds())
	imagedraw.Draw(gray, gray.Bounds(), dst, dst.Bounds().Min, imagedraw.Src)

	// Pack bits row by row, centering image horizontally within alignedMaxWidth
	imageBytes := newW / 8
	lineBytes := alignedMaxWidth / 8
	leftPadBytes := (lineBytes - imageBytes) / 2
	rightPadBytes := lineBytes - imageBytes - leftPadBytes

	data := make([]byte, 0, lineBytes*newH)
	leftPad := make([]byte, leftPadBytes)
	rightPad := make([]byte, rightPadBytes)
	for y := 0; y < newH; y++ {
		data = append(data, leftPad...)
		for bx := 0; bx < imageBytes; bx++ {
			var b8 byte
			for bit := 0; bit < 8; bit++ {
				x := bx*8 + bit
				if gray.GrayAt(x, y).Y < 128 {
					b8 |= 1 << (7 - bit)
				}
			}
			data = append(data, b8)
		}
		data = append(data, rightPad...)
	}

	xL := byte(lineBytes & 0xFF)
	xH := byte((lineBytes >> 8) & 0xFF)
	yL := byte(newH & 0xFF)
	yH := byte((newH >> 8) & 0xFF)

	var out bytes.Buffer
	out.Write([]byte{0x1B, 0x40}) // ESC @
	out.Write([]byte{0x1D, 0x76, 0x30, 0x00, xL, xH, yL, yH})
	out.Write(data)
	// Feed extra space so the receipt can be torn cleanly
	out.Write([]byte{0x1B, 0x64, 0x05}) // ESC d n -> feed n lines

	return out.Bytes(), nil
}

func normalizeBase64(input string) string {
	s := strings.TrimSpace(input)
	if idx := strings.Index(s, ","); idx != -1 {
		prefix := s[:idx]
		if strings.Contains(strings.ToLower(prefix), "base64") {
			s = s[idx+1:]
		}
	}
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
	if rem := len(s) % 4; rem != 0 {
		s += strings.Repeat("=", 4-rem)
	}
	return s
}
