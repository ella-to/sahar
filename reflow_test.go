package sahar_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"ella.to/sahar"
	"ella.to/sahar/pdf"
)

func drawPdf(t *testing.T, node *sahar.Node, i int) error {
	file, err := os.Create(fmt.Sprintf("./test_%d.pdf", i))
	if err != nil {
		return err
	}
	t.Cleanup(func() {
		file.Close()
	})
	err = pdf.Write(file, node)
	return err
}

func TestResize(t *testing.T) {

	t.Run("Test Size on a Stack", func(t *testing.T) {
		b := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,
		}

		expexted := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Size on a Stack with 2 elements without any width and height defined", func(t *testing.T) {
		b := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
			},
		}

		expexted := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 50,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 50,
				},
			},
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Size on a Group with 2 elements with width and height defined", func(t *testing.T) {
		b := &sahar.Node{
			Order:  sahar.GroupOrder,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
			},
		}

		expexted := &sahar.Node{
			Order:  sahar.GroupOrder,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Order:  sahar.StackOrder,
					Width:  50,
					Height: 100,
				},
				{
					Order:  sahar.StackOrder,
					Width:  50,
					Height: 100,
				},
			},
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Size on a Stack with 3 elements with 2 define widths defined", func(t *testing.T) {
		b := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Order: sahar.StackOrder,
					Width: 50,
				},
				{
					Order: sahar.StackOrder,
					Width: 60,
				},
				{
					Order: sahar.StackOrder,
				},
			},
		}

		expexted := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Order:  sahar.StackOrder,
					Width:  50,
					Height: 33.333333333333336,
				},
				{
					Order:  sahar.StackOrder,
					Width:  60,
					Height: 33.333333333333336,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
				},
			},
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})
}

func TestAlignment(t *testing.T) {
	t.Run("Test Alignment on a Stack", func(t *testing.T) {
		b := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,
		}

		expexted := &sahar.Node{
			Order:  sahar.StackOrder,
			Width:  100,
			Height: 100,
		}

		err := sahar.Reflow(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Group", func(t *testing.T) {
		b := &sahar.Node{
			Order:  sahar.GroupOrder,
			Width:  100,
			Height: 100,
		}

		expexted := &sahar.Node{
			Order:  sahar.GroupOrder,
			Width:  100,
			Height: 100,
		}

		err := sahar.Reflow(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children and TopLeft", func(t *testing.T) {
		b := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Left,
			Vertical:   sahar.Top,

			Children: []*sahar.Node{
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
			},
		}

		expexted := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Left,
			Vertical:   sahar.Top,
			Children: []*sahar.Node{
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      0,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      33.333333333333336,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Reflow(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children and TopCenter", func(t *testing.T) {
		b := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Center,
			Vertical:   sahar.Top,

			Children: []*sahar.Node{
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
			},
		}

		expexted := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Center,
			Vertical:   sahar.Top,
			Children: []*sahar.Node{
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      0,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      33.333333333333336,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Reflow(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children and TopRight", func(t *testing.T) {
		b := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Right,
			Vertical:   sahar.Top,

			Children: []*sahar.Node{
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
				{
					Order: sahar.StackOrder,
				},
			},
		}

		expexted := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Right,
			Vertical:   sahar.Top,
			Children: []*sahar.Node{
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      0,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      33.333333333333336,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Reflow(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children width differnt width and TopCenter", func(t *testing.T) {
		b := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Center,
			Vertical:   sahar.Top,

			Children: []*sahar.Node{
				{
					Order: sahar.StackOrder,
					Width: 50,
				},
				{
					Order: sahar.StackOrder,
					Width: 60,
				},
				{
					Order: sahar.StackOrder,
				},
			},
		}

		expexted := &sahar.Node{
			Order:      sahar.StackOrder,
			Width:      100,
			Height:     100,
			Horizontal: sahar.Center,
			Vertical:   sahar.Top,
			Children: []*sahar.Node{
				{
					Order:  sahar.StackOrder,
					Width:  50,
					Height: 33.333333333333336,
					X:      25,
					Y:      0,
				},
				{
					Order:  sahar.StackOrder,
					Width:  60,
					Height: 33.333333333333336,
					X:      20,
					Y:      33.333333333333336,
				},
				{
					Order:  sahar.StackOrder,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Reflow(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})
}
