package sahar

import (
	"fmt"
	"os"

	"github.com/golang/freetype/truetype"
)

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

		FontCache[name] = ttfFont
	}

	return nil
}
