package sahar

import "os"

// Debug function to print the layout tree
func PrintLayout(node *Node, indent int) {
	if node == nil {
		return
	}

	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	println(indentStr + "Node:")
	println(indentStr + "  Position: (" + floatToString(node.Position.X) + ", " + floatToString(node.Position.Y) + ")")
	println(indentStr + "  Size: " + floatToString(node.Width.Value) + " x " + floatToString(node.Height.Value))
	println(indentStr + "  Type: " + sizeTypeToString(node.Width.Type) + " / " + sizeTypeToString(node.Height.Type))

	for _, child := range node.Children {
		PrintLayout(child, indent+1)
	}
}

func floatToString(f float64) string {
	// Simple float to string conversion
	return "float64"
}

func sizeTypeToString(t SizeType) string {
	switch t {
	case FitType:
		return "Fit"
	case FixedType:
		return "Fixed"
	case GrowType:
		return "Grow"
	default:
		return "Unknown"
	}
}

// Example usage demonstrating the layout engine
func ExampleUsage() {
	// Create the layout tree as requested
	n := Layout(
		Box(
			// Parent container
			Sizing(
				Fixed(100), // width
				Fixed(100), // height
			),
			Direction(TopToBottom),
			// ChildGap(5),
			// Padding(10, 10, 10, 10), // top, right, bottom, left
			Alignment(Center, Top),

			// Child 1: Fixed width, grow height
			Box(
				Sizing(
					Grow(),    // width
					Fixed(10), // height
				),
			),

			// Child 2: Grow width and height
			Box(
				Sizing(
					Grow(), // width
					Grow(), // height
				),
				Alignment(Right, Bottom),

				Box(
					Sizing(
						Fixed(20), // width
						Fixed(20), // height
					),
				),
			),

			// Child 3: Fixed width, grow height
			Box(
				Sizing(
					Grow(),
					Fixed(10),
				),
			),
		),
	)

	file, err := os.Create("layout.pdf")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	WritePDF(file, n)
}
