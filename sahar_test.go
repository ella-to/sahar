package sahar

import (
	"math"
	"testing"
)

func TestText(t *testing.T) {
	t.Run("creates text node with default values", func(t *testing.T) {
		node := Text("Hello World")

		if node.Type != TextType {
			t.Errorf("expected Type to be TextType, got %v", node.Type)
		}
		if node.Value != "Hello World" {
			t.Errorf("expected Value to be 'Hello World', got %s", node.Value)
		}
		if node.Direction != LeftToRight {
			t.Errorf("expected Direction to be LeftToRight, got %v", node.Direction)
		}
		if node.Width.Type != FitType {
			t.Errorf("expected Width Type to be FitType, got %v", node.Width.Type)
		}
		if node.Height.Type != FitType {
			t.Errorf("expected Height Type to be FitType, got %v", node.Height.Type)
		}
		if node.FontLineHeight != node.FontSize/2 {
			t.Errorf("expected FontLineHeight to be FontSize/2, got %f", node.FontLineHeight)
		}
	})

	t.Run("applies text options correctly", func(t *testing.T) {
		node := Text("Test",
			FontSize(14),
			FontColor("#FF0000"),
			FontType("Arial"),
			Border(2))

		if node.FontSize != 14 {
			t.Errorf("expected FontSize to be 14, got %f", node.FontSize)
		}
		if node.FontColor != "#FF0000" {
			t.Errorf("expected FontColor to be '#FF0000', got %s", node.FontColor)
		}
		if node.FontType != "Arial" {
			t.Errorf("expected FontType to be 'Arial', got %s", node.FontType)
		}
		if node.Border != 2 {
			t.Errorf("expected Border to be 2, got %f", node.Border)
		}
		if node.FontLineHeight != 7 { // 14/2
			t.Errorf("expected FontLineHeight to be 7, got %f", node.FontLineHeight)
		}
	})
}

func TestBox(t *testing.T) {
	t.Run("creates box node with default values", func(t *testing.T) {
		node := Box()

		if node.Type != BoxType {
			t.Errorf("expected Type to be BoxType, got %v", node.Type)
		}
		if node.Direction != LeftToRight {
			t.Errorf("expected Direction to be LeftToRight, got %v", node.Direction)
		}
		if node.Width.Type != FitType {
			t.Errorf("expected Width Type to be FitType, got %v", node.Width.Type)
		}
		if node.Height.Type != FitType {
			t.Errorf("expected Height Type to be FitType, got %v", node.Height.Type)
		}
	})

	t.Run("applies node options correctly", func(t *testing.T) {
		child1 := Text("Child 1")
		child2 := Text("Child 2")

		node := Box(
			ChildGap(10),
			Alignment(Center, Middle),
			Padding(5, 10, 15, 20),
			Direction(TopToBottom),
			BackgroundColor("#00FF00"),
			BorderColor("#0000FF"),
			Border(3),
			Children(child1, child2))

		if node.ChildGap != 10 {
			t.Errorf("expected ChildGap to be 10, got %f", node.ChildGap)
		}
		if node.Horizontal != Center {
			t.Errorf("expected Horizontal to be Center, got %v", node.Horizontal)
		}
		if node.Vertical != Middle {
			t.Errorf("expected Vertical to be Middle, got %v", node.Vertical)
		}
		if node.Padding[0] != 5 || node.Padding[1] != 10 || node.Padding[2] != 15 || node.Padding[3] != 20 {
			t.Errorf("expected Padding to be [5 10 15 20], got %v", node.Padding)
		}
		if node.Direction != TopToBottom {
			t.Errorf("expected Direction to be TopToBottom, got %v", node.Direction)
		}
		if node.BackgroundColor != "#00FF00" {
			t.Errorf("expected BackgroundColor to be '#00FF00', got %s", node.BackgroundColor)
		}
		if node.BorderColor != "#0000FF" {
			t.Errorf("expected BorderColor to be '#0000FF', got %s", node.BorderColor)
		}
		if node.Border != 3 {
			t.Errorf("expected Border to be 3, got %f", node.Border)
		}
		if len(node.Children) != 2 {
			t.Errorf("expected 2 children, got %d", len(node.Children))
		}
		if child1.Parent != node {
			t.Error("expected child1 parent to be set")
		}
		if child2.Parent != node {
			t.Error("expected child2 parent to be set")
		}
	})
}

