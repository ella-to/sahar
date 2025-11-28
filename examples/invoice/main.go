package main

import (
	"fmt"
	"os"

	"ella.to/sahar"
)

// InvoiceItem represents a single line item on the invoice
type InvoiceItem struct {
	Name        string
	Description string
	Quantity    int
	UnitPrice   float64
}

func (i InvoiceItem) Total() float64 {
	return float64(i.Quantity) * i.UnitPrice
}

// Header creates the invoice header with company info and invoice number
func Header(invoiceNumber string) *sahar.Node {
	return sahar.Box(
		sahar.Direction(sahar.LeftToRight),
		sahar.Sizing(sahar.Grow(), sahar.Fit()),
		sahar.Padding(0, 0, 20, 0),

		// Company Info (Left side)
		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.Grow(), sahar.Fit()),
			sahar.ChildGap(6),

			sahar.Text(
				"ACME Corporation",
				sahar.FontType("Arial"),
				sahar.FontSize(24),
				sahar.FontColor("#2c3e50"),
			),
			sahar.Text(
				"123 Business Street",
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#7f8c8d"),
			),
			sahar.Text(
				"Toronto, ON M5V 1A1",
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#7f8c8d"),
			),
			sahar.Text(
				"Phone: (416) 555-1234",
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#7f8c8d"),
			),
		),

		// Invoice Info (Right side)
		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.Fit(), sahar.Fit()),
			sahar.Alignment(sahar.Right, sahar.Top),
			sahar.ChildGap(6),

			sahar.Text(
				"INVOICE",
				sahar.FontType("Arial"),
				sahar.FontSize(28),
				sahar.FontColor("#3498db"),
			),
			sahar.Text(
				fmt.Sprintf("Invoice #: %s", invoiceNumber),
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#2c3e50"),
			),
			sahar.Text(
				"Date: November 28, 2025",
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#2c3e50"),
			),
			sahar.Text(
				"Due Date: December 28, 2025",
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#2c3e50"),
			),
		),
	)
}

// BillTo creates the billing information section
func BillTo() *sahar.Node {
	return sahar.Box(
		sahar.Direction(sahar.TopToBottom),
		sahar.Sizing(sahar.Grow(), sahar.Fit()),
		sahar.Padding(20, 0, 20, 0),
		sahar.ChildGap(6),

		sahar.Text(
			"Bill To:",
			sahar.FontType("Arial"),
			sahar.FontSize(12),
			sahar.FontColor("#3498db"),
		),
		sahar.Text(
			"John Smith",
			sahar.FontType("Arial"),
			sahar.FontSize(11),
			sahar.FontColor("#2c3e50"),
		),
		sahar.Text(
			"456 Customer Avenue",
			sahar.FontType("Arial"),
			sahar.FontSize(10),
			sahar.FontColor("#7f8c8d"),
		),
		sahar.Text(
			"Vancouver, BC V6B 2K8",
			sahar.FontType("Arial"),
			sahar.FontSize(10),
			sahar.FontColor("#7f8c8d"),
		),
		sahar.Text(
			"Email: john.smith@email.com",
			sahar.FontType("Arial"),
			sahar.FontSize(10),
			sahar.FontColor("#7f8c8d"),
		),
	)
}

// TableHeader creates the header row for the items table
func TableHeader() *sahar.Node {
	return sahar.Box(
		sahar.Direction(sahar.LeftToRight),
		sahar.Sizing(sahar.Grow(), sahar.Fixed(30)),
		sahar.BackgroundColor("#3498db"),
		sahar.Padding(8, 10, 8, 10),
		sahar.Alignment(sahar.Left, sahar.Middle),

		// Item Name
		sahar.Box(
			sahar.Sizing(sahar.Fixed(120), sahar.Fit()),
			sahar.Text(
				"Item",
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#ffffff"),
			),
		),

		// Description
		sahar.Box(
			sahar.Sizing(sahar.Fixed(165), sahar.Fit()),
			sahar.Text(
				"Description",
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#ffffff"),
			),
		),

		// Quantity
		sahar.Box(
			sahar.Sizing(sahar.Fixed(50), sahar.Fit()),
			sahar.Alignment(sahar.Center, sahar.Middle),
			sahar.Text(
				"Qty",
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#ffffff"),
			),
		),

		// Unit Price
		sahar.Box(
			sahar.Sizing(sahar.Fixed(70), sahar.Fit()),
			sahar.Alignment(sahar.Right, sahar.Middle),
			sahar.Text(
				"Unit Price",
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#ffffff"),
			),
		),

		// Total
		sahar.Box(
			sahar.Sizing(sahar.Fixed(70), sahar.Fit()),
			sahar.Alignment(sahar.Right, sahar.Middle),
			sahar.Text(
				"Total",
				sahar.FontType("Arial"),
				sahar.FontSize(11),
				sahar.FontColor("#ffffff"),
			),
		),
	)
}

