package models

import "time"

type UserSchema struct {
	//lint:ignore U1000 This field is used by ORM
	tableName struct{}  `pg:"_users"`
	ID        string    `pg:"ID,pk,type:uuid"`
	Name      string    `pg:"NAME"`
	Username  string    `pg:"USERNAME"`
	Password  string    `pg:"PASSWORD"`
	Role      string    `pg:"ROLE"`
	CreatedAt time.Time `pg:"CREATED_AT"`
}
