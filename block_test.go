package sahar_test

import (
	"testing"

	"ella.to/sahar"
)

func TestBasicBlock(t *testing.T) {

	b := sahar.Block(
		sahar.Stack,

		sahar.Padding(5, 5, 5, 5),

		sahar.Width(100),
		sahar.Height(100),

		sahar.Alignments(sahar.Center, sahar.Middle),

		sahar.Block(
			sahar.Stack,
		),
	)

	sahar.Reflow(b)

	drawPdf(t, b, 1)
}
