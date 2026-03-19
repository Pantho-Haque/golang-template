package handlers

import (
	"time"
	"pantho/golang/internal/core"
	"pantho/golang/pkg/utils"
	userService "pantho/golang/internal/services/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	GetUsers(c *gin.Context)
}
type userHandler struct {
	ctx *core.Ctx
}

func NewUserHandler(ctx *core.Ctx) UserHandler {
	return &userHandler{
		ctx: ctx,
	}
}

func (h *userHandler) GetUsers(c *gin.Context) {
	ctx := h.ctx;
	ctx.Now = time.Now()

	resCtx := userService.ResponseCtx{}
	if err := userService.New(&resCtx).Do(ctx, &core.DoCtx{}); err != nil {
		h.ctx.Log.Error("error sharing parcel", zap.Error(err))
		utils.ServeErr(c, err)
		return
	}

	utils.ServeData(c, resCtx)
}
