package template

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/exp/shiny/iconvg"

	"github.com/jonasclaes/go-thermal-printer/pkg/escpos"
)

const (
	defaultIconSizeDots = 96
	defaultIconFeedLines = 0
)

type iconOptions struct {
	width int
	feed  int
}

func parseIconOptions(args []any) (iconOptions, error) {
	opts := iconOptions{width: defaultIconSizeDots, feed: defaultIconFeedLines}
	switch len(args) {
	case 0:
		return opts, nil
	case 1:
		w, err := toInt(args[0])
		if err != nil {
			return opts, fmt.Errorf("icon width: %w", err)
		}
		if w <= 0 {
			return opts, fmt.Errorf("icon width must be positive, got %d", w)
		}
		opts.width = w
	case 2:
		w, err := toInt(args[0])
		if err != nil {
			return opts, fmt.Errorf("icon width: %w", err)
		}
		if w <= 0 {
			return opts, fmt.Errorf("icon width must be positive, got %d", w)
		}
		opts.width = w

		feed, err := toInt(args[1])
		if err != nil {
			return opts, fmt.Errorf("icon feed: %w", err)
		}
		opts.feed = clamp(feed, 0, 255)
	default:
		return opts, fmt.Errorf("icon expects at most width and feed options")
	}
	return opts, nil
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func iconTemplateFunc(name string, args ...any) (string, error) {
	data, canonical, err := lookupIconData(name)
	if err != nil {
		return "", err
	}

	opts, err := parseIconOptions(args)
	if err != nil {
		return "", err
	}

	img, err := rasterizeIcon(data, opts.width)
	if err != nil {
		return "", fmt.Errorf("icon %s: %w", canonical, err)
	}

	raster, err := escpos.ImageToRasterBytes(img, opts.width)
	if err != nil {
		return "", fmt.Errorf("icon %s: %w", canonical, err)
	}

	if len(raster) >= 3 {
		raster[len(raster)-1] = byte(opts.feed)
	}
	raster = append(raster, 0x1B, 0x74, byte(defaultCharacterCodePage))

	return string(raster), nil
}

func rasterizeIcon(data []byte, width int) (image.Image, error) {
	if width <= 0 {
		width = defaultIconSizeDots
	}

	meta, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil, fmt.Errorf("decode metadata: %w", err)
	}

	dx, dy := meta.ViewBox.AspectRatio()
	if dx == 0 {
		dx = 1
	}
	height := int(math.Round(float64(width) * float64(dy) / float64(dx)))
	if height <= 0 {
		height = width
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}}, image.Point{}, draw.Src)

	palette := meta.Palette
	var rast iconvg.Rasterizer
	rast.SetDstImage(dst, dst.Bounds(), draw.Src)
	if err := iconvg.Decode(&rast, data, &iconvg.DecodeOptions{Palette: &palette}); err != nil {
		return nil, fmt.Errorf("decode icon: %w", err)
	}

	return dst, nil
}
