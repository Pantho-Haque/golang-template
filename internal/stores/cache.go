package stores

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type CacheStore interface {
	keyBuilder(placeholder string, values ...interface{}) string
	keyBuilderParcel(placeholder string, values ...interface{}) string
}
type cacheStore struct {
	cache *redis.Client
	log   *zap.Logger
}

const (
	PREFIX        = "prism:"
	PREFIX_PARCEL = "parcels:"
)

func NewCacheStore(cache *redis.Client, log *zap.Logger) CacheStore {
	return &cacheStore{
		cache: cache,
		log:   log,
	}
}

func (cs *cacheStore) keyBuilder(placeholder string, values ...interface{}) string {
	return fmt.Sprintf(PREFIX+placeholder, values...)
}

func (cs *cacheStore) keyBuilderParcel(placeholder string, values ...interface{}) string {
	return fmt.Sprintf(PREFIX_PARCEL+placeholder, values...)
}
