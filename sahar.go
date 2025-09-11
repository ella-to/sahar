package sahar

import (
	"math"
)

// Horizontal represents the horizontal alignment of a node.
// It can be Left, Center, or Right.
type Horizontal int

const (
	Left Horizontal = iota
	Center
	Right
)

// Vertical represents the vertical alignment of a node.
// It can be Top, Middle, or Bottom.
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

// Type represents the type of a node.
// It can be BoxType, TextType, or ImageType.
type Type int

const (
	BoxType Type = iota
	TextType
	ImageType
)

// Position represents the position of a node in the layout.
// It contains X and Y coordinates and it will be calculated by the layout engine.
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
	maxNotSet = -math.MaxFloat64
	minNotSet = math.MaxFloat64
)

// Size represents the size of a node.
// It contains the type of size (Fit, Fixed, or Grow), the value, and optional min and max values.
type Size struct {
	Type  SizeType
	Value float64
	// if the value of min or max set to -1
	// it means that the value is not set yet
	Min float64
	Max float64
}

// Node represents a layout node.
// It can be a box, text, or image.
// It contains properties for alignment, size, padding, and children nodes.
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

// Border creates a border option for nodes.
func Border(width float64) border {
	return border(width)
}

// Text creates a new text node with the specified value and options.
// A text node is used to display text with specific font, size, and color.
// It can also have a border for debugging purposes.
func Text(value string, opts ...textOpt) *Node {
	n := &Node{
		Type:      TextType,
		Direction: LeftToRight,
		Value:     value,
		Width: Size{
			Type:  FitType,
			Value: 0,
			Max:   maxNotSet,
			Min:   minNotSet,
		},
		Height: Size{
			Type:  FitType,
			Value: 0,
			Max:   maxNotSet,
			Min:   minNotSet,
		},
	}
	for _, opt := range opts {
		opt.configureText(n)
	}

	n.FontLineHeight = n.FontSize / 2

	return n
}

// Box creates a new box node with the specified options.
// A box node is a container that can hold other nodes and can have a border, background
func Box(opts ...nodeOpt) *Node {
	n := &Node{
		Type:      BoxType,
		Direction: LeftToRight,
		Width: Size{
			Type:  FitType,
			Value: 0,
			Max:   maxNotSet,
			Min:   minNotSet,
		},
		Height: Size{
			Type:  FitType,
			Value: 0,
			Max:   maxNotSet,
			Min:   minNotSet,
		},
	}

	for _, opt := range opts {
		opt.configureNode(n)
	}

	return n
}

// Image creates a new image node with the specified source and options.
func Image(src string, opts ...nodeOpt) *Node {
	n := &Node{
		Type:      ImageType,
		Direction: LeftToRight,
		Value:     src,
		Width: Size{
			Type:  FitType,
			Value: 0,
			Max:   maxNotSet,
			Min:   minNotSet,
		},
		Height: Size{
			Type:  FitType,
			Value: 0,
			Max:   maxNotSet,
			Min:   minNotSet,
		},
	}

	for _, opt := range opts {
		opt.configureNode(n)
	}

	return n
}

// FontSize sets the font size for text nodes.
func FontSize(size float64) textOpt {
	return textOptFunc(func(n *Node) {
		n.FontSize = size
	})
}

// FontType sets the font type for text nodes.
func FontType(fontType string) textOpt {
	return textOptFunc(func(n *Node) {
		n.FontType = fontType
	})
}

// FontColor sets the font color for text nodes.
func FontColor(color string) textOpt {
	return textOptFunc(func(n *Node) {
		n.FontColor = color
	})
}

// ChildGap sets the gap between child nodes in a parent node.
func ChildGap(gap float64) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.ChildGap = gap
	})
}

// Alignment sets the horizontal and vertical alignment of the node.
// Horizontal can be Left, Center, or Right.
// Vertical can be Top, Middle, or Bottom.
func Alignment(horizontal Horizontal, vertical Vertical) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.Horizontal = horizontal
		n.Vertical = vertical
	})
}

// Padding sets the padding for the node.
func Padding(top, right, bottom, left float64) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.Padding[0] = top
		n.Padding[1] = right
		n.Padding[2] = bottom
		n.Padding[3] = left
	})
}

// Min is set the value for Width or Height if they are set to either Fit or Grow
func Min(value float64) fitOpt {
	return fitOptFunc(func(s *Size) {
		s.Min = value
	})
}

