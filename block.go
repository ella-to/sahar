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

func A1() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 1685
		n.Height = 2384
	})
}

func A2() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 1190
		n.Height = 1684
	})
}

func A3() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 842
		n.Height = 1190
	})
}

func A4() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 595
		n.Height = 842
	})
}

func A4Lanscape() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 842
		n.Height = 595
	})
}

func A4Smal() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 595
		n.Height = 842
	})
}

func A5() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 420
		n.Height = 595
	})
}

func B4() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 729
		n.Height = 1032
	})
}

func B5() blockOpt {
	return blockOptFunc(func(n *Node) {
		n.Width = 516
		n.Height = 729
	})
}
