package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func ReadImageToStruct(imagePath string) (image.Image, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file %s: %w", filepath.Base(imagePath), err)
	}
	defer imageFile.Close()

	parsedImage, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file %s: %w", filepath.Base(imagePath), err)
	}

	return parsedImage, nil
}

func SavePNG(p string, img image.Image) error {
	file, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to save image %q: %s\n", p, err)
	}

	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("failed to encode image %q: %s\n", p, err)
	}

	return nil
}

func SaveJPEG(p string, img image.Image, o *jpeg.Options) error {
	file, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to save image %q: %s\n", p, err)
	}

	defer file.Close()

	err = jpeg.Encode(file, img, o)
	if err != nil {
		return fmt.Errorf("failed to encode jpeg %q: %s\n", p, err)
	}

	return nil
}

func IsBlackWhiteImage(tile image.Image) bool {
	for i := 0; i < tile.Bounds().Dx(); i++ {
		for j := 0; j < tile.Bounds().Dy(); j++ {
			c := tile.At(i, j)
			r, g, b, _ := c.RGBA()

			isBlack := r == 0 && g == 0 && b == 0
			isWhite := r == 65535 && g == 65535 && b == 65535

			if !isBlack && !isWhite {
				return false
			}
		}
	}
	return true
}

func IsFullyBlackOrWhiteImage(img image.Image) bool {
	meanValue := uint32(0)

	for i := 0; i < img.Bounds().Dx(); i++ {
		for j := 0; j < img.Bounds().Dy(); j++ {
			r, g, b, _ := img.At(i, j).RGBA()

			meanValue += (r + g + b) / 3
		}
	}

	meanValue /= uint32(img.Bounds().Dx() * img.Bounds().Dy())

	return meanValue == 0 || meanValue == 65535
}
