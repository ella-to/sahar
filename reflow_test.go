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
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,
		}

		expexted := &sahar.Node{
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Size on a Stack with 2 elements without any width and height defined", func(t *testing.T) {
		b := &sahar.Node{
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
			},
		}

		expexted := &sahar.Node{
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 50,
				},
				{
					Type:   sahar.Stack,
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
			Type:   sahar.Group,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
			},
		}

		expexted := &sahar.Node{
			Type:   sahar.Group,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Type:   sahar.Stack,
					Width:  50,
					Height: 100,
				},
				{
					Type:   sahar.Stack,
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
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Type:  sahar.Stack,
					Width: 50,
				},
				{
					Type:  sahar.Stack,
					Width: 60,
				},
				{
					Type: sahar.Stack,
				},
			},
		}

		expexted := &sahar.Node{
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,

			Children: []*sahar.Node{
				{
					Type:   sahar.Stack,
					Width:  50,
					Height: 33.333333333333336,
				},
				{
					Type:   sahar.Stack,
					Width:  60,
					Height: 33.333333333333336,
				},
				{
					Type:   sahar.Stack,
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
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,
		}

		expexted := &sahar.Node{
			Type:   sahar.Stack,
			Width:  100,
			Height: 100,
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		err = sahar.Alignment(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Group", func(t *testing.T) {
		b := &sahar.Node{
			Type:   sahar.Group,
			Width:  100,
			Height: 100,
		}

		expexted := &sahar.Node{
			Type:   sahar.Group,
			Width:  100,
			Height: 100,
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		err = sahar.Alignment(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children and TopLeft", func(t *testing.T) {
		b := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Left,
			VerticalAlignment:   sahar.Top,

			Children: []*sahar.Node{
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
			},
		}

		expexted := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Left,
			VerticalAlignment:   sahar.Top,
			Children: []*sahar.Node{
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      0,
				},
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      33.333333333333336,
				},
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		err = sahar.Alignment(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children and TopCenter", func(t *testing.T) {
		b := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Center,
			VerticalAlignment:   sahar.Top,

			Children: []*sahar.Node{
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
			},
		}

		expexted := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Center,
			VerticalAlignment:   sahar.Top,
			Children: []*sahar.Node{
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      0,
				},
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      33.333333333333336,
				},
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		err = sahar.Alignment(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children and TopRight", func(t *testing.T) {
		b := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Right,
			VerticalAlignment:   sahar.Top,

			Children: []*sahar.Node{
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
				{
					Type: sahar.Stack,
				},
			},
		}

		expexted := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Right,
			VerticalAlignment:   sahar.Top,
			Children: []*sahar.Node{
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      0,
				},
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      33.333333333333336,
				},
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		err = sahar.Alignment(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})

	t.Run("Test Alignment on a Stack with 3 children width differnt width and TopCenter", func(t *testing.T) {
		b := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Center,
			VerticalAlignment:   sahar.Top,

			Children: []*sahar.Node{
				{
					Type:  sahar.Stack,
					Width: 50,
				},
				{
					Type:  sahar.Stack,
					Width: 60,
				},
				{
					Type: sahar.Stack,
				},
			},
		}

		expexted := &sahar.Node{
			Type:                sahar.Stack,
			Width:               100,
			Height:              100,
			HorizontalAlignment: sahar.Center,
			VerticalAlignment:   sahar.Top,
			Children: []*sahar.Node{
				{
					Type:   sahar.Stack,
					Width:  50,
					Height: 33.333333333333336,
					X:      25,
					Y:      0,
				},
				{
					Type:   sahar.Stack,
					Width:  60,
					Height: 33.333333333333336,
					X:      20,
					Y:      33.333333333333336,
				},
				{
					Type:   sahar.Stack,
					Width:  100,
					Height: 33.333333333333336,
					X:      0,
					Y:      66.66666666666667,
				},
			},
		}

		err := sahar.Resize(b)
		assert.NoError(t, err)
		err = sahar.Alignment(b)
		assert.NoError(t, err)
		assert.Equal(t, expexted, b)
	})
}
