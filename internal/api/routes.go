package api

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"magic.pathao.com/beatbox/hubble"
	"magic.pathao.com/parcel/prism/internal/api/handlers"
	hubblegin "magic.pathao.com/parcel/prism/internal/api/middleware"
)

func SetupRoutes(r *gin.Engine, dh handlers.DemoHandler) *gin.RouterGroup {
	r.Use(gin.Recovery())

	pprof.Register(r)

	p := ginprometheus.NewPrometheus("prism")
	p.Use(r)

	r.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

	prom := hubble.NewAPM("prism").Histogram(nil)
	r.Use(hubblegin.Wrap(prom))

	root := r.Group("/")
	root.GET("healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := root.Group("v1")
	{
		v1.GET("/test", dh.Test)
	}

	return root
}
