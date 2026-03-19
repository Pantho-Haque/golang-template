package config

import (
	"fmt"

	// "github.com/spf13/viper"
	// _ "github.com/spf13/viper/remote"
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
	// viper.BindEnv("consul_url")
	// viper.BindEnv("consul_path")
	// viper.BindEnv("consul_path_payment_type_discount")

	// consulUrl := viper.GetString("consul_url")
	// consulPath := viper.GetString("consul_path")

	// viper.SetConfigType("yaml")
	// viper.AddRemoteProvider("consul", consulUrl, consulPath)
	// err := viper.ReadRemoteConfig()
	// if err != nil {
	// 	return nil, err
	// }

	config := &Config{
		App: AppConfig{
			Port: "8080",
		},
		Postgres: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "pathao",
			Password: "ZjYxYzZiZTBhMzFkMTQ4ZGJmOTE4MmEz",
			DB:       "on_demand",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
	}
	// err = viper.UnmarshalKey("app", &config.App)
	// if err != nil {
	// 	return nil, err
	// }

	// err = viper.UnmarshalKey("postgres", &config.Postgres)
	// if err != nil {
	// 	return nil, err
	// }

	// err = viper.UnmarshalKey("redis", &config.Redis)
	// if err != nil {
	// 	return nil, err
	// }

	config.PrintConfig()

	return config, nil
}