func TestImage(t *testing.T) {
	t.Run("creates image node with default values", func(t *testing.T) {
		node := Image("image.png")

		if node.Type != ImageType {
			t.Errorf("expected Type to be ImageType, got %v", node.Type)
		}
		if node.Value != "image.png" {
			t.Errorf("expected Value to be 'image.png', got %s", node.Value)
		}
		if node.Direction != LeftToRight {
			t.Errorf("expected Direction to be LeftToRight, got %v", node.Direction)
		}
		if node.Width.Type != FitType {
			t.Errorf("expected Width Type to be FitType, got %v", node.Width.Type)
		}
		if node.Height.Type != FitType {
			t.Errorf("expected Height Type to be FitType, got %v", node.Height.Type)
		}
	})

	t.Run("applies node options correctly", func(t *testing.T) {
		node := Image("test.jpg",
			Sizing(Fixed(100), Fixed(200)),
			Border(1))

		if node.Width.Type != FixedType || node.Width.Value != 100 {
			t.Errorf("expected Width to be Fixed(100), got Type: %v, Value: %f", node.Width.Type, node.Width.Value)
		}
		if node.Height.Type != FixedType || node.Height.Value != 200 {
			t.Errorf("expected Height to be Fixed(200), got Type: %v, Value: %f", node.Height.Type, node.Height.Value)
		}
		if node.Border != 1 {
			t.Errorf("expected Border to be 1, got %f", node.Border)
		}
	})
}

func TestSizing(t *testing.T) {
	t.Run("no arguments sets both to Fit", func(t *testing.T) {
		node := Box(Sizing())

		if node.Width.Type != FitType {
			t.Errorf("expected Width Type to be FitType, got %v", node.Width.Type)
		}
		if node.Height.Type != FitType {
			t.Errorf("expected Height Type to be FitType, got %v", node.Height.Type)
		}
	})

	t.Run("one argument sets width, height remains fit", func(t *testing.T) {
		node := Box(Sizing(Fixed(100)))

		if node.Width.Type != FixedType || node.Width.Value != 100 {
			t.Errorf("expected Width to be Fixed(100), got Type: %v, Value: %f", node.Width.Type, node.Width.Value)
		}
		if node.Height.Type != FitType {
			t.Errorf("expected Height Type to be FitType, got %v", node.Height.Type)
		}
	})

	t.Run("two arguments set both width and height", func(t *testing.T) {
		node := Box(Sizing(Fixed(100), Grow()))

		if node.Width.Type != FixedType || node.Width.Value != 100 {
			t.Errorf("expected Width to be Fixed(100), got Type: %v, Value: %f", node.Width.Type, node.Width.Value)
		}
		if node.Height.Type != GrowType {
			t.Errorf("expected Height Type to be GrowType, got %v", node.Height.Type)
		}
	})

	t.Run("panics with more than 2 arguments", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic with more than 2 sizing arguments")
			}
		}()
		Box(Sizing(Fixed(100), Fixed(200), Fixed(300)))
	})
}

func TestSizingOptions(t *testing.T) {
	t.Run("Fixed sizing", func(t *testing.T) {
		var size Size
		Fixed(150).configureSizing(&size)

		if size.Type != FixedType {
			t.Errorf("expected Type to be FixedType, got %v", size.Type)
		}
		if size.Value != 150 {
			t.Errorf("expected Value to be 150, got %f", size.Value)
		}
	})

	t.Run("Fit sizing with min/max", func(t *testing.T) {
		var size Size
		Fit(Min(10), Max(100)).configureSizing(&size)

		if size.Type != FitType {
			t.Errorf("expected Type to be FitType, got %v", size.Type)
		}
		if size.Min != 10 {
			t.Errorf("expected Min to be 10, got %f", size.Min)
		}
		if size.Max != 100 {
			t.Errorf("expected Max to be 100, got %f", size.Max)
		}
		if size.Value != 0 {
			t.Errorf("expected Value to be 0, got %f", size.Value)
		}
	})

	t.Run("Grow sizing", func(t *testing.T) {
		var size Size
		Grow().configureSizing(&size)

		if size.Type != GrowType {
			t.Errorf("expected Type to be GrowType, got %v", size.Type)
		}
		if size.Value != 0 {
			t.Errorf("expected Value to be 0, got %f", size.Value)
		}
	})
}

