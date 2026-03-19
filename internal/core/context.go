package core

import (
	"time"

	"go.uber.org/zap"
	"pantho/golang/internal/config"
	"pantho/golang/internal/stores"
)

type Ctx struct {
	Store *stores.StoreHolder
	// Cache      cache.Cache
	Config *config.Config
	Log    *zap.Logger
	Now    time.Time
}

func NewCtx(
	store *stores.StoreHolder,
	config *config.Config,
	logger *zap.Logger,
) *Ctx {
	return &Ctx{
		Store:  store,
		Config: config,
		Log:    logger,
		Now:    time.Now(),
	}
}
