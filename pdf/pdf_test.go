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

		p1 := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,
			X:      5,
			Y:      5,
			Children: []*sahar.Node{
				{
					Width:  90,
					Height: 30,
					Order:  sahar.StackOrder,
					X:      5,
					Y:      5,
				},
				{
					Width:  90,
					Height: 30,
					Order:  sahar.StackOrder,
					X:      5,
					Y:      35,
				},
				{
					Width:  90,
					Height: 30,
					Order:  sahar.StackOrder,
					X:      5,
					Y:      65,
				},
			},
		}

		p2 := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,
			X:      5,
			Y:      5,
			Children: []*sahar.Node{
				{
					Width:  90,
					Height: 30,
					Order:  sahar.StackOrder,
					X:      5,
					Y:      5,
				},
				{
					Width:  90,
					Height: 30,
					Order:  sahar.StackOrder,
					X:      5,
					Y:      35,
				},
				{
					Width:  90,
					Height: 30,
					Order:  sahar.StackOrder,
					X:      5,
					Y:      65,
				},
			},
		}

		err = pdf.Write(file, p1, p2)
		if err != nil {
			t.Fatal(err)
		}
	})
}
