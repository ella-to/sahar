# Sahar PDF Layout Library Specification

## Overview

Sahar is a Go library for creating PDF documents using a declarative, tree-based layout system. It uses a flexbox-like model with `Box`, `Text`, and `Image` nodes.

## Quick Reference

```go
import "ella.to/sahar"

// 1. Load fonts
sahar.LoadFonts("FontName", "./path/to/font.ttf")

// 2. Build layout tree
page := sahar.Layout(sahar.Box(...))

// 3. Render to PDF
sahar.RenderToPDF(writer, page)
```

## Node Types

| Type | Constructor | Purpose |
|------|-------------|---------|
| `Box` | `sahar.Box(opts...)` | Container for layout and grouping |
| `Text` | `sahar.Text(value, opts...)` | Display text |
| `Image` | `sahar.Image(path, opts...)` | Display image (PNG, JPG, GIF) |

## Box Options

```go
sahar.Box(
    // SIZE - Choose one combination:
    sahar.Sizing(sahar.Fixed(width), sahar.Fixed(height)),  // Exact size in points
    sahar.Sizing(sahar.Grow(), sahar.Grow()),               // Fill available space
    sahar.Sizing(sahar.Fit(), sahar.Fit()),                 // Shrink to content (default)
    sahar.Sizing(sahar.A4()...),                            // A4 page: 595.28 x 841.89 pt
    sahar.Sizing(sahar.USLetter()...),                      // US Letter: 612 x 792 pt
    
    // DIRECTION - How children are arranged:
    sahar.Direction(sahar.LeftToRight),  // Horizontal (default)
    sahar.Direction(sahar.TopToBottom),  // Vertical
    
    // ALIGNMENT - Position children within box:
    sahar.Alignment(sahar.Left, sahar.Top),       // Horizontal: Left|Center|Right
    sahar.Alignment(sahar.Center, sahar.Middle),  // Vertical: Top|Middle|Bottom
    sahar.Alignment(sahar.Right, sahar.Bottom),
    
    // SPACING:
    sahar.Padding(top, right, bottom, left),  // Inner spacing (points)
    sahar.ChildGap(gap),                       // Space between children (points)
    
    // VISUAL:
    sahar.BackgroundColor("#RRGGBB"),  // Hex color
    sahar.Border(width),               // Border width in points
    sahar.BorderColor("#RRGGBB"),      // Hex color
    
    // CHILDREN - Nested nodes:
    sahar.Box(...),
    sahar.Text(...),
    sahar.Image(...),
)
```

## Text Options

```go
sahar.Text("content",
    sahar.FontType("Arial"),     // Font name (must be loaded)
    sahar.FontSize(12),          // Size in points
    sahar.FontColor("#RRGGBB"),  // Hex color
    sahar.Border(1),             // Debug border
)
```

## Image Options

```go
sahar.Image("./path/to/image.png",
    sahar.Sizing(sahar.Fixed(width), sahar.Fixed(height)),  // Required for images
    sahar.Border(1),  // Optional border
)
```

## Sizing Reference

| Function | Behavior |
|----------|----------|
| `Fixed(n)` | Exact size of `n` points |
| `Grow()` | Expand to fill remaining space |
| `Fit()` | Shrink to fit content |
| `Fit(Min(n))` | Fit content, minimum `n` points |
| `Fit(Max(n))` | Fit content, maximum `n` points |

## Page Sizes (in points)

| Preset | Width | Height |
|--------|-------|--------|
| `A4()` | 595.28 | 841.89 |
| `USLetter()` | 612 | 792 |
| `USLegal()` | 612 | 1008 |

## Layout Patterns

### Full Page with Header/Content/Footer

```go
sahar.Box(
    sahar.Direction(sahar.TopToBottom),
    sahar.Sizing(sahar.A4()...),
    sahar.Padding(40, 40, 40, 40),
    
    // Header - fixed height
    sahar.Box(
        sahar.Sizing(sahar.Grow(), sahar.Fixed(60)),
        sahar.Alignment(sahar.Left, sahar.Middle),
        sahar.Text("Header", sahar.FontSize(24)),
    ),
    
    // Content - grows to fill
    sahar.Box(
        sahar.Sizing(sahar.Grow(), sahar.Grow()),
        // ... content
    ),
    
    // Footer - fixed height
    sahar.Box(
        sahar.Sizing(sahar.Grow(), sahar.Fixed(40)),
        sahar.Alignment(sahar.Center, sahar.Middle),
        sahar.Text("Page 1", sahar.FontSize(10)),
    ),
)
```

### Horizontal Row with Columns

```go
sahar.Box(
    sahar.Direction(sahar.LeftToRight),
    sahar.ChildGap(20),
    
    // Left column - fixed width
    sahar.Box(
        sahar.Sizing(sahar.Fixed(200), sahar.Grow()),
        // ... sidebar content
    ),
    
    // Right column - grows
    sahar.Box(
        sahar.Sizing(sahar.Grow(), sahar.Grow()),
        // ... main content
    ),
)
```

### Centered Content

```go
sahar.Box(
    sahar.Sizing(sahar.Fixed(400), sahar.Fixed(300)),
    sahar.Alignment(sahar.Center, sahar.Middle),
    
    sahar.Text("Centered", sahar.FontSize(20)),
)
```

### Logo with Text (Business Card Style)

