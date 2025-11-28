package sahar

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name         string
		hex          string
		defaultColor string
		wantR        int
		wantG        int
		wantB        int
		wantErr      bool
	}{
		{
			name:         "valid hex with hash",
			hex:          "#FF0000",
			defaultColor: "#000000",
			wantR:        255,
			wantG:        0,
			wantB:        0,
			wantErr:      false,
		},
		{
			name:         "valid hex without hash",
			hex:          "00FF00",
			defaultColor: "#000000",
			wantR:        0,
			wantG:        255,
			wantB:        0,
			wantErr:      false,
		},
		{
			name:         "valid blue color",
			hex:          "#0000FF",
			defaultColor: "#000000",
			wantR:        0,
			wantG:        0,
			wantB:        255,
			wantErr:      false,
		},
		{
			name:         "lowercase hex",
			hex:          "#aabbcc",
			defaultColor: "#000000",
			wantR:        170,
			wantG:        187,
			wantB:        204,
			wantErr:      false,
		},
		{
			name:         "empty string uses default",
			hex:          "",
			defaultColor: "#FFFFFF",
			wantR:        255,
			wantG:        255,
			wantB:        255,
			wantErr:      false,
		},
		{
			name:         "invalid hex - too short",
			hex:          "#FFF",
			defaultColor: "#000000",
			wantErr:      true,
		},
		{
			name:         "invalid hex - too long",
			hex:          "#FFFFFFF",
			defaultColor: "#000000",
			wantErr:      true,
		},
		{
			name:         "invalid hex - non-hex characters",
			hex:          "#GGHHII",
			defaultColor: "#000000",
			wantErr:      true,
		},
		{
			name:         "mixed case hex",
			hex:          "#AaBbCc",
			defaultColor: "#000000",
			wantR:        170,
			wantG:        187,
			wantB:        204,
			wantErr:      false,
		},
		{
			name:         "black color",
			hex:          "#000000",
			defaultColor: "#FFFFFF",
			wantR:        0,
			wantG:        0,
			wantB:        0,
			wantErr:      false,
		},
		{
			name:         "white color",
			hex:          "#FFFFFF",
			defaultColor: "#000000",
			wantR:        255,
			wantG:        255,
			wantB:        255,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, err := hexToRGB(tt.hex, tt.defaultColor)
			if (err != nil) != tt.wantErr {
				t.Errorf("hexToRGB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if r != tt.wantR {
					t.Errorf("hexToRGB() r = %v, want %v", r, tt.wantR)
				}
				if g != tt.wantG {
					t.Errorf("hexToRGB() g = %v, want %v", g, tt.wantG)
				}
				if b != tt.wantB {
					t.Errorf("hexToRGB() b = %v, want %v", b, tt.wantB)
				}
			}
		})
	}
}

func TestMapFontName(t *testing.T) {
	tests := []struct {
		name     string
		fontType string
		want     string
	}{
		{"arial lowercase", "arial", "Arial"},
		{"Arial mixed case", "Arial", "Arial"},
		{"ARIAL uppercase", "ARIAL", "Arial"},
		{"helvetica", "helvetica", "Arial"},
		{"Helvetica", "Helvetica", "Arial"},
		{"times", "times", "Times"},
		{"Times New Roman", "times new roman", "Times"},
		{"courier", "courier", "Courier"},
		{"Courier New", "courier new", "Courier"},
		{"unknown font defaults to Arial", "CustomFont", "Arial"},
		{"empty string defaults to Arial", "", "Arial"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapFontName(tt.fontType); got != tt.want {
				t.Errorf("mapFontName(%q) = %v, want %v", tt.fontType, got, tt.want)
			}
		})
	}
}

