package config

import "time"

type RedisConfig struct {
	Host string
	Port int
	DB   int
	TTL  time.Duration
}
