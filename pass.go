package sahar

import (
	"math"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// FontCache stores loaded fonts to avoid reloading
var FontCache = make(map[string]*truetype.Font)

// Layout performs multi-pass layout calculation on the node tree
func Layout(root *Node) *Node {
	if root == nil {
		return nil
	}

	// Pass 1: Fit Sizing Width (Min and Max)
	calculateFitWidths(root)

	// Pass 2: Grow & Shrink Widths (Min and Max)
	calculateGrowWidths(root)
	shrinkWidths(root)

	// Pass 3: Wrap Text
	wrapText(root)

	// Pass 4: Fit Sizing Heights (Min and Max)
	calculateFitHeights(root)

	// Pass 5: Grow & Shrink Heights (Min and Max)
	calculateGrowHeights(root)
	shrinkHeights(root)

	// Pass 6: Positions & Alignments
	calculatePositions(root)

	return root
}

// Pass 1: Calculate fit widths bottom-up
func calculateFitWidths(node *Node) {
	// First, calculate fit widths for all children
	for _, child := range node.Children {
		calculateFitWidths(child)
	}

	// Then calculate this node's fit width
	if node.Width.Type == FitType {
		var contentWidth float64

		if len(node.Children) == 0 {
			// Leaf node - content width depends on type
			if node.Type == TextType {
				contentWidth = measureTextWidth(node.Value, node.FontSize, node.FontType)
			} else {
				contentWidth = 0
			}
		} else if node.Direction == LeftToRight {
			// Horizontal layout: sum children widths + gaps
			for i, child := range node.Children {
				contentWidth += getActualWidth(child)
				if i < len(node.Children)-1 {
					contentWidth += node.ChildGap
				}
			}
		} else {
			// Vertical layout: max child width
			for _, child := range node.Children {
				childWidth := getActualWidth(child)
				if childWidth > contentWidth {
					contentWidth = childWidth
				}
			}
		}

		// Add padding
		contentWidth += node.Padding[1] + node.Padding[3] // right + left

		// Apply min/max constraints
		if node.Width.Min != MinNotSet && contentWidth < node.Width.Min {
			contentWidth = node.Width.Min
		}
		if node.Width.Max != MaxNotSet && contentWidth > node.Width.Max {
			contentWidth = node.Width.Max
		}

		node.Width.Value = contentWidth
	}
}

// Pass 2: Calculate grow widths top-down
func calculateGrowWidths(node *Node) {
	if len(node.Children) == 0 {
		return
	}

	var availableWidth float64
	if node.Width.Type == FixedType || node.Width.Type == FitType {
		availableWidth = node.Width.Value - node.Padding[1] - node.Padding[3] // subtract padding
	}

	if node.Direction == LeftToRight {
		// Calculate space taken by fixed and fit children
		var usedWidth float64
		var growCount int

		for i, child := range node.Children {
			if child.Width.Type == GrowType {
				growCount++
			} else {
				usedWidth += getActualWidth(child)
			}
			if i < len(node.Children)-1 {
				usedWidth += node.ChildGap
			}
		}

		// Distribute remaining space among grow children
		if growCount > 0 {
			remainingWidth := availableWidth - usedWidth
			if remainingWidth > 0 {
				growWidth := remainingWidth / float64(growCount)
				for _, child := range node.Children {
					if child.Width.Type == GrowType {
						child.Width.Value = math.Max(0, growWidth)
					}
				}
			}
		}
	} else {
		// Vertical layout: all children can use full width
		for _, child := range node.Children {
			if child.Width.Type == GrowType {
				child.Width.Value = math.Max(0, availableWidth)
			}
		}
	}

	// Recursively process children
	for _, child := range node.Children {
		calculateGrowWidths(child)
	}
}

// Shrink widths when content exceeds available space
func shrinkWidths(node *Node) {
	if len(node.Children) == 0 {
		return
	}

	var availableWidth float64
	if node.Width.Type == FixedType || node.Width.Type == FitType {
		availableWidth = node.Width.Value - node.Padding[1] - node.Padding[3]
	}

	if node.Direction == LeftToRight {
		// Calculate total required width
		var totalRequiredWidth float64
		for i, child := range node.Children {
			totalRequiredWidth += getActualWidth(child)
			if i < len(node.Children)-1 {
				totalRequiredWidth += node.ChildGap
			}
		}

		// If content exceeds available space, shrink proportionally
		if totalRequiredWidth > availableWidth && availableWidth > 0 {
			shrinkRatio := availableWidth / totalRequiredWidth
			for _, child := range node.Children {
				newWidth := getActualWidth(child) * shrinkRatio
				// Respect minimum constraints
				if child.Width.Min != MinNotSet && newWidth < child.Width.Min {
					newWidth = child.Width.Min
				}
				child.Width.Value = newWidth
			}
		}
	}

	// Recursively process children
	for _, child := range node.Children {
		shrinkWidths(child)
	}
}

// Pass 3: Wrap text based on available width
func wrapText(node *Node) {
	if node.Type == TextType && node.Value != "" {
		availableWidth := node.Width.Value - node.Padding[1] - node.Padding[3]
		if availableWidth > 0 {
			node.Value = wrapTextToWidth(node.Value, availableWidth, node.FontSize, node.FontType)
		}
	}

	// Recursively process children
	for _, child := range node.Children {
		wrapText(child)
	}
}

// Pass 4: Calculate fit heights bottom-up
func calculateFitHeights(node *Node) {
	// First, calculate fit heights for all children
	for _, child := range node.Children {
		calculateFitHeights(child)
	}

	// Then calculate this node's fit height
	if node.Height.Type == FitType {
		var contentHeight float64

		if len(node.Children) == 0 {
			// Leaf node - content height depends on type
			if node.Type == TextType {
				contentHeight = measureTextHeight(node.Value, node.FontSize, node.FontType)
			} else {
				contentHeight = 0
			}
		} else if node.Direction == TopToBottom {
			// Vertical layout: sum children heights + gaps
			for i, child := range node.Children {
				contentHeight += getActualHeight(child)
				if i < len(node.Children)-1 {
					contentHeight += node.ChildGap
				}
			}
		} else {
			// Horizontal layout: max child height
			for _, child := range node.Children {
				childHeight := getActualHeight(child)
				if childHeight > contentHeight {
					contentHeight = childHeight
				}
			}
		}

		// Add padding
		contentHeight += node.Padding[0] + node.Padding[2] // top + bottom

		// Apply min/max constraints
		if node.Height.Min != MinNotSet && contentHeight < node.Height.Min {
			contentHeight = node.Height.Min
		}
		if node.Height.Max != MaxNotSet && contentHeight > node.Height.Max {
			contentHeight = node.Height.Max
		}

		node.Height.Value = contentHeight
	}
}

// Pass 5: Calculate grow heights top-down
func calculateGrowHeights(node *Node) {
	if len(node.Children) == 0 {
		return
	}

	var availableHeight float64
	if node.Height.Type == FixedType || node.Height.Type == FitType {
		availableHeight = node.Height.Value - node.Padding[0] - node.Padding[2] // subtract padding
	}

	if node.Direction == TopToBottom {
		// Calculate space taken by fixed and fit children
		var usedHeight float64
		var growCount int

		for i, child := range node.Children {
			if child.Height.Type == GrowType {
				growCount++
			} else {
				usedHeight += getActualHeight(child)
			}
			if i < len(node.Children)-1 {
				usedHeight += node.ChildGap
			}
		}

		// Distribute remaining space among grow children
		if growCount > 0 {
			remainingHeight := availableHeight - usedHeight
			if remainingHeight > 0 {
				growHeight := remainingHeight / float64(growCount)
				for _, child := range node.Children {
					if child.Height.Type == GrowType {
						child.Height.Value = math.Max(0, growHeight)
					}
				}
			}
		}
	} else {
		// Horizontal layout: all children can use full height
		for _, child := range node.Children {
			if child.Height.Type == GrowType {
				child.Height.Value = math.Max(0, availableHeight)
			}
		}
	}

	// Recursively process children
	for _, child := range node.Children {
		calculateGrowHeights(child)
	}
}

// Shrink heights when content exceeds available space
func shrinkHeights(node *Node) {
	if len(node.Children) == 0 {
		return
	}

	var availableHeight float64
	if node.Height.Type == FixedType || node.Height.Type == FitType {
		availableHeight = node.Height.Value - node.Padding[0] - node.Padding[2]
	}

	if node.Direction == TopToBottom {
		// Calculate total required height
		var totalRequiredHeight float64
		for i, child := range node.Children {
			totalRequiredHeight += getActualHeight(child)
			if i < len(node.Children)-1 {
				totalRequiredHeight += node.ChildGap
			}
		}

		// If content exceeds available space, shrink proportionally
		if totalRequiredHeight > availableHeight && availableHeight > 0 {
			shrinkRatio := availableHeight / totalRequiredHeight
			for _, child := range node.Children {
				newHeight := getActualHeight(child) * shrinkRatio
				// Respect minimum constraints
				if child.Height.Min != MinNotSet && newHeight < child.Height.Min {
					newHeight = child.Height.Min
				}
				child.Height.Value = newHeight
			}
		}
	}

	// Recursively process children
	for _, child := range node.Children {
		shrinkHeights(child)
	}
}

// Pass 6: Calculate positions and apply alignments
func calculatePositions(node *Node) {
	// Set root position if not set
	if node.Parent == nil {
		// Root node starts at origin
		node.Position.X = 0
		node.Position.Y = 0
	}

	// Calculate positions for children
	var currentX, currentY float64

	// Start from the padded area
	contentX := node.Position.X + node.Padding[3] // left padding
	contentY := node.Position.Y + node.Padding[0] // top padding
	contentWidth := node.Width.Value - node.Padding[1] - node.Padding[3]
	contentHeight := node.Height.Value - node.Padding[0] - node.Padding[2]

	if node.Direction == LeftToRight {
		// Calculate total content width for alignment
		var totalContentWidth float64
		for i, child := range node.Children {
			totalContentWidth += getActualWidth(child)
			if i < len(node.Children)-1 {
				totalContentWidth += node.ChildGap
			}
		}

		// Apply horizontal alignment
		switch node.Horizontal {
		case Left:
			currentX = contentX
		case Center:
			currentX = contentX + (contentWidth-totalContentWidth)/2
		case Right:
			currentX = contentX + contentWidth - totalContentWidth
		}

		// Position children horizontally
		for _, child := range node.Children {
			// Apply vertical alignment for each child
			switch node.Vertical {
			case Top:
				currentY = contentY
			case Middle:
				currentY = contentY + (contentHeight-getActualHeight(child))/2
			case Bottom:
				currentY = contentY + contentHeight - getActualHeight(child)
			}

			child.Position.X = currentX
			child.Position.Y = currentY
			currentX += getActualWidth(child) + node.ChildGap
		}
	} else {
		// TopToBottom direction
		// Calculate total content height for alignment
		var totalContentHeight float64
		for i, child := range node.Children {
			totalContentHeight += getActualHeight(child)
			if i < len(node.Children)-1 {
				totalContentHeight += node.ChildGap
			}
		}

		// Apply vertical alignment
		switch node.Vertical {
		case Top:
			currentY = contentY
		case Middle:
			currentY = contentY + (contentHeight-totalContentHeight)/2
		case Bottom:
			currentY = contentY + contentHeight - totalContentHeight
		}

		// Position children vertically
		for _, child := range node.Children {
			// Apply horizontal alignment for each child
			switch node.Horizontal {
			case Left:
				currentX = contentX
			case Center:
				currentX = contentX + (contentWidth-getActualWidth(child))/2
			case Right:
				currentX = contentX + contentWidth - getActualWidth(child)
			}

			child.Position.X = currentX
			child.Position.Y = currentY
			currentY += getActualHeight(child) + node.ChildGap
		}
	}

	// Recursively calculate positions for children
	for _, child := range node.Children {
		calculatePositions(child)
	}
}

// getFontFace returns a font.Face for the given font type and size
func getFontFace(fontType string, fontSize float64) font.Face {
	ttfFont, exists := FontCache[fontType]
	if !exists {
		return nil
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size: fontSize,
		DPI:  72, // Standard DPI
	})
}

