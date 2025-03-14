package export

import (
	"aegis/assessment-test/core/entity"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

func GenerateTransactionExcel(transaction *entity.TransactionResponse) (string, error) {
	outputDir := "export_excel"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}

	filename := fmt.Sprintf("transaction_%s.xlsx", transaction.ID)
	filepath := filepath.Join(outputDir, filename)

	f := excelize.NewFile()
	sheetName := "Transaction"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"Transaction ID", "UserId", "Total Amount", "Date"}
	for i, h := range headers {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s1", col)
		f.SetCellValue(sheetName, cell, h)
	}

	f.SetCellValue(sheetName, "A2", transaction.ID)
	f.SetCellValue(sheetName, "B2", transaction.UserID)
	f.SetCellValue(sheetName, "C2", transaction.TotalAmount)
	f.SetCellValue(sheetName, "D2", transaction.CreatedAt)

	productHeaders := []string{"Product ID", "Product Name", "Quantity", "Price", "Subtotal"}
	for i, h := range productHeaders {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s4", col) // Baris 4 untuk header produk
		f.SetCellValue(sheetName, cell, h)
	}

	// Isi data produk
	startRow := 5
	for i, product := range *transaction.Details {
		row := startRow + i
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.ProductID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.ProductName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), product.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), product.Price)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), product.Subtotal)
	}

	if err := f.SaveAs(filepath); err != nil {
		return "", fmt.Errorf("failed to save Excel file: %w", err)
	}

	return filepath, nil
}
