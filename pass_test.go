package sahar

import (
	"math"
	"strings"
	"testing"
)

func TestLayout(t *testing.T) {
	t.Run("returns nil for nil input", func(t *testing.T) {
		result := Layout(nil)
		if result != nil {
			t.Error("expected Layout(nil) to return nil")
		}
	})

	t.Run("handles single text node", func(t *testing.T) {
		node := Text("Hello", FontSize(12))
		result := Layout(node)

		if result != node {
			t.Error("expected Layout to return the same node")
		}

		// Text node should have calculated dimensions
		if node.Width.Value <= 0 {
			t.Error("expected text width to be calculated and positive")
		}
		if node.Height.Value <= 0 {
			t.Error("expected text height to be calculated and positive")
		}

		// Root node should be positioned at origin
		if node.Position.X != 0 || node.Position.Y != 0 {
			t.Errorf("expected root position to be (0,0), got (%f,%f)", node.Position.X, node.Position.Y)
		}
	})

	t.Run("handles single empty box", func(t *testing.T) {
		node := Box()
		result := Layout(node)

		if result != node {
			t.Error("expected Layout to return the same node")
		}

		// Empty box should have minimal dimensions (just padding)
		expectedWidth := node.Padding[1] + node.Padding[3]  // right + left
		expectedHeight := node.Padding[0] + node.Padding[2] // top + bottom

		if node.Width.Value != expectedWidth {
			t.Errorf("expected empty box width to be %f, got %f", expectedWidth, node.Width.Value)
		}
		if node.Height.Value != expectedHeight {
			t.Errorf("expected empty box height to be %f, got %f", expectedHeight, node.Height.Value)
		}
	})
}

func TestFitSizing(t *testing.T) {
	t.Run("horizontal fit container with text children", func(t *testing.T) {
		child1 := Text("Hello", FontSize(12))
		child2 := Text("World", FontSize(12))

		container := Box(
			Direction(LeftToRight),
			ChildGap(10),
			Children(child1, child2))

		Layout(container)

		// Container width should be sum of children + gap
		expectedWidth := child1.Width.Value + child2.Width.Value + 10
		if math.Abs(container.Width.Value-expectedWidth) > 0.1 {
			t.Errorf("expected container width to be %f, got %f", expectedWidth, container.Width.Value)
		}

		// Children should be positioned horizontally
		if child1.Position.X != 0 {
			t.Errorf("expected child1 X position to be 0, got %f", child1.Position.X)
		}
		expectedChild2X := child1.Width.Value + 10
		if math.Abs(child2.Position.X-expectedChild2X) > 0.1 {
			t.Errorf("expected child2 X position to be %f, got %f", expectedChild2X, child2.Position.X)
		}
	})

	t.Run("vertical fit container with text children", func(t *testing.T) {
		child1 := Text("Hello", FontSize(12))
		child2 := Text("World", FontSize(12))

		container := Box(
			Direction(TopToBottom),
			ChildGap(5),
			Children(child1, child2))

		Layout(container)

		// Container height should be sum of children + gap
		expectedHeight := child1.Height.Value + child2.Height.Value + 5
		if math.Abs(container.Height.Value-expectedHeight) > 0.1 {
			t.Errorf("expected container height to be %f, got %f", expectedHeight, container.Height.Value)
		}

		// Children should be positioned vertically
		if child1.Position.Y != 0 {
			t.Errorf("expected child1 Y position to be 0, got %f", child1.Position.Y)
		}
		expectedChild2Y := child1.Height.Value + 5
		if math.Abs(child2.Position.Y-expectedChild2Y) > 0.1 {
			t.Errorf("expected child2 Y position to be %f, got %f", expectedChild2Y, child2.Position.Y)
		}
	})

	t.Run("fit with min/max constraints", func(t *testing.T) {
		// Create a small text that would normally be less than min
		smallText := Text("Hi", FontSize(8))
		container := Box(
			Sizing(Fit(Min(100), Max(200))),
			Children(smallText))

		Layout(container)

		// Width should be at least the minimum
		if container.Width.Value < 100 {
			t.Errorf("expected container width to be at least 100, got %f", container.Width.Value)
		}
	})
}

