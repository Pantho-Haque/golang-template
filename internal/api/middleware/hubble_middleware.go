package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	ctx "magic.pathao.com/beatbox/hubble/context"
)

type context struct {
	*ctx.Context
}

func (c *context) middleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		start := time.Now()

		// Process the request
		ginCtx.Next()

		// After request
		path := ginCtx.FullPath()
		method := ginCtx.Request.Method
		status := ginCtx.Writer.Status()

		// Push metrics
		c.Push(path, method, status, start)
	}
}

func Wrap(prom *ctx.Prom) gin.HandlerFunc {
	c := &context{prom.GetContext()}
	return c.middleware()
}

func Use(prom *ctx.Prom) gin.HandlerFunc {
	c := &context{prom.GetContext()}
	return c.middleware()
}
