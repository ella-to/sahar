package sahar_test

import (
	"testing"

	"ella.to/sahar"
)

func TestBasicBlock(t *testing.T) {
	b := sahar.Block(
		sahar.Stack,
		sahar.Margin(5, 5, 5, 5),
		sahar.A4(),

		sahar.Block(
			sahar.Stack,
			sahar.Height(50),
		),
		sahar.Block(
			sahar.Stack,

			sahar.Block(
				sahar.Stack,
			),
			sahar.Block(
				sahar.Stack,
			),
			sahar.Block(
				sahar.Stack,
			),
		),
		sahar.Block(
			sahar.Stack,
			sahar.Height(50),
		),
	)

	sahar.Reflow(b)

	drawPdf(t, b, 1)
}
