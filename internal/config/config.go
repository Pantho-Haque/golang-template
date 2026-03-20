package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres DatabaseConfig
	Redis    RedisConfig
}

func (c *Config) PrintConfig() {
	fmt.Println("App: ", c.App)
	fmt.Println("Postgres: ", c.Postgres)
	fmt.Println("Redis: ", c.Redis)
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	config := &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
		Postgres: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", ""),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			DB:       getEnv("POSTGRES_DB", ""),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnvInt("REDIS_PORT", 6379),
			DB:   getEnvInt("REDIS_DB", 0),
			TTL:  time.Duration(getEnvInt("REDIS_TTL_MINUTES", 10)) * time.Minute,
		},
	}
	
	config.PrintConfig()

	return config, nil
}


func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}