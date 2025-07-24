package sahar

import (
	"io"

	"codeberg.org/go-pdf/fpdf"
)

func writePdf(pdf *fpdf.Fpdf, node *Node) {
	pdf.Rect(
		node.Position.X,
		node.Position.Y,
		node.Width.Value,
		node.Height.Value,
		"D",
	)

	for _, child := range node.Children {
		writePdf(pdf, child)
	}
}

func WritePDF(w io.Writer, node *Node) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	writePdf(pdf, node)

	pdf.Output(w)
}
