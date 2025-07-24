package sahar

import "math"

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
	Direction     direction
	Type          Type
	Position      Position
	ChildGap      float64 // Space between children
	Width, Height Size
	Padding       [4]float64 // Top, Right, Bottom, Left
	Horizontal    Horizontal
	Vertical      Vertical
	Parent        *Node
	Children      []*Node
}

var _ nodeOpt = (*Node)(nil)

func (n *Node) configureNode(node *Node) {
	n.Parent = node
	node.Children = append(node.Children, n)
}

//
// Utitlities
//

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

//
//
//
