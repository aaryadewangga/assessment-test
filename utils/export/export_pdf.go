package export

import (
	"aegis/assessment-test/core/entity"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func GenerateTransactionPDF(trx *entity.TransactionResponse) (string, error) {
	exportDir := "export_pdf"
	if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create export directory: %v", err)
	}

	filePath := filepath.Join(exportDir, fmt.Sprintf("transaction_%s.pdf", trx.ID))

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Transaction Details")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Transaction ID: "+trx.ID)
	pdf.Ln(8)
	pdf.Cell(40, 10, "User ID: "+trx.UserID)
	pdf.Ln(8)
	pdf.Cell(40, 10, "Total Amount: $"+formatFloat(trx.TotalAmount))
	pdf.Ln(8)
	pdf.Cell(40, 10, "Created At: "+trx.CreatedAt.String())
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(50, 8, "Product", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Quantity", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Price", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Subtotal", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	for _, detail := range *trx.Details {
		pdf.CellFormat(50, 8, detail.ProductName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, formatInt(detail.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, formatFloat(detail.Price), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, formatFloat(detail.Subtotal), "1", 1, "C", false, 0, "")
	}

	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", fmt.Errorf("failed to save PDF file: %v", err)
	}

	return filePath, nil
}

func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func formatInt(value int) string {
	return fmt.Sprintf("%d", value)
}
