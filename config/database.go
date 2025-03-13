package config

import (
	"context"

	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
)

func NewConnPG() *pg.DB {
	cfg := Get()
	return NewWithOption(cfg)
}

func NewWithOption(cfg *Config) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.DbPgAddres,
		User:     cfg.DbPgUsername,
		Password: cfg.DbPgPassword,
		Database: cfg.DbPgName,
		OnConnect: func(ctx context.Context, cn *pg.Conn) error {
			// set the timezone
			log.Debugln("Set database connection timezone to Asia/Jakarta")
			_, err := cn.ExecContext(ctx, "set time zone ?", cfg.DBTimezone)
			return err
		},
		PoolSize: 3,
	})

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf(`failed to ping DB instance, err: %s`, err.Error())
	}

	return db
}
