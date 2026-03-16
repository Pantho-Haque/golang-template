package conn

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"magic.pathao.com/parcel/prism/internal/config"
)

func ConnectRedis(cfg *config.Config, log *zap.Logger) (*redis.Client, error) {
	rc := cfg.Redis

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		DB:   rc.DB,
	})

	// Ping the Redis server to check if the connection is established
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("successfully connected to Redis")
	return client, nil
}
