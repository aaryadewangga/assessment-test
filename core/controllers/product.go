package controllers

import (
	"aegis/assessment-test/core/constant"
	"aegis/assessment-test/core/entity"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/core/repository/models"
	"aegis/assessment-test/utils/middleware"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (
	RoleAdmin string = "admin"
)

type ProductController struct {
	productRepo repository.ProductRepository
}

func NewProductController(
	productRepo repository.ProductRepository,
) *ProductController {
	return &ProductController{productRepo: productRepo}
}

func (p *ProductController) AddNewProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		if !middleware.IsAdmin(c) {
			logrus.Errorf("access denied")
			return c.JSON(
				http.StatusForbidden,
				constant.UnauthorizeError(constant.CodeErrForbidden, "access denied", nil))
		}
		req := entity.AddNewProductRequest{}
		c.Bind(&req)
		err := c.Validate(&req)
		if err != nil {
			logrus.Errorf("err validate request=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "there is some problem from input", err))
		}

		product, err := p.productRepo.InsertNewProduct(context.Background(), &models.ProductSchema{
			ProductName: req.ProductName,
			Price:       req.Price,
			Stock:       req.Stock,
		})
		if err != nil {
			logrus.Errorf("err insert new product=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to insert new product", err))
		}

		resp := entity.AddNewProductResponse{
			Id:          product.ID,
			ProductName: product.ProductName,
			Price:       product.Price,
			Stock:       product.Stock,
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success insert new product", resp))
	}
}

func (p *ProductController) GetProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		if id == "" {
			return p.getAllProduct()(c)
		}
		return p.getProductById(id)(c)
	}
}

func (p *ProductController) getAllProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := p.productRepo.GetAllProducts(context.Background())
		if err != nil {
			logrus.Errorf("err get all products=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to get all products", err))
		}

		if len(*products) == 0 {
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "product not found", err))
		}

		tmpRes := entity.AddNewProductResponse{}
		resp := entity.GetAllProductResponse{}
		for _, val := range *products {
			tmpRes.Id = val.ID
			tmpRes.ProductName = val.ProductName
			tmpRes.Price = val.Price
			tmpRes.Stock = val.Stock

			resp.Products = append(resp.Products, tmpRes)
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success get all products", resp))
	}
}

func (p *ProductController) getProductById(id string) echo.HandlerFunc {
	return func(c echo.Context) error {
		product, err := p.productRepo.GetProductByID(context.Background(), id)
		if err != nil {
			logrus.Errorf("err get product by id=%s err=%s", id, err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, fmt.Sprintf("failed to get product by id=%s", id), err))
		}

		resp := entity.AddNewProductResponse{
			Id:          product.ID,
			ProductName: product.ProductName,
			Price:       product.Price,
			Stock:       product.Stock,
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success get product by id", resp))
	}
}

func (p *ProductController) UpdateProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		if !middleware.IsAdmin(c) {
			logrus.Errorf("access denied")
			return c.JSON(
				http.StatusForbidden,
				constant.UnauthorizeError(constant.CodeErrForbidden, "access denied", nil))
		}
		id := c.QueryParam("id")
		product := models.ProductSchema{}
		err := c.Bind(&product)
		if err != nil {
			logrus.Errorf("err validate request=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "there is some problem from input", err))
		}

		product.ID = id
		updatedProduct, err := p.productRepo.UpdateProduct(context.Background(), &product)
		if err != nil {
			logrus.Errorf("err update product by id=%s err=%s", id, err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, fmt.Sprintf("failed to update product by id=%s", id), err))
		}

		resp := entity.AddNewProductResponse{
			Id:          updatedProduct.ID,
			ProductName: updatedProduct.ProductName,
			Price:       updatedProduct.Price,
			Stock:       updatedProduct.Stock,
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success update product by id", resp))
	}
}

func (p *ProductController) DeleteProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		if !middleware.IsAdmin(c) {
			logrus.Errorf("access denied")
			return c.JSON(
				http.StatusForbidden,
				constant.UnauthorizeError(constant.CodeErrForbidden, "access denied", nil))
		}
		id := c.QueryParam("id")
		err := p.productRepo.DeleteProduct(context.Background(), id)
		if err != nil {
			logrus.Errorf("err delete product by id=%s err=%s", id, err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, fmt.Sprintf("failed to delete product by id=%s", id), err))
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, fmt.Sprintf("success delete product by id=%s", id), nil))
	}
}
