package sahar

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"codeberg.org/go-pdf/fpdf"
)

// RenderToPDF renders the node tree to a PDF and writes it to the provided writer
func RenderToPDF(writer io.Writer, nodes ...*Node) error {
	if len(nodes) == 0 {
		return fmt.Errorf("there is no node to render")
	}

	// Create a new PDF document
	pdf := fpdf.New("P", "pt", "A4", "")

	for _, node := range nodes {
		pdf.AddPageFormat("P", fpdf.SizeType{
			Wd: node.Width.Value,
			Ht: node.Height.Value,
		})

		// Set default font if not already set
		pdf.SetFont("Arial", "", 12)

		// Render the node tree
		if err := renderNode(pdf, node); err != nil {
			return fmt.Errorf("failed to render node: %w", err)
		}
	}

	// Write PDF to the writer
	return pdf.Output(writer)
}

// renderNode recursively renders a node and its children
func renderNode(pdf *fpdf.Fpdf, node *Node) error {
	if node == nil {
		return nil
	}

	switch node.Type {
	case BoxType:
		if err := renderBox(pdf, node); err != nil {
			return err
		}
	case TextType:
		if err := renderText(pdf, node); err != nil {
			return err
		}
	case ImageType:
		if err := renderImage(pdf, node); err != nil {
			return err
		}
	}

	// Render children
	for _, child := range node.Children {
		if err := renderNode(pdf, child); err != nil {
			return err
		}
	}

	return nil
}

// renderBox renders a box node (draws a rectangle)
func renderBox(pdf *fpdf.Fpdf, node *Node) error {
	x := node.Position.X
	y := node.Position.Y
	width := node.Width.Value
	height := node.Height.Value

	// Set line width and color for the box border
	pdf.SetLineWidth(node.Border)

	r, g, b, err := hexToRGB(node.BorderColor, "#000000")
	if err != nil {
		return fmt.Errorf("invalid border color: %w", err)
	}
	pdf.SetDrawColor(r, g, b) // Set border color

	r, g, b, err = hexToRGB(node.BackgroundColor, "#FFFFFF")
	if err != nil {
		return fmt.Errorf("invalid background color: %w", err)
	}
	pdf.SetFillColor(r, g, b) // Set background color

	// Draw rectangle (border only, no fill for now)
	pdf.Rect(x, y, width, height, "FD")

	return nil
}

// renderText renders a text node
func renderText(pdf *fpdf.Fpdf, node *Node) error {
	if node.Value == "" {
		return nil
	}

	x := node.Position.X
	y := node.Position.Y
	width := node.Width.Value
	height := node.Height.Value

	if node.Border > 0 {
		renderBox(pdf, node)
	}

	// Set font if specified
	if node.FontType != "" && node.FontSize > 0 {
		// Map common font names to FPDF font names
		fontName := mapFontName(node.FontType)
		pdf.SetFont(fontName, "", node.FontSize)
	}

	r, g, b, err := hexToRGB(node.FontColor, "#000000")
	if err != nil {
		return fmt.Errorf("invalid font color: %w", err)
	}

	// Set text color
	pdf.SetTextColor(r, g, b)

	// Handle multi-line text
	lines := strings.Split(node.Value, "\n")

	// Get line height
	_, lineHeight := pdf.GetFontSize()
	lineSpacing := lineHeight * 1.2 // Add some line spacing

	// Calculate starting Y position based on vertical alignment
	totalTextHeight := float64(len(lines)) * lineSpacing
	var startY float64

	switch node.Vertical {
	case Top:
		startY = y + lineHeight // Start from top + font height
	case Middle:
		startY = y + (height-totalTextHeight)/2 + lineHeight
	case Bottom:
		startY = y + height - totalTextHeight + lineHeight
	default:
		startY = y + lineHeight
	}

	// Render each line
	for i, line := range lines {
		if line == "" {
			continue
		}

		lineY := startY + float64(i)*lineSpacing

		// Calculate X position based on horizontal alignment
		var lineX float64
		textWidth := pdf.GetStringWidth(line)

		switch node.Horizontal {
		case Left:
			lineX = x + node.Padding[3] // Add left padding
		case Center:
			lineX = x + (width-textWidth)/2
		case Right:
			lineX = x + width - textWidth - node.Padding[1] // Subtract right padding
		default:
			lineX = x + node.Padding[3]
		}

		// Ensure text doesn't go outside the node bounds
		if lineX < x {
			lineX = x
		}
		if lineX+textWidth > x+width {
			lineX = x + width - textWidth
		}

		// Draw the text
		pdf.Text(lineX, lineY, line)
	}

	return nil
}

