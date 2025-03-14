package models

import "time"

type TransactionSchema struct {
	//lint:ignore U1000 This field is used by ORM
	tableName   struct{}  `pg:"_transactions,alias:t"`
	ID          string    `pg:"ID,pk,type:uuid"`
	UserID      string    `pg:"USER_ID,type:uuid"`
	TotalAmount float64   `pg:"TOTAL_AMOUNT"`
	CreatedAt   time.Time `pg:"CREATED_AT"`

	User    *UserSchema               `pg:"rel:has-one,fk:USER_ID"`
	Details []TransactionDetailSchema `pg:"rel:has-many,fk:TRANSACTION_ID"`
}

type TransactionDetailSchema struct {
	//lint:ignore U1000 This field is used by ORM
	tableName     struct{} `pg:"_transaction_details,alias:d"`
	ID            string   `pg:"ID,pk,type:uuid"`
	TransactionID string   `pg:"TRANSACTION_ID,type:uuid"`
	ProductID     string   `pg:"PRODUCT_ID,type:uuid"`
	ProductName   string   `pg:"PRODUCT_NAME"`
	Price         float64  `pg:"PRICE"`
	Quantity      int      `pg:"QUANTITY"`
	Subtotal      float64  `pg:"SUBTOTAL"`

	Product *ProductSchema `pg:"rel:has-one,fk:PRODUCT_ID"`
}