func TestPresetSizes(t *testing.T) {
	t.Run("A4 preset", func(t *testing.T) {
		node := Box(Sizing(A4()...))

		if node.Width.Type != FixedType || node.Width.Value != 595.28 {
			t.Errorf("expected Width to be Fixed(595.28), got Type: %v, Value: %f", node.Width.Type, node.Width.Value)
		}
		if node.Height.Type != FixedType || node.Height.Value != 841.89 {
			t.Errorf("expected Height to be Fixed(841.89), got Type: %v, Value: %f", node.Height.Type, node.Height.Value)
		}
	})

	t.Run("USLetter preset", func(t *testing.T) {
		node := Box(Sizing(USLetter()...))

		if node.Width.Type != FixedType || node.Width.Value != 612 {
			t.Errorf("expected Width to be Fixed(612), got Type: %v, Value: %f", node.Width.Type, node.Width.Value)
		}
		if node.Height.Type != FixedType || node.Height.Value != 792 {
			t.Errorf("expected Height to be Fixed(792), got Type: %v, Value: %f", node.Height.Type, node.Height.Value)
		}
	})

	t.Run("USLegal preset", func(t *testing.T) {
		node := Box(Sizing(USLegal()...))

		if node.Width.Type != FixedType || node.Width.Value != 612 {
			t.Errorf("expected Width to be Fixed(612), got Type: %v, Value: %f", node.Width.Type, node.Width.Value)
		}
		if node.Height.Type != FixedType || node.Height.Value != 1008 {
			t.Errorf("expected Height to be Fixed(1008), got Type: %v, Value: %f", node.Height.Type, node.Height.Value)
		}
	})
}

func TestAlignmentDef(t *testing.T) {
	tests := []struct {
		name       string
		horizontal Horizontal
		vertical   Vertical
	}{
		{"Left Top", Left, Top},
		{"Center Middle", Center, Middle},
		{"Right Bottom", Right, Bottom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Box(Alignment(tt.horizontal, tt.vertical))

			if node.Horizontal != tt.horizontal {
				t.Errorf("expected Horizontal to be %v, got %v", tt.horizontal, node.Horizontal)
			}
			if node.Vertical != tt.vertical {
				t.Errorf("expected Vertical to be %v, got %v", tt.vertical, node.Vertical)
			}
		})
	}
}

func TestDirection(t *testing.T) {
	t.Run("LeftToRight direction", func(t *testing.T) {
		node := Box(Direction(LeftToRight))
		if node.Direction != LeftToRight {
			t.Errorf("expected Direction to be LeftToRight, got %v", node.Direction)
		}
	})

	t.Run("TopToBottom direction", func(t *testing.T) {
		node := Box(Direction(TopToBottom))
		if node.Direction != TopToBottom {
			t.Errorf("expected Direction to be TopToBottom, got %v", node.Direction)
		}
	})
}

func TestPaddingDef(t *testing.T) {
	node := Box(Padding(1, 2, 3, 4))

	expected := [4]float64{1, 2, 3, 4}
	if node.Padding != expected {
		t.Errorf("expected Padding to be %v, got %v", expected, node.Padding)
	}
}

func TestChildGap(t *testing.T) {
	node := Box(ChildGap(25))

	if node.ChildGap != 25 {
		t.Errorf("expected ChildGap to be 25, got %f", node.ChildGap)
	}
}

func TestColors(t *testing.T) {
	t.Run("BackgroundColor", func(t *testing.T) {
		node := Box(BackgroundColor("#FF00FF"))
		if node.BackgroundColor != "#FF00FF" {
			t.Errorf("expected BackgroundColor to be '#FF00FF', got %s", node.BackgroundColor)
		}
	})

	t.Run("BorderColor", func(t *testing.T) {
		node := Box(BorderColor("#00FFFF"))
		if node.BorderColor != "#00FFFF" {
			t.Errorf("expected BorderColor to be '#00FFFF', got %s", node.BorderColor)
		}
	})

	t.Run("FontColor for text", func(t *testing.T) {
		node := Text("test", FontColor("#FFFF00"))
		if node.FontColor != "#FFFF00" {
			t.Errorf("expected FontColor to be '#FFFF00', got %s", node.FontColor)
		}
	})
}

