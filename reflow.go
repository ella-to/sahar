package sahar

func UpdateChildrenWidthHeight(node *Node) {
	children := node.Children

	parentWidth := node.Width - node.Margin[1] - node.Margin[3]
	parentHeight := node.Height - node.Margin[0] - node.Margin[2]

	noWidths := make([]*Node, 0)
	noHeights := make([]*Node, 0)

	var occupiedWidth float64
	var occupiedHeight float64

	isGroup := node.Type == Group

	for _, child := range children {
		occupiedWidth += child.Width + child.Margin[1] + child.Margin[3]
		occupiedHeight += child.Height + child.Margin[0] + child.Margin[2]

		if child.Width == 0 {
			noWidths = append(noWidths, child)

			// it basically means that all the children inside a non group node
			// that has not set their width will expand their width to the maximum
			// of the container minus the margin left and right
			if !isGroup {
				child.Width = parentWidth - (child.Margin[1] + child.Margin[3])
			}
		}

		if child.Height == 0 {
			noHeights = append(noHeights, child)

			// it basically means that all the children inside a group node
			// that has not set their height will expand their height to the maximum
			// of the container minus the margin top and bottom
			if isGroup {
				child.Height = parentHeight - (child.Margin[0] + child.Margin[2])
			}
		}
	}

	childrenCount := len(children)
	autoWidth := (parentWidth - occupiedWidth) / float64(childrenCount)
	autoHeight := (parentHeight - occupiedHeight) / float64(childrenCount)

	for _, child := range noWidths {
		if child.Width == 0 {
			child.Width = autoWidth
		}
	}

	for _, child := range noHeights {
		if child.Height == 0 {
			child.Height = autoHeight
		}
	}

	for _, child := range children {
		UpdateChildrenWidthHeight(child)
	}
}

func UpdateRootXY(node *Node) {
	node.X += node.Margin[3]
	node.Y += node.Margin[0]
}

func UpdateChildrenXY(node *Node, isRoot bool) {
	isGroup := node.Type == Group

	parentWidth := node.Width
	parentHeight := node.Height

	if isRoot {
		parentWidth = parentWidth - node.Margin[1] - node.Margin[3]
		parentHeight = parentHeight - node.Margin[0] - node.Margin[2]
	}

	x := node.X
	y := node.Y

	for _, child := range node.Children {
		switch node.HorizontalAlignment {
		case Left:
			child.X = x + child.Margin[3]
		case Center:
			child.X = x + (parentWidth-child.Width)/2 + child.Margin[3]
		case Right:
			child.X = x + parentWidth - child.Width - child.Margin[1]
		}

		switch node.VerticalAlignment {
		case Top:
			child.Y = y + child.Margin[0]
		case Middle:
			child.Y = y + (parentHeight-child.Height)/2 + child.Margin[0]
		case Bottom:
			child.Y = y + parentHeight - child.Height - child.Margin[2]
		}

		if isGroup {
			x = child.X + child.Width + child.Margin[1]
		} else {
			y = child.Y + child.Height + child.Margin[2]
		}

		UpdateChildrenXY(child, false)
	}
}
