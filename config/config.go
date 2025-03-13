package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// app base conf
	AppName string `envconfig:"APP_NAME" default:""`
	AppPort string `envconfig:"APP_PORT" default:""`

	// db
	DbPgAddres   string `envconfig:"DB_PG_ADDR" default:""`
	DbPgUsername string `envconfig:"DB_PG_USER" default:""`
	DbPgPassword string `envconfig:"DB_PG_PASS" default:""`
	DbPgName     string `envconfig:"DB_PG_DBNAME" default:""`
	DBTimezone   string `envconfig:"DB_PG_TZ" default:"Asia/Jakarta"`

	// JWT secret
	JWTPublicKey string `envconfig:"JWT_PUBLIC_KEY" default:""`
}

func Get() *Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return &cfg
}
