package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"go.uber.org/zap"
	"pantho/golang/internal/services"
)

type UserHandler interface {
	GetUsers(c *gin.Context)
}
type userHandler struct {
	userService *services.UserService
	log *zap.Logger
}

func NewUserHandler(
	userService *services.UserService,
	log *zap.Logger,
) UserHandler {
	return &userHandler{
		userService: userService,
		log: log,
	}
}

func (h *userHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		h.log.Error("failed to get users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}
	c.JSON(http.StatusOK, *users)
}