func TestRenderToPDF(t *testing.T) {
	t.Run("returns error for empty nodes", func(t *testing.T) {
		var buf bytes.Buffer
		err := RenderToPDF(&buf)
		if err == nil {
			t.Error("expected error for empty nodes, got nil")
		}
		if !strings.Contains(err.Error(), "no node") {
			t.Errorf("expected error message about no nodes, got: %v", err)
		}
	})

	t.Run("renders single empty box", func(t *testing.T) {
		node := Box(Sizing(Fixed(100), Fixed(100)))
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("expected PDF output, got empty buffer")
		}
		// Check PDF magic number
		if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF")) {
			t.Error("output does not appear to be a valid PDF")
		}
	})

	t.Run("renders text node", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(200), Fixed(100)),
			Children(
				Text("Hello World", FontSize(12)),
			),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("expected PDF output, got empty buffer")
		}
	})

	t.Run("renders multiple pages", func(t *testing.T) {
		page1 := Box(Sizing(Fixed(100), Fixed(100)))
		page2 := Box(Sizing(Fixed(100), Fixed(100)))
		Layout(page1)
		Layout(page2)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, page1, page2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("expected PDF output, got empty buffer")
		}
	})

	t.Run("renders nested boxes", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(300), Fixed(200)),
			BackgroundColor("#FFFFFF"),
			Children(
				Box(
					Sizing(Fixed(100), Fixed(50)),
					BackgroundColor("#FF0000"),
				),
				Box(
					Sizing(Fixed(100), Fixed(50)),
					BackgroundColor("#00FF00"),
				),
			),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("renders text with alignment", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(200), Fixed(100)),
			Alignment(Center, Middle),
			Children(
				Text("Centered", FontSize(14)),
			),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("renders multi-line text", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(200), Fixed(100)),
			Children(
				Text("Line 1\nLine 2\nLine 3", FontSize(12)),
			),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("renders empty text node without error", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(100), Fixed(100)),
			Children(Text("")),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("handles invalid font color gracefully", func(t *testing.T) {
		node := Text("Test", FontSize(12), FontColor("invalid"))
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err == nil {
			t.Error("expected error for invalid font color")
		}
	})
}