// Max is set the value for Width or Height if they are set to either Fit or Grow
func Max(value float64) fitOpt {
	return fitOptFunc(func(s *Size) {
		s.Max = value
	})
}

// Fit mostly used in Parent nodes to let the layout engine know that it can expand to fit the
// children, This is the default value of all nodes
func Fit(opts ...fitOpt) sizingOpt {
	return sizingOptFunc(func(s *Size) {
		s.Type = FitType
		s.Value = 0 // Fit does not have a specific value, it just fits the content
		s.Max = maxNotSet
		s.Min = minNotSet

		for _, opt := range opts {
			opt.configureFit(s)
		}
	})
}

// Fixed accept a value and make sure the element has the value
func Fixed(value float64) sizingOpt {
	return sizingOptFunc(func(s *Size) {
		s.Type = FixedType
		s.Value = value
	})
}

// Grow is a way to set either width or height to gorw and fill the remaining of the space
func Grow() sizingOpt {
	return sizingOptFunc(func(s *Size) {
		s.Type = GrowType
		s.Value = 0 // Grow does not have a specific value, it just fills available space
	})
}

// A4 is a custom sizing and should be used with Sizing
// please use it as spread for example: Sizing(A4()...)
func A4() []sizingOpt {
	return []sizingOpt{
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 595.28 // A4 width in points
			s.Max = 595.28   // A4 width in points
			s.Min = 0        // No minimum size
		}),
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 841.89 // A4 height in points
			s.Max = 841.89   // A4 height in points
			s.Min = 0        // No minimum size
		}),
	}
}

// USLetter is a custom sizing and should be used with Sizing
// please use it as spread for example: Sizing(USLetter()...)
func USLetter() []sizingOpt {
	return []sizingOpt{
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 612 // USLetter width in points
			s.Max = 612   // USLetter width in points
			s.Min = 0     // No minimum size
		}),
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 792 // USLetter height in points
			s.Max = 792   // USLetter height in points
			s.Min = 0     // No minimum size
		}),
	}
}

// USLegal is a custom sizing and should be used with Sizing
// please use it as spread for example: Sizing(USLegal()...)
func USLegal() []sizingOpt {
	return []sizingOpt{
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 612 // USLegal width in points
			s.Max = 612   // USLegal width in points
			s.Min = 0     // No minimum size
		}),
		sizingOptFunc(func(s *Size) {
			s.Type = FixedType
			s.Value = 1008 // USLegal height in points
			s.Max = 1008   // USLegal height in points
			s.Min = 0      // No minimum size
		}),
	}
}

// Sizing accept either 0, 1 or 2. If you pass more than 2 arguments it panics
// 0: default value to use Fit type for both width and height. They will expand to fit the children size
// 1: Set the Width of the node and set the height to fit
// 2: Set both Width and Height of the node
func Sizing(opts ...sizingOpt) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		switch len(opts) {
		case 0:
			n.Width = Size{
				Type:  FitType,
				Value: 0,
				Max:   maxNotSet,
				Min:   minNotSet,
			}
			n.Height = Size{
				Type:  FitType,
				Value: 0,
				Max:   maxNotSet,
				Min:   minNotSet,
			}
		case 1:
			opts[0].configureSizing(&n.Width)
			n.Height = Size{
				Type:  FitType,
				Value: 0,
				Max:   maxNotSet,
				Min:   minNotSet,
			}
		case 2:
			opts[0].configureSizing(&n.Width)
			opts[1].configureSizing(&n.Height)
		default:
			panic("sizing expects 0, 1, or 2 options")
		}
	})
}

// Direction use for rendering the direction of the children. There is two possibilities
// either TopToBottom ot LeftToRight
func Direction(dir direction) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.Direction = dir
	})
}

// BackgroundColor set the background color and the value should be hex color
// either it it has to be either 6 or 7 (including # at the beginning of the color value)
func BackgroundColor(color string) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.BackgroundColor = color
	})
}

// BorderColor set the border color and the value should be hex color
// either it it has to be either 6 or 7 (including # at the beginning of the color value)
func BorderColor(color string) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		n.BorderColor = color
	})
}

// Children appends new children to the parent node
// this is useful if you have a dynamic component
func Children(nodes ...*Node) nodeOpt {
	return nodeOptFunc(func(n *Node) {
		for _, node := range nodes {
			node.configureNode(n)
		}
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
