# Sahar Layout Engine

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/ella.to/sahar.svg)](https://pkg.go.dev/ella.to/sahar)
[![Go Report Card](https://goreportcard.com/badge/ella.to/sahar)](https://goreportcard.com/report/ella.to/sahar)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

_A powerful, declarative layout engine for Go that enables you to create professional PDF documents with ease_

</div>

## 🚀 Overview

Sahar is a modern layout engine designed for Go developers who need to generate complex PDF documents programmatically. Built with a component-based architecture, it provides an intuitive API for creating everything from simple reports to sophisticated document layouts.

### ✨ Key Features

- **🎯 Declarative Syntax** - Build layouts using composable, reusable components
- **📐 Flexible Sizing** - Support for fixed, fit-to-content, and grow sizing modes
- **🎨 Rich Typography** - Complete font control with size, family, color, and line height
- **📦 Smart Layout** - Automatic positioning with precise alignment control
- **🎭 Visual Styling** - Borders, backgrounds, padding, and spacing
- **🏗️ Nested Structures** - Create complex hierarchical document layouts
- **📄 PDF Export** - Direct rendering to PDF format
- **🔧 Type Safety** - Fully typed Go API with compile-time safety

## 📦 Installation

```bash
go get ella.to/sahar
```

## 🏃 Quick Start

Here's a simple example that creates a business card:

```go
package main

import (
    "os"
    "ella.to/sahar"
)

func main() {
    // Load fonts (required for text rendering)
    err := sahar.LoadFonts("Arial", "./Arial.ttf")
    if err != nil {
        panic(err)
    }

    // Create a business card layout
    card := sahar.Layout(
        sahar.Box(
            sahar.Sizing(sahar.Fixed(350), sahar.Fixed(200)),
            sahar.Direction(sahar.TopToBottom),
            sahar.Padding(20, 20, 20, 20),
            sahar.BackgroundColor("#ffffff"),
            sahar.Border(1),
            sahar.BorderColor("#e0e0e0"),

            // Header with logo and company info
            sahar.Box(
                sahar.Direction(sahar.LeftToRight),
                sahar.Alignment(sahar.Left, sahar.Middle),
                sahar.ChildGap(15),

                sahar.Image("./logo.png",
                    sahar.Sizing(sahar.Fixed(40), sahar.Fixed(40)),
                ),

                sahar.Box(
                    sahar.Direction(sahar.TopToBottom),
                    sahar.ChildGap(5),

                    sahar.Text("Tech Solutions Inc.",
                        sahar.FontType("Arial"),
                        sahar.FontSize(16),
                        sahar.FontColor("#2c3e50"),
                    ),

                    sahar.Text("Innovation at Scale",
                        sahar.FontType("Arial"),
                        sahar.FontSize(10),
                        sahar.FontColor("#7f8c8d"),
                    ),
                ),
            ),

            // Contact information
            sahar.Box(
                sahar.Sizing(sahar.Grow(), sahar.Grow()),
                sahar.Direction(sahar.TopToBottom),
                sahar.Alignment(sahar.Right, sahar.Bottom),
                sahar.ChildGap(3),

                sahar.Text("John Doe",
                    sahar.FontType("Arial"),
                    sahar.FontSize(14),
                    sahar.FontColor("#2c3e50"),
                ),
                sahar.Text("Senior Developer",
                    sahar.FontType("Arial"),
                    sahar.FontSize(11),
                    sahar.FontColor("#34495e"),
                ),
                sahar.Text("john.doe@techsolutions.com",
                    sahar.FontType("Arial"),
                    sahar.FontSize(10),
                    sahar.FontColor("#7f8c8d"),
                ),
            ),
        ),
    )

    // Export to PDF
    file, err := os.Create("business_card.pdf")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    err = sahar.RenderToPDF(file, card)
    if err != nil {
        panic(err)
    }
}
```

## 🎯 Core Concepts

### Node Types

Sahar uses three fundamental node types to build layouts:

| Type      | Purpose                      | Use Cases                           |
| --------- | ---------------------------- | ----------------------------------- |
| **Box**   | Container for other nodes    | Sections, panels, layout containers |
| **Text**  | Text content with typography | Headings, paragraphs, labels        |
| **Image** | Image content                | Logos, photos, charts, diagrams     |

### Sizing System

The flexible sizing system adapts to different layout needs:

| Type      | Behavior             | Best For                              |
| --------- | -------------------- | ------------------------------------- |
| **Fixed** | Exact dimensions     | Known sizes, consistent layouts       |
| **Fit**   | Size to content      | Dynamic content, responsive design    |
| **Grow**  | Fill available space | Flexible sections, responsive layouts |

### Layout Directions

Control how child elements are arranged:

- **LeftToRight** - Horizontal arrangement (row)
- **TopToBottom** - Vertical arrangement (column)

### Alignment Options

Precise control over element positioning:

| Horizontal | Vertical | Description          |
| ---------- | -------- | -------------------- |
| `Left`     | `Top`    | Align to start edges |
| `Center`   | `Middle` | Center alignment     |
| `Right`    | `Bottom` | Align to end edges   |

## 📚 Examples & Use Cases

### 1. Invoice Layout

```go
func CreateInvoice() *sahar.Node {
    return sahar.Box(
        sahar.Sizing(sahar.A4()...),
        sahar.Direction(sahar.TopToBottom),
        sahar.Padding(40, 40, 40, 40),

        // Header
        sahar.Box(
            sahar.Direction(sahar.LeftToRight),
            sahar.Alignment(sahar.Left, sahar.Top),
            sahar.ChildGap(20),

            sahar.Box(
                sahar.Direction(sahar.TopToBottom),
                sahar.Text("INVOICE",
                    sahar.FontSize(24),
                    sahar.FontColor("#2c3e50"),
                ),
                sahar.Text("Invoice #INV-001",
                    sahar.FontSize(12),
                    sahar.FontColor("#7f8c8d"),
                ),
            ),

            sahar.Box(
                sahar.Sizing(sahar.Grow(), sahar.Fit()),
                sahar.Alignment(sahar.Right, sahar.Top),
                sahar.Text("Due: December 31, 2023",
                    sahar.FontSize(12),
                    sahar.FontColor("#e74c3c"),
                ),
            ),
        ),

        // Items table would go here...
    )
}
```

### 2. Report with Charts

```go
func CreateReport() *sahar.Node {
    return sahar.Box(
        sahar.Sizing(sahar.USLetter()...),
        sahar.Direction(sahar.TopToBottom),
        sahar.ChildGap(30),

        // Title
        sahar.Text("Quarterly Report Q4 2023",
            sahar.FontSize(20),
            sahar.FontColor("#2c3e50"),
            sahar.Alignment(sahar.Center, sahar.Top),
        ),

        // Content grid
        sahar.Box(
            sahar.Direction(sahar.LeftToRight),
            sahar.ChildGap(20),

            // Left column - text content
            sahar.Box(
                sahar.Sizing(sahar.Fixed(300), sahar.Fit()),
                sahar.Direction(sahar.TopToBottom),
                sahar.ChildGap(15),

                sahar.Text("Executive Summary",
                    sahar.FontSize(16),
                    sahar.FontColor("#34495e"),
                ),
                sahar.Text("Revenue increased by 23% compared to Q3...",
                    sahar.FontSize(11),
                    sahar.FontColor("#2c3e50"),
                ),
            ),

            // Right column - chart
            sahar.Image("./chart.png",
                sahar.Sizing(sahar.Grow(), sahar.Fixed(200)),
            ),
        ),
    )
}
```

### 3. Form Layout

```go
func CreateForm() *sahar.Node {
    return sahar.Box(
        sahar.Direction(sahar.TopToBottom),
        sahar.ChildGap(15),
        sahar.Padding(30, 30, 30, 30),

        formField("Full Name:", "John Doe"),
        formField("Email:", "john@example.com"),
        formField("Phone:", "+1 (555) 123-4567"),
    )
}

func formField(label, value string) *sahar.Node {
    return sahar.Box(
        sahar.Direction(sahar.LeftToRight),
        sahar.ChildGap(10),

        sahar.Text(label,
            sahar.Sizing(sahar.Fixed(100), sahar.Fit()),
            sahar.FontSize(12),
            sahar.FontColor("#2c3e50"),
        ),
        sahar.Box(
            sahar.Sizing(sahar.Grow(), sahar.Fit()),
            sahar.Border(1),
            sahar.BorderColor("#bdc3c7"),
            sahar.Padding(8, 10, 8, 10),

            sahar.Text(value,
                sahar.FontSize(12),
                sahar.FontColor("#34495e"),
            ),
        ),
    )
}
```

## 🎨 Styling Guide

### Typography

```go
// Font styling options
sahar.Text("Styled Text",
    sahar.FontSize(16),              // Size in points
    sahar.FontType("Arial"),         // Font family
    sahar.FontColor("#2c3e50"),      // Hex color
)
```

### Spacing & Layout

```go
// Layout control
sahar.Box(
    sahar.Padding(10, 15, 10, 15),      // Top, Right, Bottom, Left
    sahar.ChildGap(20),                 // Space between children
    sahar.Direction(sahar.TopToBottom), // Layout direction
    sahar.Alignment(sahar.Center, sahar.Middle),
)
```

### Visual Design

```go
// Visual styling
sahar.Box(
    sahar.BackgroundColor("#f8f9fa"),   // Background color
    sahar.Border(2),                    // Border width
    sahar.BorderColor("#dee2e6"),       // Border color
)
```

### Advanced Sizing

```go
// Flexible sizing with constraints
sahar.Box(
    sahar.Sizing(
        sahar.Fixed(300),               // Fixed width
        sahar.Fit(                      // Height fits content
            sahar.Min(100),             // Minimum height
            sahar.Max(500),             // Maximum height
        ),
    ),
)
```

### Page Presets

```go
// Standard page sizes
sahar.Box(sahar.Sizing(sahar.A4()...))        // A4: 595.28 x 841.89 points
sahar.Box(sahar.Sizing(sahar.USLetter()...))  // US Letter: 612 x 792 points
sahar.Box(sahar.Sizing(sahar.USLegal()...))   // US Legal: 612 x 1008 points
```

## 📖 API Reference

### Core Functions

| Function   | Signature                         | Description                   |
| ---------- | --------------------------------- | ----------------------------- |
| `Box()`    | `Box(...nodeOpt) *Node`           | Creates a container node      |
| `Text()`   | `Text(string, ...textOpt) *Node`  | Creates a text node           |
| `Image()`  | `Image(string, ...nodeOpt) *Node` | Creates an image node         |
| `Layout()` | `Layout(*Node) *Node`             | Processes layout calculations |

### Sizing Functions

| Function  | Parameters  | Description             |
| --------- | ----------- | ----------------------- |
| `Fixed()` | `float64`   | Sets exact dimensions   |
| `Fit()`   | `...fitOpt` | Size to fit content     |
| `Grow()`  | -           | Expand to fill space    |
| `Min()`   | `float64`   | Sets minimum constraint |
| `Max()`   | `float64`   | Sets maximum constraint |

### Layout Options

| Function      | Parameters                 | Description                   |
| ------------- | -------------------------- | ----------------------------- |
| `Direction()` | `direction`                | Sets layout direction         |
| `Alignment()` | `Horizontal, Vertical`     | Sets alignment                |
| `Padding()`   | `top, right, bottom, left` | Sets internal spacing         |
| `ChildGap()`  | `float64`                  | Sets spacing between children |

### Typography

| Function      | Parameters | Description              |
| ------------- | ---------- | ------------------------ |
| `FontSize()`  | `float64`  | Sets font size in points |
| `FontType()`  | `string`   | Sets font family         |
| `FontColor()` | `string`   | Sets text color (hex)    |

### Visual Styling

| Function            | Parameters | Description                 |
| ------------------- | ---------- | --------------------------- |
| `BackgroundColor()` | `string`   | Sets background color (hex) |
| `Border()`          | `float64`  | Sets border width           |
| `BorderColor()`     | `string`   | Sets border color (hex)     |

### PDF Generation

| Function        | Parameters            | Description                         |
| --------------- | --------------------- | ----------------------------------- |
| `LoadFonts()`   | `...string`           | Loads font files (name, path pairs) |
| `RenderToPDF()` | `io.Writer, ...*Node` | Renders nodes to PDF                |

## 🔧 Advanced Usage

### Multi-Page Documents

```go
func CreateMultiPageDocument() {
    page1 := sahar.Layout(createCoverPage())
    page2 := sahar.Layout(createContentPage())
    page3 := sahar.Layout(createSummaryPage())

    file, _ := os.Create("document.pdf")
    defer file.Close()

    sahar.RenderToPDF(file, page1, page2, page3)
}
```

### Dynamic Content

```go
func CreateDynamicList(items []string) *sahar.Node {
    children := make([]*sahar.Node, len(items))
    for i, item := range items {
        children[i] = sahar.Text(item, sahar.FontSize(12))
    }

    return sahar.Box(
        sahar.Direction(sahar.TopToBottom),
        sahar.ChildGap(5),
        sahar.Children(children...),
    )
}
```

### Component Composition

```go
func Header(title string) *sahar.Node {
    return sahar.Box(
        sahar.Direction(sahar.LeftToRight),
        sahar.Padding(0, 0, 20, 0),
        sahar.Text(title, sahar.FontSize(18)),
    )
}

func Section(title, content string) *sahar.Node {
    return sahar.Box(
        sahar.Direction(sahar.TopToBottom),
        sahar.ChildGap(10),
        Header(title),
        sahar.Text(content, sahar.FontSize(12)),
    )
}
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
