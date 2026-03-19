package api

import (
	"github.com/gin-gonic/gin"
	"pantho/golang/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine, uh handlers.UserHandler) *gin.RouterGroup {

	root := r.Group("/")
	root.GET("healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := root.Group("v1")
	{
		v1.GET("/users", uh.GetUsers)

	}

	return root
}
