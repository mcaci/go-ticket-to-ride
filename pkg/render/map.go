package render

import (
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"os"
)

// MapWriters writes the NRGBA image to imgWriter and the frames to gifWriter (if frames present).
// gifWriter may be nil to skip GIF output. imgWriter must be non-nil.
func MapWriters(layer *image.NRGBA, frames []*image.Paletted, imgWriter, gifWriter io.Writer, frameDelay int) error {
	if layer == nil {
		return nil // nothing to render
	}
	if imgWriter == nil {
		return nil
	}
	if err := jpeg.Encode(imgWriter, layer, nil); err != nil {
		return err
	}

	if len(frames) == 0 || gifWriter == nil {
		return nil
	}
	delay := make([]int, len(frames))
	for i := range delay {
		delay[i] = frameDelay
	}
	g := gif.GIF{
		Image: frames,
		Delay: delay,
	}
	if err := gif.EncodeAll(gifWriter, &g); err != nil {
		return err
	}
	return nil
}

// Map writes the NRGBA image to a JPEG and the frames to a GIF (if frames present).
// It is a convenience wrapper around MapWriters that opens the output files.
func Map(layer *image.NRGBA, frames []*image.Paletted, outImagePath, outGifPath string, frameDelay int) error {
	if layer == nil {
		return nil
	}
	out, err := os.Create(outImagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	var gifFile *os.File
	if len(frames) > 0 {
		gifFile, err = os.Create(outGifPath)
		if err != nil {
			return err
		}
		defer gifFile.Close()
	}

	return MapWriters(layer, frames, out, gifFile, frameDelay)
}
