package core

import (
	"time"

	"go.uber.org/zap"
	"pantho/golang/internal/cache"
	"pantho/golang/internal/config"
	"pantho/golang/internal/stores"
)

type Ctx struct {
	Store  *stores.StoreHolder
	Cache  cache.CacheMethods
	Config *config.Config
	Log    *zap.Logger
	Now    time.Time
}

func NewCtx(
	store *stores.StoreHolder,
	cache cache.CacheMethods,
	config *config.Config,
	logger *zap.Logger,
) *Ctx {
	return &Ctx{
		Store:  store,
		Cache:  cache,
		Config: config,
		Log:    logger,
		Now:    time.Now(),
	}
}
