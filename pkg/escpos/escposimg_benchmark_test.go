package escpos

import (
	"image"
	"image/color"
	"testing"
)

func generateTestImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := uint8((x * 255) / w)
			g := uint8((y * 255) / h)
			b := uint8(((x + y) * 255) / (w + h))
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func benchmarkImageToRasterBytes(b *testing.B, w, h int) {
	img := generateTestImage(w, h)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ImageToRasterBytes(img, 384); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkImageToRasterBytesSmall(b *testing.B) {
	benchmarkImageToRasterBytes(b, 384, 256)
}

func BenchmarkImageToRasterBytesMedium(b *testing.B) {
	benchmarkImageToRasterBytes(b, 1024, 768)
}

func BenchmarkImageToRasterBytesLarge(b *testing.B) {
	benchmarkImageToRasterBytes(b, 2048, 1536)
}
