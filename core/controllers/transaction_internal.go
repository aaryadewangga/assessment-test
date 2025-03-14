package controllers

import (
	"aegis/assessment-test/core/entity"
	"aegis/assessment-test/core/repository/models"
	"context"
	"fmt"
)

func (t *TransactionController) fetchTransactionDetailsById(id string) (*entity.TransactionResponse, error) {
	transaction, err := t.transactionRepo.GetTransactionByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id=%s: %w", id, err)
	}

	details, err := t.transactionRepo.GetTransactionDetailsByTransactionID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction details id=%s: %w", id, err)
	}

	var detailsList []entity.TransactionProductResponse
	for _, v := range details {
		detailsList = append(detailsList, entity.TransactionProductResponse{
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			Quantity:    v.Quantity,
			Price:       v.Price,
			Subtotal:    v.Subtotal,
		})
	}

	resp := &entity.TransactionResponse{
		ID:          transaction.ID,
		UserID:      transaction.UserID,
		TotalAmount: transaction.TotalAmount,
		CreatedAt:   transaction.CreatedAt,
		Details:     &detailsList,
	}

	return resp, nil
}

func (t *TransactionController) parseTransactionDetails(details []models.TransactionDetailSchema) *[]entity.TransactionProductResponse {
	if len(details) == 0 {
		return nil
	}

	var detailsList []entity.TransactionProductResponse
	for _, d := range details {
		detailsList = append(detailsList, entity.TransactionProductResponse{
			ProductID:   d.ProductID,
			ProductName: d.ProductName,
			Quantity:    d.Quantity,
			Price:       d.Price,
			Subtotal:    d.Subtotal,
		})
	}
	return &detailsList
}

func (t *TransactionController) parseCreateTransactionRequest(req *entity.CreateTransactionRequest) []models.TransactionDetailSchema {
	var trxDetails []models.TransactionDetailSchema

	for _, product := range req.Products {
		trxDetails = append(trxDetails, models.TransactionDetailSchema{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		})
	}

	return trxDetails
}
