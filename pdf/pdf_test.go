package pdf_test

import (
	"os"
	"testing"

	"ella.to/sahar"
	"ella.to/sahar/pdf"
)

func TestPdf(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		file, err := os.Create("./test.pdf")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		node := &sahar.Node{
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,
			Margin: [4]float64{5, 5, 5, 5},
			X:      5,
			Y:      5,
			Children: []*sahar.Node{
				{
					Width:  90,
					Height: 30,
					Type:   sahar.Stack,
					X:      5,
					Y:      5,
				},
				{
					Width:  90,
					Height: 30,
					Type:   sahar.Stack,
					X:      5,
					Y:      35,
				},
				{
					Width:  90,
					Height: 30,
					Type:   sahar.Stack,
					X:      5,
					Y:      65,
				},
			},
		}

		err = pdf.Write(file, node)
		if err != nil {
			t.Fatal(err)
		}
	})
}
