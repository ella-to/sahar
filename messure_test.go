package sahar

import "testing"

func TestMessureString(t *testing.T) {
	// Load the font
	fnt, err := loadFont("testdata/Roboto-Regular.ttf", 12)
	if err != nil {
		t.Fatal(err)
	}

	// Measure the string
	width, height := measureString("H", fnt)
	if width != 9 || height != 15 {
		t.Fatalf("invalid width or height: %d, %d", width, height)
	}
}
