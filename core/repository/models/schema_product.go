package models

type ProductSchema struct {
	//lint:ignore U1000 This field is used by ORM
	tableName   struct{} `pg:"_products"`
	ID          int      `pg:"ID,pk"`
	ProductName string   `pg:"PRODUCT_NAME"`
	Price       string   `pg:"PRICE"`
	Stock       string   `pg:"STOCK"`
}
