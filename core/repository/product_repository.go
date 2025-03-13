package repository

import (
	"aegis/assessment-test/core/repository/models"
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
)

type ProductRepository interface {
	InsertNewProduct(ctx context.Context, product *models.ProductSchema) (*models.ProductSchema, error)
	GetAllProducts(ctx context.Context) (*[]models.ProductSchema, error)
	GetProductByID(ctx context.Context, id int) (*models.ProductSchema, error)
	UpdateProduct(ctx context.Context, product *models.ProductSchema) (*models.ProductSchema, error)
	DeleteProduct(ctx context.Context, id int) error
}

type productRepository struct {
	db *pg.DB
}

func NewProductRepository(db *pg.DB) ProductRepository {
	return &productRepository{db: db}
}

func (p *productRepository) InsertNewProduct(ctx context.Context, product *models.ProductSchema) (*models.ProductSchema, error) {
	res, err := p.getProductByProductName(ctx, product.ProductName)
	if err != pg.ErrNoRows || res != nil {
		return nil, errors.New("new product has same name")
	}

	_, err = p.db.Model(product).Context(ctx).Insert()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productRepository) GetAllProducts(ctx context.Context) (*[]models.ProductSchema, error) {
	var products []models.ProductSchema
	err := p.db.Model(&products).
		Context(ctx).
		Order("PRODUCT_NAME ASC").
		Select()
	return &products, err
}

func (p *productRepository) GetProductByID(ctx context.Context, id int) (*models.ProductSchema, error) {
	product := new(models.ProductSchema)
	err := p.db.Model(product).
		Context(ctx).
		Where("? = ?", pg.Ident("ID"), id).
		Limit(1).
		Select()
	if err != nil {
		return nil, err
	}

	return product, err
}

func (p *productRepository) getProductByProductName(ctx context.Context, productName string) (*models.ProductSchema, error) {
	product := new(models.ProductSchema)
	err := p.db.Model(product).
		Context(ctx).
		Where("? = ?", pg.Ident("PRODUCT_NAME"), productName).
		Limit(1).
		Select()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productRepository) UpdateProduct(ctx context.Context, product *models.ProductSchema) (*models.ProductSchema, error) {
	_, err := p.db.Model(product).
		Context(ctx).
		Where("? = ?", pg.Ident("ID"), product.ID).
		Returning("*").
		Update()

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productRepository) DeleteProduct(ctx context.Context, id int) error {
	_, err := p.db.Model((*models.ProductSchema)(nil)).Context(ctx).
		Where("? = ?", pg.Ident("ID"), id).
		Delete()
	return err
}
