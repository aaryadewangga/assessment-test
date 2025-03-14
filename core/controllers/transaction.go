package controllers

import (
	"aegis/assessment-test/core/constant"
	"aegis/assessment-test/core/entity"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/utils/export"
	"aegis/assessment-test/utils/middleware"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
}

func NewTransactionController(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
) *TransactionController {
	return &TransactionController{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (t *TransactionController) CreateTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userData, err := middleware.GetTokenClaims(c)
		if err != nil {
			logrus.Errorf("err claims data token=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "failed claims token data", nil))
		}

		req := entity.CreateTransactionRequest{}
		c.Bind(&req)
		err = c.Validate(&req)
		if err != nil {
			logrus.Errorf("err validate request=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "there is some problem from input", err))
		}

		user, err := t.userRepo.GetUserById(context.Background(), userData.UserId)
		if err != nil {
			logrus.Errorf("err get user data=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "user not registered", err))
		}

		trx, err := t.transactionRepo.CreateTransaction(context.Background(), user.ID, t.parseCreateTransactionRequest(&req))
		if err != nil {
			logrus.Errorf("err create trx data=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to create transactions", err))
		}

		resp := entity.TransactionResponse{
			ID:          trx.ID,
			UserID:      trx.UserID,
			TotalAmount: trx.TotalAmount,
			CreatedAt:   trx.CreatedAt,
			Details:     t.parseTransactionDetails(trx.Details),
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func (t *TransactionController) GetAllTransactions() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactions, err := t.transactionRepo.GetAllTransactions()
		if err != nil {
			logrus.Errorf("err get all transactions: %s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to get transactions", err))
		}

		var resp []entity.TransactionResponse
		for _, trx := range transactions {
			resp = append(resp, entity.TransactionResponse{
				ID:          trx.ID,
				UserID:      trx.UserID,
				TotalAmount: trx.TotalAmount,
				CreatedAt:   trx.CreatedAt,
			})
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success get all transactions", resp))
	}
}

func (t *TransactionController) GetTransactionDetailsById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		resp, err := t.fetchTransactionDetailsById(id)
		if err != nil {
			logrus.Errorf("error fetching transaction details: %s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to get transactions", err))
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success get details transactions", resp))
	}
}

func (t *TransactionController) GetTransactionPDF() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		data, err := t.fetchTransactionDetailsById(id)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to get details transactions", err))
		}

		// Set response sebagai file PDF
		c.Response().Header().Set("Content-Type", "application/pdf")
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=transaction_%s.pdf", id))

		// Generate PDF
		err = export.GenerateTransactionPDF(data, c.Response().Writer)
		if err != nil {
			logrus.Errorf("failed to generate transaction PDF: %s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to generate pdf", err))
		}

		return nil
	}
}

func (t *TransactionController) GetTransactionExcel() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, constant.InternalServerError(constant.CodeErrInternalServer, "transaction ID is required", nil))
		}

		data, err := t.fetchTransactionDetailsById(id)
		if err != nil {
			logrus.Errorf("failed to fetch transaction: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, constant.InternalServerError(constant.CodeErrInternalServer, "failed to get transaction details", err))
		}
		if data == nil {
			logrus.Warnf("transaction not found: ID=%s", id)
			return c.JSON(http.StatusNotFound, constant.InternalServerError(constant.CodeErrInternalServer, "transaction not found", nil))
		}

		filepath, err := export.GenerateTransactionExcel(data)
		if err != nil {
			logrus.Errorf("failed to generate transaction Excel: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, constant.InternalServerError(constant.CodeErrInternalServer, "failed to generate excel", err))
		}

		file, err := os.Open(filepath)
		if err != nil {
			logrus.Errorf("failed to open Excel file: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, constant.InternalServerError(constant.CodeErrInternalServer, "failed to open file", err))
		}
		defer file.Close()

		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath))
		c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		return c.Stream(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file)
	}
}
