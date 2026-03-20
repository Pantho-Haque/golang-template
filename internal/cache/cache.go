// internal/cache/cache.go
package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"pantho/golang/internal/config"
)

func NewCacheStore(client *redis.Client, log *zap.Logger, cfg *config.Config) CacheMethods {
	return &CacheStore{
		client: client,
		log:    log,
		cfg:    cfg,
	}
}


type CacheStore struct {
	client *redis.Client
	log    *zap.Logger
	cfg    *config.Config
}


func (cs *CacheStore) ctx() context.Context {
	return context.Background()
}

