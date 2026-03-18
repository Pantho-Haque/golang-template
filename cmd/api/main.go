package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"pantho/golang/internal/api"
	"pantho/golang/internal/api/handlers"
	// "pantho/golang/internal/config"
	// "pantho/golang/internal/conn"
	// "pantho/golang/internal/providers"
	"pantho/golang/internal/services"
	// "pantho/golang/internal/stores"
	"pantho/golang/pkg"
)

func main() {
	fx.New(
		fx.Provide(
			// dependencies
			pkg.CustomLogger,
			// config.LoadConfig,

			// connections
			// conn.ConnectPostgres,
			// conn.ConnectRedis,

			// handlers
			handlers.NewDemoHandler,

			// services
			services.NewPalantirService,

			// providers
			// providers.NewHttpProvider,
			// providers.NewDemoProvider,

			// stores
			// stores.NewCacheStore,
			// stores.NewParcelStore,
			// stores.NewUserStore,

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

func GinHttpServer(lc fx.Lifecycle, log *zap.Logger) *gin.Engine {
	r := gin.Default()

	srv := &http.Server{
		Addr:    ":" + "8080",
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
