package entity

type AddNewProductRequest struct {
	ProductName string `json:"productName" validate:"required"`
	Price       string `json:"price" validate:"required"`
	Stock       string `json:"stock" validate:"required"`
}

type AddNewProductResponse struct {
	Id          int    `json:"id"`
	ProductName string `json:"productName"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
}

type GetAllProductResponse struct {
	Products []AddNewProductResponse `json:"products"`
}
