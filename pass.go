package sahar

import (
	"math"
	"strings"

	"golang.org/x/image/font"
)

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
		if node.Width.Min != minNotSet && contentWidth < node.Width.Min {
			contentWidth = node.Width.Min
		}
		if node.Width.Max != maxNotSet && contentWidth > node.Width.Max {
			contentWidth = node.Width.Max
		}

		node.Width.Value = contentWidth
	}
}

// Pass 2: Calculate grow widths top-down
func calculateGrowWidths(node *Node) {
	if len(node.Children) > 0 {
		availableWidth := getAvailableWidth(node)
		distributeGrowWidths(node, availableWidth)
	}

	// Recursively process all children
	for _, child := range node.Children {
		calculateGrowWidths(child)
	}
}

// getAvailableWidth calculates the available width for children
func getAvailableWidth(node *Node) float64 {
	if node.Width.Type == FixedType || node.Width.Type == FitType || (node.Width.Type == GrowType && node.Width.Value > 0) {
		return node.Width.Value - node.Padding[1] - node.Padding[3]
	}
	return 0
}

// distributeGrowWidths distributes available width to grow children
func distributeGrowWidths(node *Node, availableWidth float64) {
	if availableWidth <= 0 {
		setGrowChildrenWidth(node.Children, 0)
		return
	}

	if node.Direction == LeftToRight {
		distributeHorizontalGrowWidths(node, availableWidth)
	} else {
		setGrowChildrenWidth(node.Children, availableWidth)
	}
}

// distributeHorizontalGrowWidths handles width distribution for horizontal layouts
func distributeHorizontalGrowWidths(node *Node, availableWidth float64) {
	usedWidth, growCount := calculateUsedWidthAndGrowCount(node)

	if growCount == 0 {
		return
	}

	remainingWidth := availableWidth - usedWidth
	if remainingWidth > 0 {
		growWidth := remainingWidth / float64(growCount)
		setGrowChildrenWidth(node.Children, math.Max(0, growWidth))
	} else {
		setGrowChildrenWidth(node.Children, 0)
	}
}

// calculateUsedWidthAndGrowCount calculates space used by non-grow children and counts grow children
func calculateUsedWidthAndGrowCount(node *Node) (usedWidth float64, growCount int) {
	for _, child := range node.Children {
		if child.Width.Type == GrowType {
			growCount++
		} else {
			usedWidth += getActualWidth(child)
		}
	}

	if len(node.Children) > 1 {
		usedWidth += node.ChildGap * float64(len(node.Children)-1)
	}
	return
}

