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

func Write(w io.Writer, node *sahar.Node) error {
	pdf := gopdf.GoPdf{}

	pageSize := *gopdf.PageSizeA4
	pageSize.H = float64(node.Height)
	pageSize.W = float64(node.Width)

	pdf.Start(gopdf.Config{PageSize: pageSize})
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

	background := getBackgroundColor(node)

	drawRect(pdf, node.X, node.Y, node.Width, node.Height, background)

	for _, child := range node.Children {
		write(pdf, child)
	}
}

func drawRect(pdf *gopdf.GoPdf, x, y, width, height float64, backgroundColor string) {
	pdf.SetLineWidth(0.0)

	if backgroundColor != "" {
		r, g, b, err := hexToRGB(backgroundColor)
		if err == nil {
			pdf.SetFillColor(r, g, b)
		}
	}

	pdf.RectFromUpperLeftWithStyle(x, y, width, height, "F")
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