func TestRenderToPDFWithOptions(t *testing.T) {
	t.Run("returns error for nil root", func(t *testing.T) {
		var buf bytes.Buffer
		err := RenderToPDFWithOptions(nil, &buf, DefaultPDFOptions())
		if err == nil {
			t.Error("expected error for nil root, got nil")
		}
	})

	t.Run("renders with landscape orientation", func(t *testing.T) {
		node := Box(Sizing(Fixed(100), Fixed(100)))
		Layout(node)

		opts := DefaultPDFOptions()
		opts.Landscape = true

		var buf bytes.Buffer
		err := RenderToPDFWithOptions(node, &buf, opts)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("renders with custom page size", func(t *testing.T) {
		node := Box(Sizing(Fixed(100), Fixed(100)))
		Layout(node)

		opts := DefaultPDFOptions()
		opts.PageSize = "Letter"

		var buf bytes.Buffer
		err := RenderToPDFWithOptions(node, &buf, opts)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("renders with custom margins", func(t *testing.T) {
		node := Box(Sizing(Fixed(100), Fixed(100)))
		Layout(node)

		opts := PDFOptions{
			MarginTop:    50,
			MarginRight:  50,
			MarginBottom: 50,
			MarginLeft:   50,
		}

		var buf bytes.Buffer
		err := RenderToPDFWithOptions(node, &buf, opts)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("renders with custom font", func(t *testing.T) {
		node := Box(Sizing(Fixed(100), Fixed(100)))
		Layout(node)

		opts := PDFOptions{
			DefaultFont:     "Times",
			DefaultFontSize: 14,
		}

		var buf bytes.Buffer
		err := RenderToPDFWithOptions(node, &buf, opts)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestDefaultPDFOptions(t *testing.T) {
	opts := DefaultPDFOptions()

	if opts.Landscape {
		t.Error("expected Landscape to be false")
	}
	if opts.PageSize != "A4" {
		t.Errorf("expected PageSize to be A4, got %s", opts.PageSize)
	}
	if opts.MarginTop != 72 {
		t.Errorf("expected MarginTop to be 72, got %f", opts.MarginTop)
	}
	if opts.MarginRight != 72 {
		t.Errorf("expected MarginRight to be 72, got %f", opts.MarginRight)
	}
	if opts.MarginBottom != 72 {
		t.Errorf("expected MarginBottom to be 72, got %f", opts.MarginBottom)
	}
	if opts.MarginLeft != 72 {
		t.Errorf("expected MarginLeft to be 72, got %f", opts.MarginLeft)
	}
	if opts.DefaultFont != "Arial" {
		t.Errorf("expected DefaultFont to be Arial, got %s", opts.DefaultFont)
	}
	if opts.DefaultFontSize != 12 {
		t.Errorf("expected DefaultFontSize to be 12, got %f", opts.DefaultFontSize)
	}
}

func TestDetectImageType(t *testing.T) {
	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := detectImageType("/nonexistent/file.png")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})

	t.Run("returns error for unsupported image type", func(t *testing.T) {
		// Create a temp file with text content
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(tmpFile, []byte("this is not an image"), 0o644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		_, err = detectImageType(tmpFile)
		if err == nil {
			t.Error("expected error for unsupported image type")
		}
		if !strings.Contains(err.Error(), "unsupported") {
			t.Errorf("expected error about unsupported type, got: %v", err)
		}
	})

	t.Run("detects PNG image", func(t *testing.T) {
		// PNG magic bytes
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		pngData := append(pngHeader, make([]byte, 504)...)

		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.png")
		err := os.WriteFile(tmpFile, pngData, 0o644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		imgType, err := detectImageType(tmpFile)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if imgType != "PNG" {
			t.Errorf("expected PNG, got %s", imgType)
		}
	})

	t.Run("detects JPEG image", func(t *testing.T) {
		// JPEG magic bytes
		jpgData := []byte{0xFF, 0xD8, 0xFF, 0xE0}
		jpgData = append(jpgData, make([]byte, 508)...)

		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.jpg")
		err := os.WriteFile(tmpFile, jpgData, 0o644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		imgType, err := detectImageType(tmpFile)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if imgType != "JPG" {
			t.Errorf("expected JPG, got %s", imgType)
		}
	})

	t.Run("detects GIF image", func(t *testing.T) {
		// GIF magic bytes
		gifData := []byte("GIF89a")
		gifData = append(gifData, make([]byte, 506)...)

		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.gif")
		err := os.WriteFile(tmpFile, gifData, 0o644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		imgType, err := detectImageType(tmpFile)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if imgType != "GIF" {
			t.Errorf("expected GIF, got %s", imgType)
		}
	})
}

func TestRenderBoxWithoutVisualProperties(t *testing.T) {
	t.Run("box without border or background renders empty", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(100), Fixed(100)),
			// No border, no background color
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("box with only border", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(100), Fixed(100)),
			Border(2),
			BorderColor("#000000"),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("box with only background", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(100), Fixed(100)),
			BackgroundColor("#FF0000"),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("box with invalid border color", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(100), Fixed(100)),
			Border(2),
			BorderColor("invalid"),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err == nil {
			t.Error("expected error for invalid border color")
		}
	})

	t.Run("box with invalid background color", func(t *testing.T) {
		node := Box(
			Sizing(Fixed(100), Fixed(100)),
			BackgroundColor("invalid"),
		)
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err == nil {
			t.Error("expected error for invalid background color")
		}
	})
}

func TestTextRenderingWithDifferentAlignments(t *testing.T) {
	alignments := []struct {
		name       string
		horizontal Horizontal
		vertical   Vertical
	}{
		{"left-top", Left, Top},
		{"center-top", Center, Top},
		{"right-top", Right, Top},
		{"left-middle", Left, Middle},
		{"center-middle", Center, Middle},
		{"right-middle", Right, Middle},
		{"left-bottom", Left, Bottom},
		{"center-bottom", Center, Bottom},
		{"right-bottom", Right, Bottom},
	}

	for _, align := range alignments {
		t.Run(align.name, func(t *testing.T) {
			node := Box(
				Sizing(Fixed(200), Fixed(100)),
				Alignment(align.horizontal, align.vertical),
				Children(
					Text("Test", FontSize(12)),
				),
			)
			Layout(node)

			var buf bytes.Buffer
			err := RenderToPDF(&buf, node)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestTextWithBorder(t *testing.T) {
	t.Run("text with border renders correctly", func(t *testing.T) {
		node := Text("Hello", FontSize(12), Border(1))
		Layout(node)

		var buf bytes.Buffer
		err := RenderToPDF(&buf, node)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
