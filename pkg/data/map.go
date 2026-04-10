package data

import (
	"image"
	"image/draw"
	"os"
)

// LoadMap opens and decodes an image from path and returns it as *image.NRGBA.
// It uses image.Decode to support multiple image formats.
func LoadMap(path string) (*image.NRGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	nr := image.NewNRGBA(img.Bounds())
	draw.Draw(nr, nr.Bounds(), img, image.Point{}, draw.Over)
	return nr, nil
}
