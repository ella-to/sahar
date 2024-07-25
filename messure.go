package sahar

import (
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func loadFont(filename string, size float64) (font.Face, error) {
	// Read the font file
	fontBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse the font
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	// Create a font face
	fontFace, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return fontFace, nil
}

func measureString(s string, fnt font.Face) (int, int) {
	d := &font.Drawer{
		Dst:  nil, // Destination image, not needed for measurement
		Src:  nil, // Source image, not needed for measurement
		Face: fnt,
	}
	// Calculate the width
	width := d.MeasureString(s).Ceil()
	// Calculate the height
	height := (fnt.Metrics().Ascent + fnt.Metrics().Descent).Ceil()
	return width, height
}
