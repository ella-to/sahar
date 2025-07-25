package main

import (
	"os"

	"ella.to/sahar"
)

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
			// for debugging purposes border can be set
			// sahar.Border(1),
			sahar.Sizing(
				sahar.A4()...,
			),
			sahar.Alignment(sahar.Center, sahar.Middle),

			sahar.Text(
				"123 Hello World!",
				sahar.FontType("Arial"),
				sahar.FontSize(20),
				// for debugging purposes border can be set
				// sahar.Border(1),
			),
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

	err = sahar.RenderToPDF(root, pdfFile)
	if err != nil {
		panic(err)
	}
}
