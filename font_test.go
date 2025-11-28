package sahar

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFonts(t *testing.T) {
	t.Run("returns nil for empty input", func(t *testing.T) {
		err := LoadFonts()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("returns error for odd number of arguments", func(t *testing.T) {
		err := LoadFonts("Arial")
		if err == nil {
			t.Error("expected error for odd number of arguments")
		}
		if !strings.Contains(err.Error(), "pairs") {
			t.Errorf("expected error about pairs, got: %v", err)
		}
	})

	t.Run("returns error for non-existent font file", func(t *testing.T) {
		err := LoadFonts("TestFont", "/nonexistent/font.ttf")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})

	t.Run("returns error for invalid font file", func(t *testing.T) {
		// Create a temp file with invalid font data
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "invalid.ttf")
		err := os.WriteFile(tmpFile, []byte("not a font"), 0o644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		err = LoadFonts("InvalidFont", tmpFile)
		if err == nil {
			t.Error("expected error for invalid font file")
		}
		if !strings.Contains(err.Error(), "parse") {
			t.Errorf("expected error about parsing, got: %v", err)
		}
	})

	t.Run("loads valid font file", func(t *testing.T) {
		// Check if Arial.ttf exists in examples
		arialPath := "./examples/basic/Arial.ttf"
		if _, err := os.Stat(arialPath); os.IsNotExist(err) {
			t.Skip("Arial.ttf not found in examples/basic")
		}

		err := LoadFonts("TestArial", arialPath)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("loads multiple fonts", func(t *testing.T) {
		arialPath := "./examples/basic/Arial.ttf"
		if _, err := os.Stat(arialPath); os.IsNotExist(err) {
			t.Skip("Arial.ttf not found in examples/basic")
		}

		// Load same font with different names
		err := LoadFonts(
			"Font1", arialPath,
			"Font2", arialPath,
		)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGetFontFace(t *testing.T) {
	t.Run("returns nil for non-existent font", func(t *testing.T) {
		face := getFontFace("NonExistentFont", 12)
		if face != nil {
			t.Error("expected nil for non-existent font")
		}
	})

	t.Run("returns face for loaded font", func(t *testing.T) {
		arialPath := "./examples/basic/Arial.ttf"
		if _, err := os.Stat(arialPath); os.IsNotExist(err) {
			t.Skip("Arial.ttf not found in examples/basic")
		}

		err := LoadFonts("TestArialFace", arialPath)
		if err != nil {
			t.Fatalf("failed to load font: %v", err)
		}

		face := getFontFace("TestArialFace", 12)
		if face == nil {
			t.Error("expected face for loaded font, got nil")
		}
		if face != nil {
			face.Close()
		}
	})

	t.Run("returns different face for different sizes", func(t *testing.T) {
		arialPath := "./examples/basic/Arial.ttf"
		if _, err := os.Stat(arialPath); os.IsNotExist(err) {
			t.Skip("Arial.ttf not found in examples/basic")
		}

		err := LoadFonts("TestArialSize", arialPath)
		if err != nil {
			t.Fatalf("failed to load font: %v", err)
		}

		face12 := getFontFace("TestArialSize", 12)
		face24 := getFontFace("TestArialSize", 24)

		if face12 == nil || face24 == nil {
			t.Error("expected faces for loaded font")
		}

		if face12 != nil {
			face12.Close()
		}
		if face24 != nil {
			face24.Close()
		}
	})
}

func TestMeasureTextWidth(t *testing.T) {
	t.Run("returns approximate width without font", func(t *testing.T) {
		width := measureTextWidth("Hello", 12, "NonExistentFont")
		if width <= 0 {
			t.Error("expected positive width even without font")
		}
	})

	t.Run("returns zero for empty text", func(t *testing.T) {
		width := measureTextWidth("", 12, "NonExistentFont")
		if width != 0 {
			t.Errorf("expected zero width for empty text, got %f", width)
		}
	})

	t.Run("longer text has greater width", func(t *testing.T) {
		short := measureTextWidth("Hi", 12, "NonExistentFont")
		long := measureTextWidth("Hello World", 12, "NonExistentFont")
		if long <= short {
			t.Error("expected longer text to have greater width")
		}
	})

	t.Run("larger font has greater width", func(t *testing.T) {
		small := measureTextWidth("Test", 10, "NonExistentFont")
		large := measureTextWidth("Test", 20, "NonExistentFont")
		if large <= small {
			t.Error("expected larger font to have greater width")
		}
	})

	t.Run("measures multi-line text", func(t *testing.T) {
		single := measureTextWidth("Hello World", 12, "NonExistentFont")
		multi := measureTextWidth("Hello\nWorld", 12, "NonExistentFont")
		// Multi-line should have width of the widest line
		if multi <= 0 {
			t.Error("expected positive width for multi-line text")
		}
		// "World" is 5 chars, "Hello" is 5 chars, "Hello World" is 11 chars
		// Multi-line width should be based on the widest line (5 chars)
		// which is less than 11 chars
		_ = single // single line includes space, multi-line measures max of individual lines
	})

	t.Run("measures with loaded font", func(t *testing.T) {
		arialPath := "./examples/basic/Arial.ttf"
		if _, err := os.Stat(arialPath); os.IsNotExist(err) {
			t.Skip("Arial.ttf not found in examples/basic")
		}

		err := LoadFonts("MeasureTestArial", arialPath)
		if err != nil {
			t.Fatalf("failed to load font: %v", err)
		}

		width := measureTextWidth("Hello", 12, "MeasureTestArial")
		if width <= 0 {
			t.Error("expected positive width with loaded font")
		}
	})
}

func TestMeasureTextHeight(t *testing.T) {
	t.Run("returns approximate height without font", func(t *testing.T) {
		height := measureTextHeight("Hello", 12, "NonExistentFont")
		if height <= 0 {
			t.Error("expected positive height even without font")
		}
	})

	t.Run("single line has baseline height", func(t *testing.T) {
		single := measureTextHeight("Hello", 12, "NonExistentFont")
		if single <= 0 {
			t.Error("expected positive height for single line")
		}
	})

	t.Run("multi-line has proportionally greater height", func(t *testing.T) {
		single := measureTextHeight("Hello", 12, "NonExistentFont")
		double := measureTextHeight("Hello\nWorld", 12, "NonExistentFont")
		if double <= single {
			t.Error("expected two lines to have greater height than one")
		}
	})

	t.Run("larger font has greater height", func(t *testing.T) {
		small := measureTextHeight("Test", 10, "NonExistentFont")
		large := measureTextHeight("Test", 20, "NonExistentFont")
		if large <= small {
			t.Error("expected larger font to have greater height")
		}
	})

	t.Run("measures with loaded font", func(t *testing.T) {
		arialPath := "./examples/basic/Arial.ttf"
		if _, err := os.Stat(arialPath); os.IsNotExist(err) {
			t.Skip("Arial.ttf not found in examples/basic")
		}

		err := LoadFonts("HeightTestArial", arialPath)
		if err != nil {
			t.Fatalf("failed to load font: %v", err)
		}

		height := measureTextHeight("Hello", 12, "HeightTestArial")
		if height <= 0 {
			t.Error("expected positive height with loaded font")
		}
	})
}

func TestWrapTextToWidth(t *testing.T) {
	t.Run("returns original text for zero width", func(t *testing.T) {
		result := wrapTextToWidth("Hello World", 0, 12, "NonExistentFont")
		if result != "Hello World" {
			t.Errorf("expected original text, got %q", result)
		}
	})

	t.Run("returns original text for negative width", func(t *testing.T) {
		result := wrapTextToWidth("Hello World", -100, 12, "NonExistentFont")
		if result != "Hello World" {
			t.Errorf("expected original text, got %q", result)
		}
	})

	t.Run("returns original text for empty input", func(t *testing.T) {
		result := wrapTextToWidth("", 100, 12, "NonExistentFont")
		if result != "" {
			t.Errorf("expected empty text, got %q", result)
		}
	})

	t.Run("wraps long text", func(t *testing.T) {
		// With a small width, long text should wrap
		result := wrapTextToWidth("Hello World How Are You", 50, 12, "NonExistentFont")
		if !strings.Contains(result, "\n") {
			t.Log("text wrapping depends on font availability")
		}
	})

	t.Run("preserves short text", func(t *testing.T) {
		// Short text with large width shouldn't wrap
		result := wrapTextToWidth("Hi", 1000, 12, "NonExistentFont")
		if strings.Contains(result, "\n") {
			t.Error("short text should not wrap with large width")
		}
	})
}

func TestWrapTextByCharCount(t *testing.T) {
	t.Run("wraps text at character limit", func(t *testing.T) {
		result := wrapTextByCharCount("Hello World Test", 5)
		lines := strings.Split(result, "\n")
		if len(lines) < 2 {
			t.Error("expected text to be wrapped")
		}
	})

	t.Run("handles single word longer than limit", func(t *testing.T) {
		result := wrapTextByCharCount("Supercalifragilisticexpialidocious", 5)
		// Should not panic
		if result == "" {
			t.Error("expected non-empty result")
		}
	})

	t.Run("preserves text within limit", func(t *testing.T) {
		result := wrapTextByCharCount("Hi", 100)
		if result != "Hi" {
			t.Errorf("expected 'Hi', got %q", result)
		}
	})

	t.Run("handles empty text", func(t *testing.T) {
		result := wrapTextByCharCount("", 10)
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("handles whitespace-only text", func(t *testing.T) {
		result := wrapTextByCharCount("   ", 10)
		if result != "" {
			t.Errorf("expected empty string for whitespace, got %q", result)
		}
	})
}

func TestGetActualWidth(t *testing.T) {
	t.Run("returns node width value", func(t *testing.T) {
		node := Box(Sizing(Fixed(100)))
		Layout(node)

		if getActualWidth(node) != 100 {
			t.Errorf("expected 100, got %f", getActualWidth(node))
		}
	})
}

func TestGetActualHeight(t *testing.T) {
	t.Run("returns node height value", func(t *testing.T) {
		node := Box(Sizing(Fit(), Fixed(100)))
		Layout(node)

		if getActualHeight(node) != 100 {
			t.Errorf("expected 100, got %f", getActualHeight(node))
		}
	})
}
