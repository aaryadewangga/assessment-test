package entity

type AddNewProductRequest struct {
	ProductName string  `json:"productName" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Stock       int     `json:"stock" validate:"required"`
}

type AddNewProductResponse struct {
	Id          string  `json:"id"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type GetAllProductResponse struct {
	Products []AddNewProductResponse `json:"products"`
}
