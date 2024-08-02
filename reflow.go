package sahar

import (
	"errors"
	"fmt"
)

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

	if IsType(node, TextType) {
		// Since we have the width and height,
		// we can process the Text node here
		text, ok := node.Attributes["text"]
		if !ok {
			return nil
		}

		fontFamily, ok := node.Attributes["font-family-src"]
		if !ok {
			return fmt.Errorf("font-family is required for text node")
		}

		fontSize, ok := node.Attributes["font-size"]
		if !ok {
			return fmt.Errorf("font-size is required for text node")
		}

		ff, err := loadFont(fontFamily.(string), fontSize.(float64))
		if err != nil {
			return err
		}

		width, height := measureString(text.(string), ff)
		node.Width = float64(width)
		node.Height = float64(height)
	}

	return nil
}

func copyParentsAttrsToChildren(node *Node, keys ...string) {
	for _, child := range node.Children {
		for _, key := range keys {
			val, ok := node.Attributes[key]
			if !ok {
				continue
			}

			if _, ok := child.Attributes[key]; !ok {
				child.Attributes[key] = val
			}
		}

		copyParentsAttrsToChildren(child, keys...)
	}
}

func Reflow(node *Node) error {
	copyParentsAttrsToChildren(
		node,
		"font-family",
		"font-size",
		"font-family-src",
		"background-color",
		"color",
		"border-width",
		"border-color",
	)

	err := Resize(node)
	if err != nil {
		return err
	}

	node.AlignChildren()

	return nil
}
