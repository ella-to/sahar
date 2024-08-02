package pdf

import (
	"io"

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
	if sahar.IsType(node, sahar.TextType) {
		fontName := node.Attributes["font-family"].(string)
		fontSize := node.Attributes["font-size"].(float64)
		text := node.Attributes["text"].(string)

		pdf.SetFont(fontName, "", fontSize)
		pdf.SetXY(node.X, node.Y)
		pdf.Cell(nil, text)

		return
	} else if sahar.IsType(node, sahar.ImageType) {
		imagePath := node.Attributes["img-src"].(string)

		pdf.Image(imagePath, node.X, node.Y, &gopdf.Rect{W: node.Width, H: node.Height})

		return
	}

	drawRect(pdf, node.X, node.Y, node.Width, node.Height)

	for _, child := range node.Children {
		write(pdf, child)
	}
}

func drawRect(pdf *gopdf.GoPdf, x, y, width, height float64) {
	pdf.SetLineWidth(0.0)

	// randomize color
	pdf.SetFillColor(generateRandomColor()) //setup fill color
	pdf.RectFromUpperLeftWithStyle(x, y, width, height, "FD")
}

func generateRandomColor() (uint8, uint8, uint8) {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return uint8(r), uint8(g), uint8(b)
}
