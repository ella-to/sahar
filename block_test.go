package sahar_test

import (
	"testing"

	"ella.to/sahar"
	"github.com/stretchr/testify/assert"
)

func TestBasicBlock(t *testing.T) {

	b := sahar.Stack(
		sahar.Padding(5, 5, 5, 5),

		sahar.FontFamily("Arial", "./testdata/Arial.ttf"),

		sahar.Width(200),
		sahar.Height(200),

		sahar.Stack(
			sahar.Alignments(sahar.Center, sahar.Middle),
			sahar.BackgroundColor("#FF0000"),

			sahar.Text(
				"Hello, World!",
				sahar.Color("#000000"),
				sahar.FontSize(8),
			),

			sahar.Image("./testdata/Sample.jpeg"),
		),
	)

	err := sahar.Reflow(b)
	assert.NoError(t, err)

	drawPdf(t, b, 1)
}