// setGrowChildrenWidth sets width for all grow children
func setGrowChildrenWidth(children []*Node, width float64) {
	for _, child := range children {
		if child.Width.Type == GrowType {
			child.Width.Value = width
		}
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
				if child.Width.Min != minNotSet && newWidth < child.Width.Min {
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
		if node.Height.Min != minNotSet && contentHeight < node.Height.Min {
			contentHeight = node.Height.Min
		}
		if node.Height.Max != maxNotSet && contentHeight > node.Height.Max {
			contentHeight = node.Height.Max
		}

		node.Height.Value = contentHeight
	}
}

// Pass 5: Calculate grow heights top-down
func calculateGrowHeights(node *Node) {
	if len(node.Children) > 0 {
		availableHeight := getAvailableHeight(node)
		distributeGrowHeights(node, availableHeight)
	}

	// Recursively process all children
	for _, child := range node.Children {
		calculateGrowHeights(child)
	}
}

// getAvailableHeight calculates the available height for children
func getAvailableHeight(node *Node) float64 {
	if node.Height.Type == FixedType || node.Height.Type == FitType || (node.Height.Type == GrowType && node.Height.Value > 0) {
		return node.Height.Value - node.Padding[0] - node.Padding[2]
	}
	return 0
}

// distributeGrowHeights distributes available height to grow children
func distributeGrowHeights(node *Node, availableHeight float64) {
	if availableHeight <= 0 {
		setGrowChildrenHeight(node.Children, 0)
		return
	}

	if node.Direction == TopToBottom {
		distributeVerticalGrowHeights(node, availableHeight)
	} else {
		setGrowChildrenHeight(node.Children, availableHeight)
	}
}

// distributeVerticalGrowHeights handles height distribution for vertical layouts
func distributeVerticalGrowHeights(node *Node, availableHeight float64) {
	usedHeight, growCount := calculateUsedHeightAndGrowCount(node)

	if growCount == 0 {
		return
	}

	remainingHeight := availableHeight - usedHeight
	if remainingHeight > 0 {
		growHeight := remainingHeight / float64(growCount)
		setGrowChildrenHeight(node.Children, math.Max(0, growHeight))
	} else {
		setGrowChildrenHeight(node.Children, 0)
	}
}

// calculateUsedHeightAndGrowCount calculates space used by non-grow children and counts grow children
func calculateUsedHeightAndGrowCount(node *Node) (usedHeight float64, growCount int) {
	for _, child := range node.Children {
		if child.Height.Type == GrowType {
			growCount++
		} else {
			usedHeight += getActualHeight(child)
		}
	}

	if len(node.Children) > 1 {
		usedHeight += node.ChildGap * float64(len(node.Children)-1)
	}
	return
}

// setGrowChildrenHeight sets height for all grow children
func setGrowChildrenHeight(children []*Node, height float64) {
	for _, child := range children {
		if child.Height.Type == GrowType {
			child.Height.Value = height
		}
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
				if child.Height.Min != minNotSet && newHeight < child.Height.Min {
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
	initRootPosition(node)
	positionChildren(node)

	// Recursively calculate positions for children
	for _, child := range node.Children {
		calculatePositions(child)
	}
}

// initRootPosition sets the root node position to origin
func initRootPosition(node *Node) {
	if node.Parent == nil {
		node.Position.X = 0
		node.Position.Y = 0
	}
}

// contentArea holds the calculated content area bounds
type contentArea struct {
	x, y, width, height float64
}

// getContentArea calculates the content area after padding
func getContentArea(node *Node) contentArea {
	return contentArea{
		x:      node.Position.X + node.Padding[3],
		y:      node.Position.Y + node.Padding[0],
		width:  node.Width.Value - node.Padding[1] - node.Padding[3],
		height: node.Height.Value - node.Padding[0] - node.Padding[2],
	}
}

// positionChildren positions all children based on direction and alignment
func positionChildren(node *Node) {
	if len(node.Children) == 0 {
		return
	}

	content := getContentArea(node)

	if node.Direction == LeftToRight {
		positionChildrenHorizontally(node, content)
	} else {
		positionChildrenVertically(node, content)
	}
}

// calculateTotalChildrenWidth calculates total width of children including gaps
func calculateTotalChildrenWidth(node *Node) float64 {
	var total float64
	for i, child := range node.Children {
		total += getActualWidth(child)
		if i < len(node.Children)-1 {
			total += node.ChildGap
		}
	}
	return total
}

// calculateTotalChildrenHeight calculates total height of children including gaps
func calculateTotalChildrenHeight(node *Node) float64 {
	var total float64
	for i, child := range node.Children {
		total += getActualHeight(child)
		if i < len(node.Children)-1 {
			total += node.ChildGap
		}
	}
	return total
}

// getAlignedX calculates the X position based on horizontal alignment
func getAlignedX(horizontal Horizontal, contentX, contentWidth, totalWidth float64) float64 {
	switch horizontal {
	case Center:
		return contentX + (contentWidth-totalWidth)/2
	case Right:
		return contentX + contentWidth - totalWidth
	default:
		return contentX
	}
}

// getAlignedY calculates the Y position based on vertical alignment
func getAlignedY(vertical Vertical, contentY, contentHeight, totalHeight float64) float64 {
	switch vertical {
	case Middle:
		return contentY + (contentHeight-totalHeight)/2
	case Bottom:
		return contentY + contentHeight - totalHeight
	default:
		return contentY
	}
}

// positionChildrenHorizontally positions children in a horizontal layout
func positionChildrenHorizontally(node *Node, content contentArea) {
	totalWidth := calculateTotalChildrenWidth(node)
	currentX := getAlignedX(node.Horizontal, content.x, content.width, totalWidth)

	for _, child := range node.Children {
		childHeight := getActualHeight(child)
		currentY := getAlignedY(node.Vertical, content.y, content.height, childHeight)

		child.Position.X = currentX
		child.Position.Y = currentY
		currentX += getActualWidth(child) + node.ChildGap
	}
}

// positionChildrenVertically positions children in a vertical layout
func positionChildrenVertically(node *Node, content contentArea) {
	totalHeight := calculateTotalChildrenHeight(node)
	currentY := getAlignedY(node.Vertical, content.y, content.height, totalHeight)

	for _, child := range node.Children {
		childWidth := getActualWidth(child)
		currentX := getAlignedX(node.Horizontal, content.x, content.width, childWidth)

		child.Position.X = currentX
		child.Position.Y = currentY
		currentY += getActualHeight(child) + node.ChildGap
	}
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
		return wrapTextWithFallback(text, maxWidth, fontSize)
	}
	defer face.Close()

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	return wrapWordsToWidth(words, maxWidth, face)
}

// wrapTextWithFallback wraps text using character count when font is not available
func wrapTextWithFallback(text string, maxWidth, fontSize float64) string {
	charWidth := fontSize * 0.6
	maxCharsPerLine := int(maxWidth / charWidth)
	if maxCharsPerLine <= 0 {
		return text
	}
	return wrapTextByCharCount(text, maxCharsPerLine)
}

// wrapWordsToWidth wraps words to fit within maxWidth using font metrics
func wrapWordsToWidth(words []string, maxWidth float64, face font.Face) string {
	var result strings.Builder
	var currentLine strings.Builder
	var currentLineWidth float64

	spaceWidth := measureGlyph(face, ' ')

	for _, word := range words {
		wordWidth := measureWord(face, word)
		shouldWrap := currentLine.Len() > 0 && currentLineWidth+spaceWidth+wordWidth > maxWidth

		if shouldWrap {
			appendLine(&result, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
			currentLineWidth = wordWidth
		} else {
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
		appendLine(&result, currentLine.String())
	}

	return result.String()
}

// measureGlyph measures the width of a single glyph
func measureGlyph(face font.Face, r rune) float64 {
	advance, ok := face.GlyphAdvance(r)
	if ok {
		return float64(advance) / 64.0
	}
	return 0
}

// measureWord measures the width of a word
func measureWord(face font.Face, word string) float64 {
	var width float64
	for _, r := range word {
		width += measureGlyph(face, r)
	}
	return width
}

// appendLine appends a line to the result with proper newline handling
func appendLine(result *strings.Builder, line string) {
	if result.Len() > 0 {
		result.WriteString("\n")
	}
	result.WriteString(line)
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