// measureTextWidth measures the width of text using the specified font
func measureTextWidth(text string, fontSize float64, fontType string) float64 {
	face := getFontFace(fontType, fontSize)
	if face == nil {
		// Fallback to approximation if font not available
		return fontSize * 0.6 * float64(len(text))
	}
	defer face.Close()

	// Split text into lines and measure the widest line
	lines := strings.Split(text, "\n")
	var maxWidth float64

	for _, line := range lines {
		var lineWidth float64
		for _, r := range line {
			advance, ok := face.GlyphAdvance(r)
			if ok {
				lineWidth += float64(advance) / 64.0 // Convert from fixed.Int26_6 to float64
			}
		}
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return maxWidth
}

// measureTextHeight measures the height of text using the specified font
func measureTextHeight(text string, fontSize float64, fontType string) float64 {
	face := getFontFace(fontType, fontSize)
	if face == nil {
		// Fallback to approximation if font not available
		lines := strings.Count(text, "\n") + 1
		return fontSize * 1.2 * float64(lines)
	}
	defer face.Close()

	metrics := face.Metrics()
	lineHeight := float64(metrics.Height) / 64.0 // Convert from fixed.Int26_6 to float64

	// Count lines in text
	lines := strings.Count(text, "\n") + 1

	return lineHeight * float64(lines)
}

// wrapTextToWidth wraps text to fit within the specified width
func wrapTextToWidth(text string, maxWidth float64, fontSize float64, fontType string) string {
	if maxWidth <= 0 {
		return text
	}

	face := getFontFace(fontType, fontSize)
	if face == nil {
		// Fallback to character-based wrapping if font not available
		charWidth := fontSize * 0.6
		maxCharsPerLine := int(maxWidth / charWidth)
		if maxCharsPerLine <= 0 {
			return text
		}
		return wrapTextByCharCount(text, maxCharsPerLine)
	}
	defer face.Close()

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var result strings.Builder
	var currentLine strings.Builder
	var currentLineWidth float64

	for _, word := range words {
		// Measure word width
		var wordWidth float64
		for _, r := range word {
			advance, ok := face.GlyphAdvance(r)
			if ok {
				wordWidth += float64(advance) / 64.0
			}
		}

		// Measure space width (if not first word on line)
		var spaceWidth float64
		if currentLine.Len() > 0 {
			advance, ok := face.GlyphAdvance(' ')
			if ok {
				spaceWidth = float64(advance) / 64.0
			}
		}

		// Check if adding this word would exceed the line width
		if currentLine.Len() > 0 && currentLineWidth+spaceWidth+wordWidth > maxWidth {
			// Start a new line
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
			currentLineWidth = wordWidth
		} else {
			// Add word to current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
				currentLineWidth += spaceWidth
			}
			currentLine.WriteString(word)
			currentLineWidth += wordWidth
		}
	}

	// Add the last line
	if currentLine.Len() > 0 {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(currentLine.String())
	}

	return result.String()
}

// wrapTextByCharCount is a fallback function for character-based wrapping
func wrapTextByCharCount(text string, maxCharsPerLine int) string {
	words := strings.Fields(text)
	var result strings.Builder
	var currentLine strings.Builder

	for _, word := range words {
		// Check if adding this word would exceed the line length
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= maxCharsPerLine {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			// Start a new line
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	// Add the last line
	if currentLine.Len() > 0 {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(currentLine.String())
	}

	return result.String()
}

// Helper functions
func getActualWidth(node *Node) float64 {
	return node.Width.Value
}

func getActualHeight(node *Node) float64 {
	return node.Height.Value
}