// TableRow creates a single row for an invoice item
func TableRow(item InvoiceItem, isAlternate bool) *sahar.Node {
	bgColor := "#ffffff"
	if isAlternate {
		bgColor = "#f8f9fa"
	}

	return sahar.Box(
		sahar.Direction(sahar.LeftToRight),
		sahar.Sizing(sahar.Grow(), sahar.Fit()),
		sahar.BackgroundColor(bgColor),
		sahar.Padding(8, 10, 8, 10),
		sahar.Alignment(sahar.Left, sahar.Middle),

		// Item Name
		sahar.Box(
			sahar.Sizing(sahar.Fixed(120), sahar.Fit()),
			sahar.Text(
				item.Name,
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#2c3e50"),
			),
		),

		// Description
		sahar.Box(
			sahar.Sizing(sahar.Fixed(165), sahar.Fit()),
			sahar.Text(
				item.Description,
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#7f8c8d"),
			),
		),

		// Quantity
		sahar.Box(
			sahar.Sizing(sahar.Fixed(50), sahar.Fit()),
			sahar.Alignment(sahar.Center, sahar.Middle),
			sahar.Text(
				fmt.Sprintf("%d", item.Quantity),
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#2c3e50"),
			),
		),

		// Unit Price
		sahar.Box(
			sahar.Sizing(sahar.Fixed(70), sahar.Fit()),
			sahar.Alignment(sahar.Right, sahar.Middle),
			sahar.Text(
				fmt.Sprintf("$%.2f", item.UnitPrice),
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#2c3e50"),
			),
		),

		// Total
		sahar.Box(
			sahar.Sizing(sahar.Fixed(70), sahar.Fit()),
			sahar.Alignment(sahar.Right, sahar.Middle),
			sahar.Text(
				fmt.Sprintf("$%.2f", item.Total()),
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#2c3e50"),
			),
		),
	)
}

// ItemsTable creates the complete items table
func ItemsTable(items []InvoiceItem) *sahar.Node {
	// Build item rows
	rows := make([]*sahar.Node, len(items))
	for i, item := range items {
		rows[i] = TableRow(item, i%2 == 1)
	}

	return sahar.Box(
		sahar.Direction(sahar.TopToBottom),
		sahar.Sizing(sahar.Grow(), sahar.Fit()),
		TableHeader(),
		sahar.Children(rows...),
	)
}

// SummaryRow creates a single summary row (subtotal, tax, total)
func SummaryRow(label string, value float64, isBold bool) *sahar.Node {
	fontSize := 10.0
	fontColor := "#2c3e50"

	if isBold {
		fontSize = 12
		fontColor = "#2c3e50"
	}

	return sahar.Box(
		sahar.Direction(sahar.LeftToRight),
		sahar.Sizing(sahar.Grow(), sahar.Fit()),
		sahar.Padding(4, 0, 4, 0),
		sahar.Alignment(sahar.Right, sahar.Middle),
		sahar.ChildGap(20),

		sahar.Text(
			label,
			sahar.FontType("Arial"),
			sahar.FontSize(fontSize),
			sahar.FontColor(fontColor),
		),

		sahar.Box(
			sahar.Sizing(sahar.Fixed(100), sahar.Fit()),
			sahar.Alignment(sahar.Right, sahar.Middle),
			sahar.Text(
				fmt.Sprintf("$%.2f", value),
				sahar.FontType("Arial"),
				sahar.FontSize(fontSize),
				sahar.FontColor(fontColor),
			),
		),
	)
}

