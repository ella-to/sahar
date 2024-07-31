package sahar

import "fmt"

type HorizontalAlignment int

const (
	Left HorizontalAlignment = iota
	Center
	Right
)

type VerticalAlignment int

const (
	Top VerticalAlignment = iota
	Middle
	Bottom
)

type Type int

const (
	Stack Type = iota
	Group
	Image
	Text
)

type Node struct {
	X, Y                float64
	Width, Height       float64
	Margin              [4]float64
	Padding             [4]float64
	Type                Type
	HorizontalAlignment HorizontalAlignment
	VerticalAlignment   VerticalAlignment
	Attributes          map[string]any
	Children            []*Node
}

// GetContainer returns the x, y, width, and height of the container of the node.
// The container is the area where the children of the node are placed.
func (n *Node) GetContainer() (x, y, w, h float64) {
	x = n.X + n.Padding[3]
	y = n.Y + n.Padding[0]

	w = n.Width - n.Padding[1] - n.Padding[3]
	h = n.Height - n.Padding[0] - n.Padding[2]

	return
}

// GetChildrenContainer used to calculate the max width and height of the
// all children for current node. It uses x and y coordinate of the parent
func (n *Node) GetChildrenContainer() (x, y, w, h float64) {
	cx, cy, cw, ch := n.GetContainer()

	isStack := n.Type == Stack
	isGroup := n.Type == Group

	for _, child := range n.Children {
		if isStack {
			h += child.Height
		} else {
			w += child.Width
		}
	}

	if isStack {
		w = cw
	} else if isGroup {
		h = ch
	}

	// align the position x and y
	switch n.HorizontalAlignment {
	case Left:
		x = cx
	case Center:
		x = cx + (cw-w)/2
	case Right:
		x = cx + cw - w
	}

	switch n.VerticalAlignment {
	case Top:
		y = cy
	case Middle:
		y = cy + (ch-h)/2
	case Bottom:
		y = cy + ch - h
	}

	return
}

func (n *Node) AlignChildren() {
	cx, cy, cw, ch := n.GetChildrenContainer()
	isStack := n.Type == Stack
	isGroup := n.Type == Group

	x := cx
	y := cy

	for _, child := range n.Children {
		if isGroup {
			child.X = x
			x += child.Width

			switch n.VerticalAlignment {
			case Top:
				child.Y = y
			case Middle:
				child.Y = y + (ch-child.Height)/2
			case Bottom:
				child.Y = y + ch - child.Height
			}
		} else if isStack {
			child.Y = y
			y += child.Height

			switch n.HorizontalAlignment {
			case Left:
				child.X = x
			case Center:
				child.X = x + (cw-child.Width)/2
			case Right:
				child.X = x + cw - child.Width
			}
		}

		child.AlignChildren()
	}
}

func Block(typ Type, opts ...blockOpt) *Node {
	block := &Node{
		Type:                typ,
		Margin:              [4]float64{0, 0, 0, 0},
		Padding:             [4]float64{0, 0, 0, 0},
		HorizontalAlignment: Left,
		VerticalAlignment:   Top,
		Attributes:          make(map[string]any),
		Children:            make([]*Node, 0),
	}

	for _, opt := range opts {
		opt.configureNode(block)
	}

	return block
}

var _ blockOpt = (*Node)(nil)

func (n *Node) configureNode(node *Node) {
	node.Children = append(node.Children, n)
}

//
// Block options
//

type blockOpt interface {
	configureNode(*Node)
}

type blockOptFunc func(*Node)

func (f blockOptFunc) configureNode(n *Node) {
	f(n)
}

func Margin(top, right, bottom, left float64) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Margin = [4]float64{top, right, bottom, left}
	})
}

func Padding(top, right, bottom, left float64) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Padding = [4]float64{top, right, bottom, left}
	})
}

func Horizontal(h HorizontalAlignment) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.HorizontalAlignment = h
	})
}

func Vertical(v VerticalAlignment) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.VerticalAlignment = v
	})
}

func Alignments(horizontal HorizontalAlignment, vertical VerticalAlignment) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.HorizontalAlignment = horizontal
		n.VerticalAlignment = vertical
	})
}

func Attr(key string, value any) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Attributes[key] = value
	})
}

func Height(height float64) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Height = height
	})
}

func Width(width float64) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = width
	})
}

func X(x float64) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.X = x
	})
}

func Y(y float64) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Y = y
	})
}

func XY(x, y float64) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.X = x
		n.Y = y
	})
}

func FontSize(size float64) blockOpt {
	return Attr("font-size", size)
}

func FontFamily(family string) blockOpt {
	return Attr("font-family", family)
}

func FontWeight(weight string) blockOpt {
	return Attr("font-weight", weight)
}

func TextColor(color string) blockOpt {
	return Attr("color", color)
}

func BackgroundColor(color string) blockOpt {
	return Attr("background-color", color)
}

func Src(src string) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Attributes["src"] = src
	})
}

func TextValue(format string, args ...any) blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Attributes["value"] = fmt.Sprintf(format, args...)
	})
}

type blockSizeOpt interface {
	configureNodeSize(*Node)
}

type size struct {
	w, h float64
}

var _ blockSizeOpt = (*size)(nil)
var _ blockOpt = (*size)(nil)

func (s *size) configureNodeSize(n *Node) {
	n.Width = s.w
	n.Height = s.h
}

func (s *size) configureNode(n *Node) {
	n.Width = s.w
	n.Height = s.h
}

func Size(w, h float64) *size {
	return &size{w, h}
}

func A1() *size {
	return &size{1685, 2384}
}

func A2() *size {
	return &size{1190, 1684}
}

func A3() *size {
	return &size{842, 1190}
}

func A4() *size {
	return &size{595, 842}
}

func A4Lanscape() *size {
	return &size{842, 595}
}

func A4Smal() *size {
	return &size{595, 842}
}

func A5() *size {
	return &size{420, 595}
}

func B4() *size {
	return &size{729, 1032}
}

func B5() *size {
	return &size{516, 729}
}