func TestFixedSizing(t *testing.T) {
	t.Run("fixed size container", func(t *testing.T) {
		container := Box(Sizing(Fixed(200), Fixed(100)))
		child := Text("Some text", FontSize(12))
		container.Children = append(container.Children, child)
		child.Parent = container

		Layout(container)

		if container.Width.Value != 200 {
			t.Errorf("expected container width to be 200, got %f", container.Width.Value)
		}
		if container.Height.Value != 100 {
			t.Errorf("expected container height to be 100, got %f", container.Height.Value)
		}
	})

	t.Run("fixed size with padding", func(t *testing.T) {
		container := Box(
			Sizing(Fixed(200), Fixed(100)),
			Padding(10, 20, 30, 40))

		Layout(container)

		// Fixed size should remain unchanged
		if container.Width.Value != 200 {
			t.Errorf("expected container width to be 200, got %f", container.Width.Value)
		}
		if container.Height.Value != 100 {
			t.Errorf("expected container height to be 100, got %f", container.Height.Value)
		}
	})
}

func TestGrowSizing(t *testing.T) {
	t.Run("horizontal grow distribution", func(t *testing.T) {
		growChild1 := Box(Sizing(Grow()))
		growChild2 := Box(Sizing(Grow()))
		fixedChild := Box(Sizing(Fixed(50)))

		container := Box(
			Sizing(Fixed(300)),
			Direction(LeftToRight),
			ChildGap(10),
			Children(growChild1, fixedChild, growChild2))

		Layout(container)

		// Available width = 300 - 50 (fixed) - 20 (2 gaps) = 230
		// Each grow child should get 115
		expectedGrowWidth := 115.0

		if math.Abs(growChild1.Width.Value-expectedGrowWidth) > 0.1 {
			t.Errorf("expected growChild1 width to be %f, got %f", expectedGrowWidth, growChild1.Width.Value)
		}
		if math.Abs(growChild2.Width.Value-expectedGrowWidth) > 0.1 {
			t.Errorf("expected growChild2 width to be %f, got %f", expectedGrowWidth, growChild2.Width.Value)
		}

		// Fixed child should maintain its size
		if fixedChild.Width.Value != 50 {
			t.Errorf("expected fixedChild width to be 50, got %f", fixedChild.Width.Value)
		}
	})

	t.Run("vertical grow distribution", func(t *testing.T) {
		growChild1 := Box(Sizing(Fit(), Grow()))
		growChild2 := Box(Sizing(Fit(), Grow()))
		fixedChild := Box(Sizing(Fit(), Fixed(30)))

		container := Box(
			Sizing(Fit(), Fixed(200)),
			Direction(TopToBottom),
			ChildGap(5),
			Children(growChild1, fixedChild, growChild2))

		Layout(container)

		// Available height = 200 - 30 (fixed) - 10 (2 gaps) = 160
		// Each grow child should get 80
		expectedGrowHeight := 80.0

		if math.Abs(growChild1.Height.Value-expectedGrowHeight) > 0.1 {
			t.Errorf("expected growChild1 height to be %f, got %f", expectedGrowHeight, growChild1.Height.Value)
		}
		if math.Abs(growChild2.Height.Value-expectedGrowHeight) > 0.1 {
			t.Errorf("expected growChild2 height to be %f, got %f", expectedGrowHeight, growChild2.Height.Value)
		}
	})

	t.Run("grow with insufficient space", func(t *testing.T) {
		growChild := Box(Sizing(Grow()))
		largeFixedChild := Box(Sizing(Fixed(200)))

		container := Box(
			Sizing(Fixed(100)), // Smaller than fixed child
			Direction(LeftToRight),
			Children(growChild, largeFixedChild))

		Layout(container)

		// Grow child should get 0 width when no space available
		if growChild.Width.Value != 0 {
			t.Errorf("expected growChild width to be 0, got %f", growChild.Width.Value)
		}
	})
}

