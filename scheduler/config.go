package scheduler

import "time"

const (
	interval = 10 * time.Second
)

type Config struct {
	JobRunInterval time.Duration
}

func NewConfig() *Config {
	return &Config{JobRunInterval: interval}
}
