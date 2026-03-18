package api

import (
	"github.com/gin-gonic/gin"
	"pantho/golang/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine, dh handlers.DemoHandler) *gin.RouterGroup {
	r.Use(gin.Recovery())

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
