package sahar

import "errors"

func Resize(node *Node) error {
	if node.Width <= 0 || node.Height <= 0 {
		return errors.New("parent node must have valid width and height")
	}

	isStack := node.IsStack()
	isGroup := node.IsGroup()

	noWidths := make([]*Node, 0)
	noHeights := make([]*Node, 0)

	remainingWidth := node.Width - node.Padding[1] - node.Padding[3]
	remainingHeight := node.Height - node.Padding[0] - node.Padding[2]

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

func Reflow(node *Node) error {
	err := Resize(node)
	if err != nil {
		return err
	}

	node.AlignChildren()
	return nil
}
