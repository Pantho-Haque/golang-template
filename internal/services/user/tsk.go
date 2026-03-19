package userService

import (
	"fmt"
	"pantho/golang/internal/core"
)

type FetchUser struct {
	resCtx *ResponseCtx
}

func (fu *FetchUser) Do(ctx *core.Ctx, doCtx *core.DoCtx) error {
	user, err := ctx.Store.UserStore.GetFirstTenUsers()
	if err != nil {
		return fmt.Errorf("user not found")
	}
	fu.resCtx.Name = "name"
	fu.resCtx.Data = user
	return nil
}