// Summary creates the summary section with subtotal, tax, and total
func Summary(items []InvoiceItem) *sahar.Node {
	var subtotal float64
	for _, item := range items {
		subtotal += item.Total()
	}

	taxRate := 0.13 // 13% tax
	tax := subtotal * taxRate
	total := subtotal + tax

	return sahar.Box(
		sahar.Direction(sahar.TopToBottom),
		sahar.Sizing(sahar.Grow(), sahar.Fit()),
		sahar.Padding(20, 0, 0, 0),
		sahar.Alignment(sahar.Right, sahar.Top),

		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.Fixed(250), sahar.Fit()),
			sahar.Padding(10, 10, 10, 10),
			sahar.Border(1),
			sahar.BorderColor("#ecf0f1"),
			sahar.ChildGap(3),

			SummaryRow("Subtotal:", subtotal, false),
			SummaryRow("Tax (13%):", tax, false),

			// Divider line
			sahar.Box(
				sahar.Sizing(sahar.Grow(), sahar.Fixed(1)),
				sahar.BackgroundColor("#bdc3c7"),
			),

			SummaryRow("Total:", total, true),
		),
	)
}

// Footer creates the footer with thank you message and payment info
func Footer() *sahar.Node {
	return sahar.Box(
		sahar.Direction(sahar.TopToBottom),
		sahar.Sizing(sahar.Grow(), sahar.Fit()),
		sahar.Padding(30, 0, 0, 0),
		sahar.ChildGap(10),
		sahar.Alignment(sahar.Center, sahar.Middle),

		sahar.Text(
			"Thank you for your business!",
			sahar.FontType("Arial"),
			sahar.FontSize(14),
			sahar.FontColor("#3498db"),
		),

		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.Grow(), sahar.Fit()),
			sahar.ChildGap(6),
			sahar.Alignment(sahar.Center, sahar.Middle),

			sahar.Text(
				"Payment Terms: Net 30 days",
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#7f8c8d"),
			),
			sahar.Text(
				"Please make checks payable to ACME Corporation",
				sahar.FontType("Arial"),
				sahar.FontSize(10),
				sahar.FontColor("#7f8c8d"),
			),
		),
	)
}

func main() {
	// Load fonts
	err := sahar.LoadFonts("Arial", "./Arial.ttf")
	if err != nil {
		panic(err)
	}

	// Invoice number
	invoiceNumber := "INV-2025-001234"

	// Sample invoice items
	items := []InvoiceItem{
		{
			Name:        "Web Development",
			Description: "Custom website design and development",
			Quantity:    1,
			UnitPrice:   2500.00,
		},
		{
			Name:        "Hosting (Annual)",
			Description: "Cloud hosting service for 12 months",
			Quantity:    1,
			UnitPrice:   299.00,
		},
		{
			Name:        "SSL Certificate",
			Description: "Premium SSL certificate for secure connections",
			Quantity:    2,
			UnitPrice:   79.99,
		},
		{
			Name:        "Domain Name",
			Description: "Domain registration for 1 year (.com)",
			Quantity:    1,
			UnitPrice:   14.99,
		},
		{
			Name:        "Support Hours",
			Description: "Technical support and maintenance hours",
			Quantity:    10,
			UnitPrice:   75.00,
		},
	}

	// Build the invoice page
	page := sahar.Layout(
		sahar.Box(
			sahar.Direction(sahar.TopToBottom),
			sahar.Sizing(sahar.A4()...),
			sahar.Padding(50, 50, 50, 50),
			sahar.BackgroundColor("#ffffff"),

			Header(invoiceNumber),
			BillTo(),
			ItemsTable(items),
			Summary(items),

			// Spacer to push footer to bottom
			sahar.Box(
				sahar.Sizing(sahar.Grow(), sahar.Grow()),
			),

			Footer(),
		),
	)

	// Write the layout to a PDF file
	pdfFile, err := os.Create("./invoice.pdf")
	if err != nil {
		panic(err)
	}
	defer pdfFile.Close()

	err = sahar.RenderToPDF(pdfFile, page)
	if err != nil {
		panic(err)
	}

	fmt.Println("Invoice PDF generated successfully: invoice.pdf")
}
