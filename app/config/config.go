package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

// Config stores all configuration of the application.
type Config struct {
	Port     int    `env:"SERVER_PORT,default=5000"`
	LogLevel string `env:"LOG_LEVEL,default=info"`
	DB       struct {
		Host     string `env:"DB_HOST,required"`
		Port     uint   `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		Database string `env:"DB_DATABASE,required"`
		Driver   string `env:"DB_DRIVER,default=postgres"`
	}
}

// LoadConfig Loads the application config
func LoadConfig() (Config, error) {
	var c Config

	ctx := context.Background()

	err := envconfig.Process(ctx, &c)
	return c, err
}
