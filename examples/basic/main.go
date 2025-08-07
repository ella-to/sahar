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
			// sahar.Sizing(
			// 	sahar.A4()...,
			// ),
			sahar.Alignment(sahar.Center, sahar.Middle),
			sahar.Direction(sahar.TopToBottom),
			// sahar.BackgroundColor("#ff0000"),

			// sahar.Box(
			// 	sahar.Border(1),
			// 	sahar.Alignment(sahar.Center, sahar.Middle),
			// 	sahar.Sizing(
			// 		sahar.Grow(),
			// 		sahar.Fixed(100),
			// 	),

			sahar.Image(
				"./logo",
				sahar.Border(0),
				sahar.BorderColor("#ff0000"),
				sahar.Sizing(
					sahar.Fixed(100),
					sahar.Fixed(100),
				),
			),
			// ),

			sahar.Text(
				"123 Hello World!",
				sahar.FontType("Arial"),
				sahar.FontSize(20),
				sahar.FontColor("#ff0000"),

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

	err = sahar.RenderToPDF(pdfFile, root)
	if err != nil {
		panic(err)
	}
}
