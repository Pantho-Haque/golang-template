package handlers

import "github.com/gin-gonic/gin"

type DemoHandler interface {
	Test(c *gin.Context)
}
type demoHandler struct {
}

func NewDemoHandler() DemoHandler {
	return &demoHandler{}
}

func (h *demoHandler) Test(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, world!"})
}
