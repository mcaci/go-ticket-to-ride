package render

import (
	"bytes"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"testing"
)

func TestMapWriters_WithFrames(t *testing.T) {
	// create a small base image
	nr := image.NewNRGBA(image.Rect(0, 0, 16, 16))
	draw.Draw(nr, nr.Bounds(), &image.Uniform{color.RGBA{R: 10, G: 20, B: 30, A: 255}}, image.Point{}, draw.Src)

	// create one paletted frame from the base image
	p := image.NewPaletted(nr.Bounds(), palette.Plan9)
	draw.Draw(p, p.Bounds(), nr, image.Point{}, draw.Over)
	frames := []*image.Paletted{p}

	var bufImg, bufGif bytes.Buffer
	if err := MapWriters(nr, frames, &bufImg, &bufGif, 5); err != nil {
		t.Fatalf("MapWriters failed: %v", err)
	}
	if bufImg.Len() == 0 {
		t.Fatalf("expected non-empty JPEG output")
	}
	if bufGif.Len() == 0 {
		t.Fatalf("expected non-empty GIF output")
	}
}

func TestMapWriters_NoFrames(t *testing.T) {
	// create a small base image
	nr := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	draw.Draw(nr, nr.Bounds(), &image.Uniform{color.RGBA{R: 100, G: 50, B: 25, A: 255}}, image.Point{}, draw.Src)

	var bufImg, bufGif bytes.Buffer
	if err := MapWriters(nr, nil, &bufImg, &bufGif, 5); err != nil {
		t.Fatalf("MapWriters failed: %v", err)
	}
	if bufImg.Len() == 0 {
		t.Fatalf("expected non-empty JPEG output")
	}
	if bufGif.Len() != 0 {
		t.Fatalf("expected empty GIF output when frames are nil")
	}
}
