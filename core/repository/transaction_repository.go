package repository

import (
	"aegis/assessment-test/core/repository/models"
	"context"
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, userID string, details []models.TransactionDetailSchema) (*models.TransactionSchema, error)
	GetAllTransactions() ([]models.TransactionSchema, error)
	GetTransactionByID(ctx context.Context, id string) (*models.TransactionSchema, error)
	GetTransactionDetailsByTransactionID(ctx context.Context, transactionID string) ([]models.TransactionDetailSchema, error)
}

type transactionRepository struct {
	db *pg.DB
}

func NewTransactionRepository(db *pg.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, userID string, details []models.TransactionDetailSchema) (*models.TransactionSchema, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	var totalAmount float64
	for i := range details {
		var product models.ProductSchema
		err := tx.Model(&product).Where("? = ?", pg.Ident("ID"), details[i].ProductID).Select()
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < details[i].Quantity {
			return nil, errors.New("insufficient stock for product: " + product.ProductName)
		}

		details[i].ProductName = product.ProductName
		details[i].Price = product.Price
		details[i].Subtotal = product.Price * float64(details[i].Quantity)
		totalAmount += product.Price * float64(details[i].Quantity)
	}

	trx := &models.TransactionSchema{
		UserID:      userID,
		TotalAmount: totalAmount,
	}
	_, err = tx.Model(trx).Insert()
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = trx.ID
	}

	for i := range details {
		_, err := tx.Model((*models.ProductSchema)(nil)).
			Set("? = ? - ?", pg.Ident("STOCK"), pg.Ident("STOCK"), details[i].Quantity).
			Where("? = ?", pg.Ident("ID"), details[i].ProductID).
			Update()
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Model(&details).Insert()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return trx, nil
}

func (r *transactionRepository) GetAllTransactions() ([]models.TransactionSchema, error) {
	var transactions []models.TransactionSchema
	err := r.db.Model(&transactions).
		Relation("User").
		Relation("Details.Product").
		Select()
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, id string) (*models.TransactionSchema, error) {
	transaction := new(models.TransactionSchema)

	err := r.db.Model(transaction).
		Context(ctx).
		Where("? = ?", pg.Ident("ID"), id).
		Limit(1).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) GetTransactionDetailsByTransactionID(ctx context.Context, transactionID string) ([]models.TransactionDetailSchema, error) {
	var details []models.TransactionDetailSchema

	err := r.db.Model(&details).
		Context(ctx).
		Where("? = ?", pg.Ident("TRANSACTION_ID"), transactionID).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction details: %w", err)
	}

	return details, nil
}
