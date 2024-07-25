package pdf

import (
	"io"

	"math/rand"

	"github.com/signintech/gopdf"

	"ella.to/sahar"
)

func Write(w io.Writer, node *sahar.Node) error {
	pdf := gopdf.GoPdf{}

	pageSize := *gopdf.PageSizeA4
	pageSize.H = float64(node.Height)
	pageSize.W = float64(node.Width)

	pdf.Start(gopdf.Config{PageSize: pageSize})
	pdf.AddPage()

	for _, child := range node.Children {
		write(&pdf, child)
	}

	_, err := pdf.WriteTo(w)
	return err
}

func write(pdf *gopdf.GoPdf, node *sahar.Node) {
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