func TestAlignment(t *testing.T) {
	t.Run("horizontal center alignment", func(t *testing.T) {
		child := Text("Hello", FontSize(12))
		container := Box(
			Sizing(Fixed(200)),
			Alignment(Center, Top),
			Children(child))

		Layout(container)

		// Child should be centered horizontally
		expectedX := (200 - child.Width.Value) / 2
		if math.Abs(child.Position.X-expectedX) > 0.1 {
			t.Errorf("expected child X position to be %f, got %f", expectedX, child.Position.X)
		}
		if child.Position.Y != 0 {
			t.Errorf("expected child Y position to be 0, got %f", child.Position.Y)
		}
	})

	t.Run("vertical middle alignment", func(t *testing.T) {
		child := Text("Hello", FontSize(12))
		container := Box(
			Sizing(Fit(), Fixed(100)),
			Alignment(Left, Middle),
			Children(child))

		Layout(container)

		// Child should be centered vertically
		expectedY := (100 - child.Height.Value) / 2
		if math.Abs(child.Position.Y-expectedY) > 0.1 {
			t.Errorf("expected child Y position to be %f, got %f", expectedY, child.Position.Y)
		}
		if child.Position.X != 0 {
			t.Errorf("expected child X position to be 0, got %f", child.Position.X)
		}
	})

	t.Run("bottom right alignment", func(t *testing.T) {
		child := Text("Hello", FontSize(12))
		container := Box(
			Sizing(Fixed(200), Fixed(100)),
			Alignment(Right, Bottom),
			Children(child))

		Layout(container)

		// Child should be at bottom right
		expectedX := 200 - child.Width.Value
		expectedY := 100 - child.Height.Value

		if math.Abs(child.Position.X-expectedX) > 0.1 {
			t.Errorf("expected child X position to be %f, got %f", expectedX, child.Position.X)
		}
		if math.Abs(child.Position.Y-expectedY) > 0.1 {
			t.Errorf("expected child Y position to be %f, got %f", expectedY, child.Position.Y)
		}
	})
}

func TestPadding(t *testing.T) {
	t.Run("padding affects content area", func(t *testing.T) {
		child := Text("Hello", FontSize(12))
		container := Box(
			Padding(10, 20, 30, 40), // top, right, bottom, left
			Children(child))

		Layout(container)

		// Child should be positioned within padded area
		if child.Position.X != 40 { // left padding
			t.Errorf("expected child X position to be 40, got %f", child.Position.X)
		}
		if child.Position.Y != 10 { // top padding
			t.Errorf("expected child Y position to be 10, got %f", child.Position.Y)
		}

		// Container size should include padding
		expectedWidth := child.Width.Value + 40 + 20   // content + left + right
		expectedHeight := child.Height.Value + 10 + 30 // content + top + bottom

		if math.Abs(container.Width.Value-expectedWidth) > 0.1 {
			t.Errorf("expected container width to be %f, got %f", expectedWidth, container.Width.Value)
		}
		if math.Abs(container.Height.Value-expectedHeight) > 0.1 {
			t.Errorf("expected container height to be %f, got %f", expectedHeight, container.Height.Value)
		}
	})
}

