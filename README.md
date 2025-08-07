Sahar Layout Engine

A powerful, flexible layout engine for Go that enables you to create complex PDF documents using a declarative, component-based approach. Sahar provides an intuitive API for building layouts with precise control over positioning, sizing, and styling.

# Features

- Declarative Layout System: Build layouts using composable components
- Flexible Sizing: Support for fixed, fit-to-content, and grow sizing modes
- Advanced Alignment: Precise horizontal and vertical alignment control
- Rich Typography: Full font control with size, type, color, and line height
- Visual Styling: Borders, background colors, and padding support
- Nested Layouts: Create complex hierarchical structures
- PDF Generation: Direct output to PDF format

# Installation

```bash
    go get github.com/ella-to/sahar
```

# Core Concepts

### Node Types

| Type      | Description                              | Use Case                            |
| --------- | ---------------------------------------- | ----------------------------------- |
| BoxType   | Container node that can hold other nodes | Layout containers, sections, panels |
| TextType  | Text content with typography controls    | Headings, paragraphs, labels        |
| ImageType | Image content with source reference      | Photos, logos, diagrams             |

### Sizing Types

| Type      | Description         | Behavior                              |
| --------- | ------------------- | ------------------------------------- |
| FitType   | Size fits content   | Automatically adjusts to content size |
| FixedType | Fixed size value    | Uses exact specified dimensions       |
| GrowType  | Grows to fill space | Expands to use available space        |

### Alignment Options

| Horizontal | Vertical | Description                 |
| ---------- | -------- | --------------------------- |
| Left       | Top      | Align to left/top edges     |
| Center     | Middle   | Center alignment            |
| Right      | Bottom   | Align to right/bottom edges |

### Layout Directions

| Direction   | Description                    |
| ----------- | ------------------------------ |
| LeftToRight | Children arranged horizontally |
| TopToBottom | Children arranged vertically   |

# Quick Start

creating a simple business card

```golang
package main

import (
	"os"

	"ella.to/sahar"
)

func Header() *sahar.Node {
	return sahar.Box(
		sahar.Direction(sahar.LeftToRight),
		sahar.Alignment(sahar.Left, sahar.Middle),
		sahar.BorderColor("#ffffff"),
		sahar.ChildGap(10),

		sahar.Image(
			"./logo",
			sahar.Sizing(sahar.Fixed(50), sahar.Fixed(50)),
		),

		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.BorderColor("#ffffff"),
			sahar.ChildGap(2),

			// Children
			sahar.Text(
				"Company Name",
				sahar.FontType("Arial"),
				sahar.FontSize(18),
				sahar.FontColor("#5478ac"),
			),

			sahar.Text(
				"Compnay Message",
				sahar.FontType("Arial"),
				sahar.FontSize(12),
				sahar.FontColor("#717171"),
			),
		),
	)
}

func Main() *sahar.Node {
	return sahar.Box(
		sahar.BorderColor("#ffffff"),
		sahar.Sizing(sahar.Grow(), sahar.Grow()),
		sahar.Alignment(sahar.Right, sahar.Bottom),

		sahar.Box(
			sahar.BorderColor("#ffffff"),
			sahar.Direction(sahar.TopToBottom),
			sahar.Alignment(sahar.Right, sahar.Middle),
			sahar.ChildGap(2),

			sahar.Text(
				"Full Name",
				sahar.FontType("Arial"),
				sahar.FontSize(16),
				sahar.FontColor("#5478ac"),
			),
			sahar.Text(
				"Job Title",
				sahar.FontType("Arial"),
				sahar.FontSize(12),
				sahar.FontColor("#717171"),
			),
			sahar.Text(
				"Email / Other",
				sahar.FontType("Arial"),
				sahar.FontSize(12),
				sahar.FontColor("#717171"),
			),
		),
	)
}

func main() {
	//
	// load fonts
	//
	err := sahar.LoadFonts(
		"Arial", "./Arial.ttf",
	)
	if err != nil {
		panic(err)
	}

	root := sahar.Layout(
		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.Fixed(300), sahar.Fixed(150)),
			sahar.Padding(10, 10, 10, 10),

			Header(),
			Main(),
		),
	)

	//
	// Write the layout to a PDF file
	//

	pdfFile, err := os.Create("./layout.pdf")
	if err != nil {
		panic(err)
	}
	defer pdfFile.Close()

	err = sahar.RenderToPDF(pdfFile, root)
	if err != nil {
		panic(err)
	}
}
```

# Styling Options

### Typography

```golang
sahar.Text("Sample Text",
		sahar.FontSize(16),           // Font size in points
		sahar.FontType("Arial-Bold"), // Font family and weight
		sahar.FontColor("#2c3e50"),   // Hex color code
)
```

### Spacing and Layout

```golang
sahar.Box(
		sahar.Padding(10, 15, 10, 15),    // Top, Right, Bottom, Left
		sahar.ChildGap(20),               // Space between children
		sahar.Direction(sahar.TopToBottom), // Layout direction
)
```

### Visual Styling

```golang
sahar.Box(
		sahar.BackgroundColor("#f8f9fa"),  // Background color
		sahar.Border(2),                   // Border width
		sahar.BorderColor("#dee2e6"),      // Border color
)
```

### Advanced Sizing

```golang
sahar.Box(
		sahar.Sizing(
				sahar.Fixed(300),              // Fixed width: 300 points
				sahar.Fit(                     // Height fits content
						sahar.Min(100),            // Minimum height: 100 points
						sahar.Max(500),            // Maximum height: 500 points
				),
		),
)
```

# API Reference

### Core Functions

| Function | Parameters         | Description              |
| -------- | ------------------ | ------------------------ |
| Box()    | ...nodeOpt         | Creates a container node |
| Text()   | string, ...textOpt | Creates a text node      |
| Image()  | string, ...nodeOpt | Creates an image node    |

### Sizing Functions

| Function | Parameters | Description         |
| -------- | ---------- | ------------------- |
| Fixed()  | float64    | Sets fixed size     |
| Fit()    | ...fitOpt  | Size fits content   |
| Grow()   | -          | Grows to fill space |
| Min()    | float64    | Sets minimum size   |
| Max()    | float64    | Sets maximum size   |

### Layout Functions

| Function    | Parameters               | Description               |
| ----------- | ------------------------ | ------------------------- |
| Direction() | direction                | Sets layout direction     |
| Alignment() | Horizontal, Vertical     | Sets alignment            |
| Padding()   | top, right, bottom, left | Sets padding              |
| ChildGap()  | float64                  | Sets gap between children |

# License

This project is licensed under the MIT License - see the LICENSE file for details.
