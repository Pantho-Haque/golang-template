// internal/cache/methods.go
package cache

import (
	"encoding/json"
	"errors"
	"pantho/golang/internal/models"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// UserCacheData — only the fields you want to cache, not the full DB model
// type UserCacheData struct {
// 	ID    int32  `json:"id"`
// 	Name  string `json:"name"`
// 	Email string `json:"email"`
// }

// Cache is the interface your app depends on — add methods here as you implement them
type CacheMethods interface {
	GetUser() (*[]models.User, error)
	SetUser(user *[]models.User) error
	DeleteUser() error
}


// ---- User cache methods ----
func (cs *CacheStore) GetUser() (*[]models.User, error) {
	key := cs.keyBuilder(keyUser)

	data, err := cs.client.Get(cs.ctx(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss — not an error
		}
		cs.log.Error("cache get user failed", zap.Error(err))
		return nil, err
	}

	var user []models.User
	err = json.Unmarshal([]byte(data), &user);
	if err != nil {
		cs.log.Error("cache unmarshal user failed", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (cs *CacheStore) SetUser(users *[]models.User) error {
	key := cs.keyBuilder(keyUser)

	bytes, err := json.Marshal(users)
	if err != nil {
		return err
	}

	ttl := cs.cfg.Redis.TTL
	err = cs.client.Set(cs.ctx(), key, string(bytes), ttl).Err()
	if  err != nil {
		cs.log.Error("cache set user failed", zap.Error(err))
		return err
	}

	return nil
}

func (cs *CacheStore) DeleteUser() error {
	key := cs.keyBuilder(keyUser)
	err := cs.client.Del(cs.ctx(), key).Err()
	if err != nil {
		cs.log.Error("cache delete user failed", zap.Error(err))
		return err
	}
	return nil
}