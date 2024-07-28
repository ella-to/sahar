package sahar

import "errors"

func Resize(node *Node) error {
	if node.Width <= 0 || node.Height <= 0 {
		return errors.New("parent node must have valid width and height")
	}

	isStack := node.Type == Stack
	isGroup := node.Type == Group

	noWidths := make([]*Node, 0)
	noHeights := make([]*Node, 0)

	remainingWidth := node.Width
	remainingHeight := node.Height

	for _, child := range node.Children {
		if isStack {
			remainingHeight -= child.Height
		}

		if isGroup {
			remainingWidth -= child.Width
		}

		if child.Width == 0 {
			noWidths = append(noWidths, child)
		}

		if child.Height == 0 {
			noHeights = append(noHeights, child)
		}
	}

	autoWidth := remainingWidth / float64(len(noWidths))
	autoHeight := remainingHeight / float64(len(noHeights))

	for _, child := range noWidths {
		if isStack {
			child.Width = remainingWidth
		} else {
			child.Width = autoWidth
		}
	}

	for _, child := range noHeights {
		if isGroup {
			child.Height = remainingHeight
		} else {
			child.Height = autoHeight
		}
	}

	for _, child := range node.Children {
		err := Resize(child)
		if err != nil {
			return err
		}
	}

	return nil
}

func Alignment(node *Node) error {
	if node.Width <= 0 || node.Height <= 0 {
		return errors.New("parent node must have valid width and height")
	}

	isStack := node.Type == Stack
	isGroup := node.Type == Group
	horizontal := node.HorizontalAlignment
	vertical := node.VerticalAlignment

	x := node.X
	y := node.Y

	for _, child := range node.Children {
		if isStack {
			switch horizontal {
			case Left:
				child.X = x
			case Center:
				child.X = x + (node.Width-child.Width)/2
			case Right:
				child.X = x + node.Width - child.Width
			}

			switch vertical {
			case Top:
				child.Y = y
			case Middle:
				child.Y = y + (node.Height-child.Height)/2
			case Bottom:
				child.Y = y + node.Height - child.Height
			}

			y += child.Height
		} else if isGroup {
			switch horizontal {
			case Left:
				child.X = x
			case Center:
				child.X = x + (node.Width-child.Width)/2
			case Right:
				child.X = x + node.Width - child.Width
			}

			switch vertical {
			case Top:
				child.Y = y
			case Middle:
				child.Y = y + (node.Height-child.Height)/2
			case Bottom:
				child.Y = y + node.Height - child.Height
			}

			x += child.Width
		}

		err := Alignment(child)
		if err != nil {
			return err
		}
	}

	return nil
}

func Paddings(node *Node) error {
	if node.Width <= 0 || node.Height <= 0 {
		return errors.New("parent node must have valid width and height")
	}

	for _, child := range node.Children {
		child.X += node.Padding[3]
		child.Y += node.Padding[0]
		child.Width -= node.Padding[1]
		child.Height -= node.Padding[2]

		err := Paddings(child)
		if err != nil {
			return err
		}
	}

	return nil
}

func Margins(node *Node) error {
	if node.Width <= 0 || node.Height <= 0 {
		return errors.New("parent node must have valid width and height")
	}

	for _, child := range node.Children {
		child.X += node.Margin[3]
		child.Y += node.Margin[0]
		child.Width -= node.Margin[1]
		child.Height -= node.Margin[2]

		err := Margins(child)
		if err != nil {
			return err
		}
	}

	return nil
}