// renderImage renders an image node (placeholder implementation)
func renderImage(pdf *fpdf.Fpdf, node *Node) error {
	if node.Border > 0 {
		renderBox(pdf, node)
	}

	x := node.Position.X
	y := node.Position.Y
	width := node.Width.Value
	height := node.Height.Value

	imageType, err := detectImageType(node.Value)
	if err != nil {
		return fmt.Errorf("failed to detect image type: %w", err)
	}

	pdf.ImageOptions(node.Value, x, y, width, height, false, fpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: imageType,
	}, 0, "")

	return nil
}

// mapFontName maps common font names to FPDF-compatible font names
func mapFontName(fontType string) string {
	switch strings.ToLower(fontType) {
	case "arial", "helvetica":
		return "Arial"
	case "times", "times new roman":
		return "Times"
	case "courier", "courier new":
		return "Courier"
	default:
		return "Arial" // Default fallback
	}
}

// RenderToPDFWithOptions renders the node tree to PDF with additional options
func RenderToPDFWithOptions(root *Node, writer io.Writer, options PDFOptions) error {
	if root == nil {
		return fmt.Errorf("root node cannot be nil")
	}

	// Create PDF with specified options
	orientation := "P"
	if options.Landscape {
		orientation = "L"
	}

	pageSize := "A4"
	if options.PageSize != "" {
		pageSize = options.PageSize
	}

	pdf := fpdf.New(orientation, "pt", pageSize, "")

	// Set margins if specified
	if options.MarginTop > 0 || options.MarginRight > 0 || options.MarginBottom > 0 || options.MarginLeft > 0 {
		pdf.SetMargins(options.MarginLeft, options.MarginTop, options.MarginRight)
		pdf.SetAutoPageBreak(true, options.MarginBottom)
	}

	pdf.AddPage()

	// Set default font
	defaultFont := "Arial"
	defaultSize := 12.0
	if options.DefaultFont != "" {
		defaultFont = options.DefaultFont
	}
	if options.DefaultFontSize > 0 {
		defaultSize = options.DefaultFontSize
	}
	pdf.SetFont(defaultFont, "", defaultSize)

	// Render the node tree
	if err := renderNode(pdf, root); err != nil {
		return fmt.Errorf("failed to render node: %w", err)
	}

	// Write PDF to the writer
	return pdf.Output(writer)
}

// PDFOptions contains options for PDF rendering
type PDFOptions struct {
	Landscape       bool
	PageSize        string // "A4", "A3", "Letter", etc.
	MarginTop       float64
	MarginRight     float64
	MarginBottom    float64
	MarginLeft      float64
	DefaultFont     string
	DefaultFontSize float64
}

// DefaultPDFOptions returns default PDF rendering options
func DefaultPDFOptions() PDFOptions {
	return PDFOptions{
		Landscape:       false,
		PageSize:        "A4",
		MarginTop:       72, // 1 inch
		MarginRight:     72, // 1 inch
		MarginBottom:    72, // 1 inch
		MarginLeft:      72, // 1 inch
		DefaultFont:     "Arial",
		DefaultFontSize: 12,
	}
}

func detectImageType(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read first 512 bytes
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		panic(err)
	}

	// Detect MIME type
	contentType := http.DetectContentType(buf)

	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		return "JPG", nil
	case strings.HasPrefix(contentType, "image/png"):
		return "PNG", nil
	case strings.HasPrefix(contentType, "image/gif"):
		return "GIF", nil
	default:
		return "", fmt.Errorf("unsupported image type: %s", contentType)
	}
}

func hexToRGB(hex string, defaultColor string) (r, g, b int, err error) {
	if hex == "" {
		hex = defaultColor
	}

	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		err = fmt.Errorf("invalid hex color: %s", hex)
		return
	}

	rPart, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return
	}
	gPart, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return
	}
	bPart, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return
	}

	return int(rPart), int(gPart), int(bPart), nil
}
