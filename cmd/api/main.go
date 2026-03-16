package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"magic.pathao.com/parcel/prism/internal/api"
	"magic.pathao.com/parcel/prism/internal/api/handlers"
	"magic.pathao.com/parcel/prism/internal/config"
	"magic.pathao.com/parcel/prism/internal/conn"
	"magic.pathao.com/parcel/prism/internal/providers"
	"magic.pathao.com/parcel/prism/internal/services"
	"magic.pathao.com/parcel/prism/internal/stores"
	"magic.pathao.com/parcel/prism/pkg"
)

func main() {
	fx.New(
		fx.Provide(
			pkg.CustomLogger,
			config.LoadConfig,
			conn.ConnectPostgres,
			conn.ConnectRedis,
			conn.ConnectBeatbox,
			handlers.NewDemoHandler,

			// services
			services.NewPalantirService,
			services.NewSherlockService,

			// providers
			providers.NewHttpProvider,
			// providers.NewDemoProvider,

			// stores
			stores.NewCacheStore,
			stores.NewParcelStore,
			stores.NewUserStore,

			GinHttpServer,
			api.SetupRoutes,
		),
		fx.Invoke(
			func(r *gin.RouterGroup, l *zap.Logger) {
				l.Info("routes registered")
			},
			// InitSentry
			conn.ConnectSentry,
		),
	).Run()
}

func GinHttpServer(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger) *gin.Engine {
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
