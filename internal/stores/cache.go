package stores

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type CacheStore struct {
	cache *redis.Client
	log   *zap.Logger
}

const (
	PREFIX        = "template:"
)


func (cs *CacheStore) keyBuilder(placeholder string, values ...interface{}) string {
	return fmt.Sprintf(PREFIX+placeholder, values...)
}