func TestBorder(t *testing.T) {
	t.Run("Border for Box", func(t *testing.T) {
		node := Box(Border(5))
		if node.Border != 5 {
			t.Errorf("expected Border to be 5, got %f", node.Border)
		}
	})

	t.Run("Border for Text", func(t *testing.T) {
		node := Text("test", Border(2.5))
		if node.Border != 2.5 {
			t.Errorf("expected Border to be 2.5, got %f", node.Border)
		}
	})
}

func TestChildren(t *testing.T) {
	child1 := Text("Child 1")
	child2 := Box()
	child3 := Image("test.png")

	parent := Box(Children(child1, child2, child3))

	if len(parent.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(parent.Children))
	}

	if child1.Parent != parent {
		t.Error("expected child1 parent to be set")
	}
	if child2.Parent != parent {
		t.Error("expected child2 parent to be set")
	}
	if child3.Parent != parent {
		t.Error("expected child3 parent to be set")
	}

	if parent.Children[0] != child1 {
		t.Error("expected first child to be child1")
	}
	if parent.Children[1] != child2 {
		t.Error("expected second child to be child2")
	}
	if parent.Children[2] != child3 {
		t.Error("expected third child to be child3")
	}
}

func TestFontOptions(t *testing.T) {
	t.Run("FontSize", func(t *testing.T) {
		node := Text("test", FontSize(16))
		if node.FontSize != 16 {
			t.Errorf("expected FontSize to be 16, got %f", node.FontSize)
		}
		if node.FontLineHeight != 8 { // 16/2
			t.Errorf("expected FontLineHeight to be 8, got %f", node.FontLineHeight)
		}
	})

	t.Run("FontType", func(t *testing.T) {
		node := Text("test", FontType("Helvetica"))
		if node.FontType != "Helvetica" {
			t.Errorf("expected FontType to be 'Helvetica', got %s", node.FontType)
		}
	})
}

func TestComplexLayoutDef(t *testing.T) {
	// Test a more complex layout structure
	header := Text("Header", FontSize(18), FontColor("#000000"))
	content := Box(
		Direction(TopToBottom),
		ChildGap(10),
		Padding(5, 5, 5, 5),
		Children(
			Text("Paragraph 1", FontSize(12)),
			Text("Paragraph 2", FontSize(12)),
			Image("image.png", Sizing(Fixed(100), Fixed(100))),
		),
	)

	root := Box(
		Direction(TopToBottom),
		Sizing(A4()...),
		Padding(20, 20, 20, 20),
		Children(header, content),
	)

	// Test root properties
	if root.Direction != TopToBottom {
		t.Error("expected root direction to be TopToBottom")
	}
	if len(root.Children) != 2 {
		t.Errorf("expected root to have 2 children, got %d", len(root.Children))
	}

	// Test header
	if header.Parent != root {
		t.Error("expected header parent to be root")
	}
	if header.FontSize != 18 {
		t.Error("expected header font size to be 18")
	}

	// Test content
	if content.Parent != root {
		t.Error("expected content parent to be root")
	}
	if len(content.Children) != 3 {
		t.Errorf("expected content to have 3 children, got %d", len(content.Children))
	}

	// Test image sizing
	img := content.Children[2]
	if img.Type != ImageType {
		t.Error("expected third child of content to be ImageType")
	}
	if img.Width.Type != FixedType || img.Width.Value != 100 {
		t.Error("expected image width to be Fixed(100)")
	}
	if img.Height.Type != FixedType || img.Height.Value != 100 {
		t.Error("expected image height to be Fixed(100)")
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	if maxNotSet != -math.MaxFloat64 {
		t.Errorf("expected maxNotSet to be -math.MaxFloat64, got %f", maxNotSet)
	}
	if minNotSet != math.MaxFloat64 {
		t.Errorf("expected minNotSet to be math.MaxFloat64, got %f", minNotSet)
	}
}

func TestDefaultSizeValues(t *testing.T) {
	// Test that default Size values are set correctly
	node := Box()

	if node.Width.Max != maxNotSet {
		t.Errorf("expected Width.Max to be maxNotSet, got %f", node.Width.Max)
	}
	if node.Width.Min != minNotSet {
		t.Errorf("expected Width.Min to be minNotSet, got %f", node.Width.Min)
	}
	if node.Height.Max != maxNotSet {
		t.Errorf("expected Height.Max to be maxNotSet, got %f", node.Height.Max)
	}
	if node.Height.Min != minNotSet {
		t.Errorf("expected Height.Min to be minNotSet, got %f", node.Height.Min)
	}
}
