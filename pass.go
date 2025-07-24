package sahar

import "math"

// Layout performs multi-pass layout calculation on the node tree
func Layout(root *Node) *Node {
	if root == nil {
		return nil
	}

	// Pass 1: Fit Sizing Width (Min and Max)
	calculateFitWidths(root)

	// Pass 2: Grow & Shrink Widths (Min and Max)
	calculateGrowWidths(root)

	// Pass 3: Wrap Text (placeholder for now)
	// TODO: Implement text wrapping logic when text nodes are added

	// Pass 4: Fit Sizing Heights (Min and Max)
	calculateFitHeights(root)

	// Pass 5: Grow & Shrink Heights (Min and Max)
	calculateGrowHeights(root)

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
			// Leaf node - content width is 0 for boxes (would be text width for text nodes)
			contentWidth = 0
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
			// Leaf node - content height is 0 for boxes
			contentHeight = 0
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

// Helper functions
func getActualWidth(node *Node) float64 {
	return node.Width.Value
}

func getActualHeight(node *Node) float64 {
	return node.Height.Value
}