```go
sahar.Box(
    sahar.Direction(sahar.LeftToRight),
    sahar.Alignment(sahar.Left, sahar.Middle),
    sahar.ChildGap(10),
    
    sahar.Image("./logo.png",
        sahar.Sizing(sahar.Fixed(50), sahar.Fixed(50)),
    ),
    
    sahar.Box(
        sahar.Direction(sahar.TopToBottom),
        sahar.ChildGap(4),
        
        sahar.Text("Company Name",
            sahar.FontType("Arial"),
            sahar.FontSize(18),
            sahar.FontColor("#333333"),
        ),
        sahar.Text("Tagline",
            sahar.FontType("Arial"),
            sahar.FontSize(12),
            sahar.FontColor("#666666"),
        ),
    ),
)
```

### Grid Layout (2x2)

```go
sahar.Box(
    sahar.Direction(sahar.TopToBottom),
    sahar.ChildGap(10),
    
    // Row 1
    sahar.Box(
        sahar.Direction(sahar.LeftToRight),
        sahar.ChildGap(10),
        sahar.Box(sahar.Sizing(sahar.Grow(), sahar.Fixed(100)), /* cell 1 */),
        sahar.Box(sahar.Sizing(sahar.Grow(), sahar.Fixed(100)), /* cell 2 */),
    ),
    
    // Row 2
    sahar.Box(
        sahar.Direction(sahar.LeftToRight),
        sahar.ChildGap(10),
        sahar.Box(sahar.Sizing(sahar.Grow(), sahar.Fixed(100)), /* cell 3 */),
        sahar.Box(sahar.Sizing(sahar.Grow(), sahar.Fixed(100)), /* cell 4 */),
    ),
)
```

## Complete Example

```go
package main

import (
    "os"
    "ella.to/sahar"
)

func main() {
    // Load fonts
    sahar.LoadFonts("Arial", "./Arial.ttf")

    // Build page
    page := sahar.Layout(
        sahar.Box(
            sahar.Direction(sahar.TopToBottom),
            sahar.Sizing(sahar.A4()...),
            sahar.Padding(50, 50, 50, 50),
            sahar.BackgroundColor("#FFFFFF"),

            // Header
            sahar.Box(
                sahar.Direction(sahar.LeftToRight),
                sahar.Alignment(sahar.Left, sahar.Middle),
                sahar.ChildGap(15),

                sahar.Image("./logo.png",
                    sahar.Sizing(sahar.Fixed(60), sahar.Fixed(60)),
                ),
                sahar.Text("Document Title",
                    sahar.FontType("Arial"),
                    sahar.FontSize(24),
                    sahar.FontColor("#000000"),
                ),
            ),

            // Content area
            sahar.Box(
                sahar.Sizing(sahar.Grow(), sahar.Grow()),
                sahar.Direction(sahar.TopToBottom),
                sahar.ChildGap(20),
                sahar.Padding(20, 0, 20, 0),

                sahar.Text("Section 1",
                    sahar.FontType("Arial"),
                    sahar.FontSize(16),
                    sahar.FontColor("#333333"),
                ),
                sahar.Text("Body text goes here...",
                    sahar.FontType("Arial"),
                    sahar.FontSize(12),
                    sahar.FontColor("#666666"),
                ),
            ),

            // Footer
            sahar.Box(
                sahar.Sizing(sahar.Grow(), sahar.Fixed(30)),
                sahar.Alignment(sahar.Center, sahar.Bottom),

                sahar.Text("Page 1",
                    sahar.FontType("Arial"),
                    sahar.FontSize(10),
                    sahar.FontColor("#999999"),
                ),
            ),
        ),
    )

    // Render
    file, _ := os.Create("output.pdf")
    defer file.Close()
    sahar.RenderToPDF(file, page)
}
```

## Rules

1. **Always call `Layout()` before `RenderToPDF()`**
2. **Load fonts before using them in Text nodes**
3. **Images require explicit `Sizing` with `Fixed` dimensions**
4. **Colors must be hex format: `#RRGGBB`**
5. **All measurements are in points (1 inch = 72 points)**
6. **Nest children directly inside `Box()` constructor**
7. **Use `Direction(TopToBottom)` for vertical stacking**
8. **Use `Direction(LeftToRight)` for horizontal layout (default)**
9. **`Grow()` only works when parent has defined size**
10. **Multiple pages: pass multiple nodes to `RenderToPDF(writer, page1, page2, ...)`**

## Image to Code Translation Guide

When converting a visual PDF design to code:

1. **Identify the page size** → Use `Sizing(A4()...)` or `Sizing(Fixed(w), Fixed(h))`

2. **Identify major sections** (header, body, footer) → Each becomes a `Box` with `Direction(TopToBottom)`

3. **For each section, determine:**
   - Layout direction: horizontal elements → `LeftToRight`, stacked elements → `TopToBottom`
   - Size: fixed dimensions → `Fixed(n)`, fill space → `Grow()`, fit content → `Fit()`
   - Alignment: where content sits within the box

4. **For text elements:**
   - Extract: content, font size, color
   - Estimate sizes: headline ~18-24pt, body ~10-12pt, small ~8-10pt

5. **For images:**
   - Note position and approximate dimensions
   - Always use `Fixed` sizing

6. **Spacing:**
   - Gaps between elements → `ChildGap(n)`
   - Space around content → `Padding(t, r, b, l)`

7. **Build tree inside-out:** Start with leaf nodes (Text, Image), wrap in Boxes, combine into sections, wrap in page Box.
