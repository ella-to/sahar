package main

import (
	"os"

	"ella.to/sahar"
)

func Header(companyName string) *sahar.Node {
	return sahar.Box(
		sahar.Direction(sahar.LeftToRight),
		sahar.Alignment(sahar.Left, sahar.Middle),
		sahar.BorderColor("#ffffff"),
		sahar.ChildGap(10),

		sahar.Image(
			"./logo",
			sahar.Sizing(sahar.Fixed(50), sahar.Fixed(50)),
		),

		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.BorderColor("#ffffff"),
			sahar.ChildGap(2),

			// Children
			sahar.Text(
				companyName,
				sahar.FontType("Arial"),
				sahar.FontSize(18),
				sahar.FontColor("#5478ac"),
			),

			sahar.Text(
				"Compnay Message",
				sahar.FontType("Arial"),
				sahar.FontSize(12),
				sahar.FontColor("#717171"),
			),
		),
	)
}

func Main() *sahar.Node {
	return sahar.Box(
		// sahar.BorderColor("#ffffff"),
		sahar.Sizing(sahar.Grow(), sahar.Grow()),
		sahar.Alignment(sahar.Right, sahar.Bottom),

		sahar.Box(
			// sahar.BorderColor("#ffffff"),
			sahar.Direction(sahar.TopToBottom),
			sahar.Alignment(sahar.Right, sahar.Middle),
			sahar.ChildGap(2),

			sahar.Text(
				"Full Name",
				sahar.FontType("Arial"),
				sahar.FontSize(16),
				sahar.FontColor("#5478ac"),
			),
			sahar.Text(
				"Job Title",
				sahar.FontType("Arial"),
				sahar.FontSize(12),
				sahar.FontColor("#717171"),
			),
			sahar.Text(
				"Email / Other",
				sahar.FontType("Arial"),
				sahar.FontSize(12),
				sahar.FontColor("#717171"),
			),
		),
	)
}

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

	page1 := sahar.Layout(
		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.Fixed(300), sahar.Fixed(150)),
			sahar.Padding(10, 10, 10, 10),

			Header("Compnay B"),
			Main(),
		),
	)

	page2 := sahar.Layout(
		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.Fixed(300), sahar.Fixed(150)),
			sahar.Padding(10, 10, 10, 10),

			Header("Compnay A"),
			Main(),
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

	err = sahar.RenderToPDF(pdfFile, page1, page2)
	if err != nil {
		panic(err)
	}
}
