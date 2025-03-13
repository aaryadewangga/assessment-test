package models

import "time"

type User struct {
	tableName struct{}  `pg:"_users"`
	ID        int       `pg:"ID,pk"`
	Name      string    `pg:"NAME"`
	Username  string    `pg:"USERNAME"`
	Password  string    `pg:"PASSWORD"`
	Role      string    `pg:"ROLE"`
	CreatedAt time.Time `pg:"CREATED_AT"`
}
