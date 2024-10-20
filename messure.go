package sahar

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var cacheFontFaceLock sync.Mutex
var cacheFontFace = make(map[string]font.Face)

func loadFont(filename string, size float64) (font.Face, error) {
	cacheFontFaceLock.Lock()
	defer cacheFontFaceLock.Unlock()

	key := fmt.Sprintf("%s-%f", filename, size)
	if fnt, ok := cacheFontFace[key]; ok {
		return fnt, nil
	}

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
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	cacheFontFace[key] = fontFace

	return fontFace, nil
}

func measureString(s string, fnt font.Face) (int, int) {
	d := &font.Drawer{
		Dst:  nil, // Destination image, not needed for measurement
		Src:  nil, // Source image, not needed for measurement
		Face: fnt,
	}

	// Measure the width.
	advance := d.MeasureString(s)

	// Get the height from the font metrics.
	metrics := fnt.Metrics()
	height := metrics.Ascent + metrics.Descent

	width := advance.Round() // Convert fixed.Point26_6 to int
	return width, height.Round()
}
