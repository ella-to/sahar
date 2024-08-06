package pdf

import (
	"fmt"
	"io"
	"strconv"

	"math/rand"

	"github.com/signintech/gopdf"

	"ella.to/sahar"
)

func getFonts(node *sahar.Node, result map[string]string) {
	fontFamilySrc, ok := node.Attributes["font-family-src"]
	if ok {
		fontFamily, ok := node.Attributes["font-family"]
		if ok {
			result[fontFamily.(string)] = fontFamilySrc.(string)
		}
	}
	for _, child := range node.Children {
		getFonts(child, result)
	}
}

func Write(w io.Writer, nodes ...*sahar.Node) error {
	pdf := gopdf.GoPdf{}

	for i, node := range nodes {
		if i == 0 {
			pageSize := *gopdf.PageSizeA4
			pageSize.H = float64(node.Height)
			pageSize.W = float64(node.Width)

			pdf.Start(gopdf.Config{PageSize: pageSize})
		}
		pdf.AddPage()

		fontsMap := make(map[string]string)
		getFonts(node, fontsMap)

		for name, src := range fontsMap {
			err := pdf.AddTTFFont(name, src)
			if err != nil {
				return err
			}
		}

		for _, child := range node.Children {
			write(&pdf, child)
		}
	}

	_, err := pdf.WriteTo(w)
	return err
}

func write(pdf *gopdf.GoPdf, node *sahar.Node) {
	//pdf.SetFillColor(generateRandomColor()) //setup fill color

	if sahar.IsType(node, sahar.TextType) {
		fontName := node.Attributes["font-family"].(string)
		fontSize := node.Attributes["font-size"].(float64)
		text := node.Attributes["text"].(string)
		color := getColor(node)

		if color != "" {
			r, g, b, err := hexToRGB(color)
			if err == nil {
				pdf.SetFillColor(r, g, b)
			}
		}

		pdf.SetFont(fontName, "", fontSize)
		pdf.SetXY(node.X, node.Y)
		pdf.Cell(nil, text)

		return
	} else if sahar.IsType(node, sahar.ImageType) {
		imagePath := node.Attributes["img-src"].(string)

		pdf.Image(imagePath, node.X, node.Y, &gopdf.Rect{W: node.Width, H: node.Height})

		return
	}

	borderWidth, borderColor := getBorder(node)
	if borderWidth > 0 {
		pdf.SetLineWidth(borderWidth)
	}

	pdf.SetLineWidth(0.0)
	if borderColor != "" {
		r, g, b, err := hexToRGB(borderColor)
		if err == nil {
			pdf.SetStrokeColor(r, g, b)
		}
	}

	backgroundColor := getBackgroundColor(node)
	if backgroundColor != "" {
		r, g, b, err := hexToRGB(backgroundColor)
		if err == nil {
			pdf.SetFillColor(r, g, b)
		}
	}

	pdf.RectFromUpperLeftWithStyle(node.X, node.Y, node.Width, node.Height, "FD")

	for _, child := range node.Children {
		write(pdf, child)
	}
}

func getColor(node *sahar.Node) string {
	if val, ok := node.Attributes["color"]; ok {
		return val.(string)
	}
	return ""
}

func getBackgroundColor(node *sahar.Node) string {
	if val, ok := node.Attributes["background-color"]; ok {
		return val.(string)
	}
	return ""
}

func getBorder(node *sahar.Node) (float64, string) {
	width, ok := node.Attributes["border-width"]
	if !ok {
		return 0, ""
	}

	color, ok := node.Attributes["border-color"]
	if !ok {
		return 0, ""
	}

	return width.(float64), color.(string)
}

func hexToRGB(hex string) (uint8, uint8, uint8, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, fmt.Errorf("invalid hex color format")
	}

	r, err := strconv.ParseUint(hex[1:3], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	g, err := strconv.ParseUint(hex[3:5], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	b, err := strconv.ParseUint(hex[5:7], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}

	return uint8(r), uint8(g), uint8(b), nil
}

func generateRandomColor() (uint8, uint8, uint8) {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return uint8(r), uint8(g), uint8(b)
}
