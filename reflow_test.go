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

func TestUpdateChildrenWidthHeight(t *testing.T) {
	testCases := []struct {
		node         *sahar.Node
		expectedNode *sahar.Node
	}{
		{
			node: &sahar.Node{
				Type:   sahar.Stack,
				Margin: [4]float64{5, 5, 5, 5},
				Width:  100,
				Height: 100,
			},
			expectedNode: &sahar.Node{
				Width:  90,
				Height: 90,
				Margin: [4]float64{5, 5, 5, 5},
				Type:   sahar.Stack,
			},
		},
		{
			node: &sahar.Node{
				Type:   sahar.Stack,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
				Children: []*sahar.Node{
					{
						Type: sahar.Stack,
					},
				},
			},
			expectedNode: &sahar.Node{
				Type:   sahar.Stack,
				Width:  90,
				Height: 90,
				Margin: [4]float64{5, 5, 5, 5},
				Children: []*sahar.Node{
					{
						Width:  90,
						Height: 90,
						Type:   sahar.Stack,
					},
				},
			},
		},
		{
			node: &sahar.Node{
				Type:   sahar.Stack,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
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
			},
			expectedNode: &sahar.Node{
				Type:   sahar.Stack,
				Width:  90,
				Height: 90,
				Margin: [4]float64{5, 5, 5, 5},
				Children: []*sahar.Node{
					{
						Width:  90,
						Height: 30,
						Type:   sahar.Stack,
					},
					{
						Width:  90,
						Height: 30,
						Type:   sahar.Stack,
					},
					{
						Width:  90,
						Height: 30,
						Type:   sahar.Stack,
					},
				},
			},
		},
		{
			node: &sahar.Node{
				Type:   sahar.Group,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
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
			},
			expectedNode: &sahar.Node{
				Type:   sahar.Group,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
				Children: []*sahar.Node{
					{
						Width:  30,
						Height: 90,
						Type:   sahar.Stack,
					},
					{
						Width:  30,
						Height: 90,
						Type:   sahar.Stack,
					},
					{
						Width:  30,
						Height: 90,
						Type:   sahar.Stack,
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		sahar.UpdateChildrenWidthHeight(tc.node)
		assert.Equal(t, tc.expectedNode, tc.node)

		drawPdf(t, tc.node, i)
	}
}

func TestUpdateChildrenXY(t *testing.T) {
	testCases := []struct {
		node         *sahar.Node
		expectedNode *sahar.Node
	}{
		{
			node: &sahar.Node{
				Type:   sahar.Stack,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
			},
			expectedNode: &sahar.Node{
				Type:   sahar.Stack,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
				X:      5,
				Y:      5,
			},
		},
		{
			node: &sahar.Node{
				Type:   sahar.Stack,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
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
			},
			expectedNode: &sahar.Node{
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
			},
		},
		{
			node: &sahar.Node{
				Type:   sahar.Group,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
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
			},
			expectedNode: &sahar.Node{
				Type:   sahar.Group,
				Width:  100,
				Height: 100,
				Margin: [4]float64{5, 5, 5, 5},
				X:      5,
				Y:      5,
				Children: []*sahar.Node{
					{
						Width:  30,
						Height: 90,
						Type:   sahar.Stack,
						X:      5,
						Y:      5,
					},
					{
						Width:  30,
						Height: 90,
						Type:   sahar.Stack,
						X:      35,
						Y:      5,
					},
					{
						Width:  30,
						Height: 90,
						Type:   sahar.Stack,
						X:      65,
						Y:      5,
					},
				},
			},
		},
		{
			node: &sahar.Node{
				Type:                sahar.Stack,
				Width:               100,
				Height:              100,
				Margin:              [4]float64{5, 5, 5, 5},
				HorizontalAlignment: sahar.Center,
				Children: []*sahar.Node{
					{
						Type:  sahar.Stack,
						Width: 10,
					},
					{
						Type:  sahar.Stack,
						Width: 10,
					},
					{
						Type:  sahar.Stack,
						Width: 10,
					},
				},
			},
			expectedNode: &sahar.Node{
				Type:                sahar.Stack,
				Width:               100,
				Height:              100,
				Margin:              [4]float64{5, 5, 5, 5},
				X:                   5,
				Y:                   5,
				HorizontalAlignment: sahar.Center,
				Children: []*sahar.Node{
					{
						Width:  10,
						Height: 30,
						Type:   sahar.Stack,
						X:      45,
						Y:      5,
					},
					{
						Width:  10,
						Height: 30,
						Type:   sahar.Stack,
						X:      45,
						Y:      35,
					},
					{
						Width:  10,
						Height: 30,
						Type:   sahar.Stack,
						X:      45,
						Y:      65,
					},
				},
			},
		},
		{
			node: &sahar.Node{
				Type:              sahar.Group,
				Width:             100,
				Height:            100,
				Margin:            [4]float64{5, 5, 5, 5},
				VerticalAlignment: sahar.Middle,
				Children: []*sahar.Node{
					{
						Type:   sahar.Stack,
						Height: 10,
					},
					{
						Type:   sahar.Stack,
						Height: 10,
					},
					{
						Type:   sahar.Stack,
						Height: 10,
					},
				},
			},
			expectedNode: &sahar.Node{
				Type:              sahar.Group,
				Width:             100,
				Height:            100,
				Margin:            [4]float64{5, 5, 5, 5},
				X:                 5,
				Y:                 5,
				VerticalAlignment: sahar.Middle,
				Children: []*sahar.Node{
					{
						Width:  30,
						Height: 10,
						Type:   sahar.Stack,
						X:      5,
						Y:      45,
					},
					{
						Width:  30,
						Height: 10,
						Type:   sahar.Stack,
						X:      35,
						Y:      45,
					},
					{
						Width:  30,
						Height: 10,
						Type:   sahar.Stack,
						X:      65,
						Y:      45,
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		sahar.UpdateChildrenWidthHeight(tc.node)
		sahar.UpdateRootXY(tc.node)
		sahar.UpdateChildrenXY(tc.node, true)
		assert.Equal(t, tc.expectedNode, tc.node)
		drawPdf(t, tc.node, i)
	}
}
