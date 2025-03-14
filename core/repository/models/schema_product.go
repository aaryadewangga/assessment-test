package models

type ProductSchema struct {
	//lint:ignore U1000 This field is used by ORM
	tableName   struct{} `pg:"_products,alias:p"`
	ID          string   `pg:"ID,pk,type:uuid"`
	ProductName string   `pg:"PRODUCT_NAME"`
	Price       float64  `pg:"PRICE"`
	Stock       int      `pg:"STOCK"`
}
