package sahar

import (
	"fmt"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// fontCache stores loaded fonts to avoid reloading
var fontCache = make(map[string]*truetype.Font)

func LoadFonts(src ...string) error {
	if len(src) == 0 {
		return nil // No fonts to load
	}

	if len(src)%2 != 0 {
		return fmt.Errorf("src must contain pairs of names and their paths")
	}

	for i := 0; i < len(src); i += 2 {
		name := src[i]
		path := src[i+1]

		font, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load font %s from %s: %w", name, path, err)
		}

		ttfFont, err := truetype.Parse(font)
		if err != nil {
			return fmt.Errorf("failed to parse font %s from %s: %w", name, path, err)
		}

		fontCache[name] = ttfFont
	}

	return nil
}

// getFontFace returns a font.Face for the given font type and size
func getFontFace(fontType string, fontSize float64) font.Face {
	ttfFont, exists := fontCache[fontType]
	if !exists {
		return nil
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size: fontSize,
		DPI:  72, // Standard DPI
	})
}
