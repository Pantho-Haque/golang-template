package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"pantho/golang/internal/api"
	"pantho/golang/internal/api/handlers"
	// "pantho/golang/internal/services"
	"pantho/golang/internal/config"
	"pantho/golang/internal/conn"
	"pantho/golang/internal/core"
	"pantho/golang/internal/cache"
	// "pantho/golang/internal/providers"
	"pantho/golang/internal/stores"
	"pantho/golang/pkg"
)

func main() {
	fx.New(
		fx.Options(
			conn.Module,
		),
		fx.Provide(
			// dependencies
			pkg.CustomLogger,
			config.LoadConfig,

			// connections

			// handlers
			handlers.NewUserHandler,

			// services
			// services.NewUserService,

			// providers
			// providers.NewHttpProvider,
			// providers.NewDemoProvider,

			// stores
			cache.NewCacheStore,
			// stores.NewParcelStore,
			stores.NewStoreHolder,

			// core
			core.NewCtx,

			GinHttpServer,
			api.SetupRoutes,
		),
		fx.Invoke(
			func(r *gin.RouterGroup, l *zap.Logger) {
				l.Info("routes registered")
			},
		),
	).Run()
}

func GinHttpServer(lc fx.Lifecycle, log *zap.Logger, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return r
}