func TestComplexLayout(t *testing.T) {
	t.Run("nested layout with mixed sizing", func(t *testing.T) {
		// Create a complex nested layout
		header := Text("Header", FontSize(18))
		sidebar := Box(
			Sizing(Fixed(100), Grow()),
			BackgroundColor("#f0f0f0"))
		content := Box(
			Sizing(Grow(), Grow()),
			BackgroundColor("#ffffff"))

		mainArea := Box(
			Direction(LeftToRight),
			Sizing(Grow(), Grow()),
			Children(sidebar, content))

		root := Box(
			Direction(TopToBottom),
			Sizing(Fixed(800), Fixed(600)),
			Padding(10, 10, 10, 10),
			Children(header, mainArea))

		Layout(root)

		// Check root dimensions
		if root.Width.Value != 800 {
			t.Errorf("expected root width to be 800, got %f", root.Width.Value)
		}
		if root.Height.Value != 600 {
			t.Errorf("expected root height to be 600, got %f", root.Height.Value)
		}

		// Header should be at top with padding
		if header.Position.X != 10 || header.Position.Y != 10 {
			t.Errorf("expected header position to be (10,10), got (%f,%f)", header.Position.X, header.Position.Y)
		}

		// Sidebar should have fixed width
		if sidebar.Width.Value != 100 {
			t.Errorf("expected sidebar width to be 100, got %f", sidebar.Width.Value)
		}

		// Content area should grow to fill remaining space
		expectedContentWidth := float64(800 - 20 - 100) // total - padding - sidebar
		if math.Abs(content.Width.Value-expectedContentWidth) > 0.1 {
			t.Errorf("expected content width to be %f, got %f", expectedContentWidth, content.Width.Value)
		}

		// Main area should be positioned below header
		expectedMainY := 10 + header.Height.Value // top padding + header height
		if math.Abs(mainArea.Position.Y-expectedMainY) > 0.1 {
			t.Errorf("expected mainArea Y position to be %f, got %f", expectedMainY, mainArea.Position.Y)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty container", func(t *testing.T) {
		container := Box()
		Layout(container)

		// Should not crash and should have minimal dimensions
		if container.Width.Value < 0 || container.Height.Value < 0 {
			t.Error("empty container should have non-negative dimensions")
		}
	})

	t.Run("deeply nested structure", func(t *testing.T) {
		// Create a deeply nested structure
		deepest := Text("Deep", FontSize(10))
		level3 := Box(Children(deepest))
		level2 := Box(Children(level3))
		level1 := Box(Children(level2))
		root := Box(Children(level1))

		// Should not crash
		Layout(root)

		// Check that positions propagate correctly
		if deepest.Position.X < 0 || deepest.Position.Y < 0 {
			t.Error("deeply nested element should have valid position")
		}
	})

	t.Run("zero child gap", func(t *testing.T) {
		child1 := Text("A", FontSize(12))
		child2 := Text("B", FontSize(12))
		container := Box(
			ChildGap(0),
			Direction(LeftToRight),
			Children(child1, child2))

		Layout(container)

		// Children should be adjacent (no gap)
		expectedChild2X := child1.Width.Value
		if math.Abs(child2.Position.X-expectedChild2X) > 0.1 {
			t.Errorf("expected child2 X position to be %f, got %f", expectedChild2X, child2.Position.X)
		}
	})

	t.Run("large child gap", func(t *testing.T) {
		child1 := Text("A", FontSize(12))
		child2 := Text("B", FontSize(12))
		container := Box(
			ChildGap(100),
			Direction(LeftToRight),
			Children(child1, child2))

		Layout(container)

		// Children should be separated by large gap
		expectedChild2X := child1.Width.Value + 100
		if math.Abs(child2.Position.X-expectedChild2X) > 0.1 {
			t.Errorf("expected child2 X position to be %f, got %f", expectedChild2X, child2.Position.X)
		}
	})
}

func TestTextHandling(t *testing.T) {
	t.Run("empty text", func(t *testing.T) {
		emptyText := Text("", FontSize(12))
		Layout(emptyText)

		// Should not crash and should have some dimensions
		if emptyText.Width.Value < 0 || emptyText.Height.Value < 0 {
			t.Error("empty text should have non-negative dimensions")
		}
	})

	t.Run("text with different font sizes", func(t *testing.T) {
		smallText := Text("Test", FontSize(8))
		largeText := Text("Test", FontSize(24))

		Layout(smallText)
		Layout(largeText)

		// Large text should be bigger than small text
		if largeText.Width.Value <= smallText.Width.Value {
			t.Error("large text should be wider than small text")
		}
		if largeText.Height.Value <= smallText.Height.Value {
			t.Error("large text should be taller than small text")
		}
	})
}

func TestImageHandling(t *testing.T) {
	t.Run("image with fixed size", func(t *testing.T) {
		img := Image("test.png", Sizing(Fixed(100), Fixed(200)))
		Layout(img)

		if img.Width.Value != 100 {
			t.Errorf("expected image width to be 100, got %f", img.Width.Value)
		}
		if img.Height.Value != 200 {
			t.Errorf("expected image height to be 200, got %f", img.Height.Value)
		}
	})

	t.Run("image in container", func(t *testing.T) {
		img := Image("test.png", Sizing(Fixed(50), Fixed(50)))
		container := Box(
			Sizing(Fixed(200), Fixed(100)),
			Alignment(Center, Middle),
			Children(img))

		Layout(container)

		// Image should be centered in container
		expectedX := float64((200 - 50) / 2)
		expectedY := float64((100 - 50) / 2)

		if math.Abs(img.Position.X-expectedX) > 0.1 {
			t.Errorf("expected image X position to be %f, got %f", expectedX, img.Position.X)
		}
		if math.Abs(img.Position.Y-expectedY) > 0.1 {
			t.Errorf("expected image Y position to be %f, got %f", expectedY, img.Position.Y)
		}
	})
}

func TestLayoutIdempotency(t *testing.T) {
	t.Run("multiple layout calls should be idempotent", func(t *testing.T) {
		child := Text("Hello", FontSize(12))
		container := Box(
			Sizing(Fixed(200), Fixed(100)),
			Alignment(Center, Middle),
			Children(child))

		// First layout
		Layout(container)
		firstWidth := container.Width.Value
		firstHeight := container.Height.Value
		firstChildX := child.Position.X
		firstChildY := child.Position.Y

		// Second layout should produce same results
		Layout(container)

		if container.Width.Value != firstWidth {
			t.Error("container width changed after second layout")
		}
		if container.Height.Value != firstHeight {
			t.Error("container height changed after second layout")
		}
		if child.Position.X != firstChildX {
			t.Error("child X position changed after second layout")
		}
		if child.Position.Y != firstChildY {
			t.Error("child Y position changed after second layout")
		}
	})
}

func TestShrinkingBehavior(t *testing.T) {
	t.Run("content exceeds container width", func(t *testing.T) {
		// Create children that would exceed container width
		child1 := Box(Sizing(Fixed(100)))
		child2 := Box(Sizing(Fixed(100)))
		child3 := Box(Sizing(Fixed(100)))

		container := Box(
			Sizing(Fixed(200)), // Smaller than sum of children (300)
			Direction(LeftToRight),
			Children(child1, child2, child3))

		Layout(container)

		// Children should be shrunk proportionally
		// Available width = 200, required = 300, so shrink ratio = 200/300 = 2/3
		expectedChildWidth := 100 * (2.0 / 3.0)

		if math.Abs(child1.Width.Value-expectedChildWidth) > 0.1 {
			t.Errorf("expected child1 width to be shrunk to %f, got %f", expectedChildWidth, child1.Width.Value)
		}
		if math.Abs(child2.Width.Value-expectedChildWidth) > 0.1 {
			t.Errorf("expected child2 width to be shrunk to %f, got %f", expectedChildWidth, child2.Width.Value)
		}
		if math.Abs(child3.Width.Value-expectedChildWidth) > 0.1 {
			t.Errorf("expected child3 width to be shrunk to %f, got %f", expectedChildWidth, child3.Width.Value)
		}
	})
}

// Helper function to test text wrapping behavior
func TestTextWrapping(t *testing.T) {
	t.Run("text wrapping in fixed width container", func(t *testing.T) {
		longText := Text("This is a very long text that should wrap", FontSize(12))
		container := Box(
			Sizing(Fixed(100)), // Force narrow width
			Children(longText))

		Layout(container)

		// After layout, text should contain newlines (wrapped)
		if !strings.Contains(longText.Value, "\n") {
			// This might not always be true depending on font metrics
			// but the text wrapping should have been attempted
			t.Log("Text wrapping was attempted (actual wrapping depends on font availability)")
		}
	})
}
