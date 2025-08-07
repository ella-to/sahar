package sahar

import (
	"math"
)

type Horizontal int

const (
	Left Horizontal = iota
	Center
	Right
)

type Vertical int

const (
	Top Vertical = iota
	Middle
	Bottom
)

type direction int

const (
	LeftToRight direction = iota
	TopToBottom
)

type Type int

const (
	BoxType Type = iota
	TextType
	ImageType
)

type Position struct {
	X, Y float64
}

type SizeType int

const (
	// FitType is used when the size of the node should fit its content.
	FitType SizeType = iota
	// FixedType is used when the size of the node is fixed to a specific value.
	FixedType
	// GrowType is used when the size of the node should grow to fill available space.
	GrowType
)

const (
	MaxNotSet = -math.MaxFloat64
	MinNotSet = math.MaxFloat64
)

type Size struct {
	Type  SizeType
	Value float64
	// if the value of min or max set to -1
	// it means that the value is not set yet
	Min float64
	Max float64
}

type Node struct {
	Direction       direction
	Type            Type
	Value           string  // For Text nodes
	FontColor       string  // For Text nodes
	FontSize        float64 // For Text nodes
	FontType        string  // For Text nodes
	FontLineHeight  float64 // For Text nodes
	Position        Position
	ChildGap        float64 // Space between children
	Width, Height   Size
	Padding         [4]float64 // Top, Right, Bottom, Left
	Horizontal      Horizontal
	Vertical        Vertical
	Parent          *Node
	Children        []*Node
	Border          float64 // Border width for Box nodes
	BorderColor     string  // Border color for Box nodes
	BackgroundColor string  // Background color for Box nodes
}

var _ nodeOpt = (*Node)(nil)

func (n *Node) configureNode(node *Node) {
	n.Parent = node
	node.Children = append(node.Children, n)
}

//
// Utitlities
//

type border float64

var (
	_ nodeOpt = border(0)
	_ textOpt = border(0)
)

func (b border) configureNode(n *Node) {
	n.Border = float64(b)
}

func (b border) configureText(n *Node) {
	n.Border = float64(b)
}

func Border(width float64) border {
	return border(width)
}

func Text(value string, opts ...textOpt) *Node {
	n := &Node{
		Type:      TextType,
		Direction: LeftToRight,
		Value:     value,
		Width: Size{
			Type:  FitType,
			Value: 0,
			Max:   MaxNotSet,
			Min:   MinNotSet,
		},
		Height: Size{
			Type:  FitType,
			Value: 0,
			Max:   MaxNotSet,
			Min:   MinNotSet,
		},
	}
	for _, opt := range opts {
		opt.configureText(n)
	}

	n.FontLineHeight = n.FontSize / 2

	return n
}

func Box(opts ...nodeOpt) *Node {
	n := &Node{
		Type:      BoxType,
		Direction: LeftToRight,
		Width: Size{
			Type:  FitType,
			Value: 0,
			Max:   MaxNotSet,
			Min:   MinNotSet,
		},
		Height: Size{
			Type:  FitType,
			Value: 0,
			Max:   MaxNotSet,
			Min:   MinNotSet,
		},
	}

	for _, opt := range opts {
		opt.configureNode(n)
	}

	return n
}

func Image(name string, opts ...nodeOpt) *Node {
	n := &Node{
		Type:      ImageType,
		Direction: LeftToRight,
		Value:     name,
		Width: Size{
			Type:  FitType,
			Value: 0,
			Max:   MaxNotSet,
			Min:   MinNotSet,
		},
		Height: Size{
			Type:  FitType,
			Value: 0,
			Max:   MaxNotSet,
			Min:   MinNotSet,
		},
	}

	for _, opt := range opts {
		opt.configureNode(n)
	}

	return n
}

func FontSize(size float64) textOpt {
	return textOptFunc(func(n *Node) {
		n.FontSize = size
	})
}

func FontType(fontType string) textOpt {
	return textOptFunc(func(n *Node) {
		n.FontType = fontType
	})
}

func FontColor(color string) textOpt {
	return textOptFunc(func(n *Node) {
		n.FontColor = color
	})
}

func ChildGap(gap float64) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.ChildGap = gap
	})
}

func Alignment(horizontal Horizontal, vertical Vertical) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.Horizontal = horizontal
		n.Vertical = vertical
	})
}

func Padding(top, right, bottom, left float64) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.Padding[0] = top
		n.Padding[1] = right
		n.Padding[2] = bottom
		n.Padding[3] = left
	})
}

func Min(value float64) fitOpt {
	return fitOptFunc(func(s *Size) {
		s.Min = value
	})
}

func Max(value float64) fitOpt {
	return fitOptFunc(func(s *Size) {
		s.Max = value
	})
}

func Fit(opts ...fitOpt) sizingOpt {
	return sizingOptFunc(func(s *Size) {
		s.Type = FitType
		s.Value = 0 // Fit does not have a specific value, it just fits the content
		s.Max = MaxNotSet
		s.Min = MinNotSet

		for _, opt := range opts {
			opt.configureFit(s)
		}
	})
}

func Fixed(value float64) sizingOpt {
	return sizingOptFunc(func(s *Size) {
		s.Type = FixedType
		s.Value = value
	})
}

func Grow() sizingOpt {
	return sizingOptFunc(func(s *Size) {
		s.Type = GrowType
		s.Value = 0 // Grow does not have a specific value, it just fills available space
	})
}

func A4() []sizingOpt {
	return []sizingOpt{
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 595.28 // A4 width in points
			s.Max = 841.89   // A4 height in points
			s.Min = 0        // No minimum size
		}),
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 841.89 // A4 height in points
			s.Max = 595.28   // A4 width in points
			s.Min = 0        // No minimum size
		}),
	}
}

func Sizing(opts ...sizingOpt) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		switch len(opts) {
		case 0:
			n.Width = Size{
				Type:  FitType,
				Value: 0,
				Max:   MaxNotSet,
				Min:   MinNotSet,
			}
			n.Height = Size{
				Type:  FitType,
				Value: 0,
				Max:   MaxNotSet,
				Min:   MinNotSet,
			}
		case 1:
			opts[0].configureSizing(&n.Width)
			n.Height = Size{
				Type:  FitType,
				Value: 0,
				Max:   MaxNotSet,
				Min:   MinNotSet,
			}
		case 2:
			opts[0].configureSizing(&n.Width)
			opts[1].configureSizing(&n.Height)
		default:
			panic("sizing expects 0, 1, or 2 options")
		}
	})
}

func Direction(dir direction) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.Direction = dir
	})
}

func BackgroundColor(color string) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.BackgroundColor = color
	})
}

func BorderColor(color string) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.BorderColor = color
	})
}

//
// OPTIONS
//

type nodeOpt interface {
	configureNode(*Node)
}

type nodeOptFunc func(*Node)

func (f nodeOptFunc) configureNode(n *Node) {
	f(n)
}

type fitOpt interface {
	configureFit(*Size)
}

type fitOptFunc func(*Size)

func (f fitOptFunc) configureFit(s *Size) {
	f(s)
}

type sizingOpt interface {
	configureSizing(*Size)
}

type sizingOptFunc func(*Size)

func (f sizingOptFunc) configureSizing(s *Size) {
	f(s)
}

type textOpt interface {
	configureText(*Node)
}

type textOptFunc func(*Node)

func (f textOptFunc) configureText(n *Node) {
	f(n)
}

//
//
//
