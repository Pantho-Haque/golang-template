package userService

import (
	"fmt"
	"pantho/golang/internal/core"
	"time"
)

type FetchUser struct {
	resCtx *ResponseCtx
}

func (fu *FetchUser) Do(ctx *core.Ctx, doCtx *core.DoCtx) error {
	
	cached, err := ctx.Cache.GetUser()
	if err == nil && cached != nil {
		fu.resCtx.Time = fmt.Sprintf("%s", time.Since(ctx.Now))
		fu.resCtx.Data = *cached
		return nil // cache hit — skip DB entirely
	}
	
	users, err := ctx.Store.UserStore.GetFirstTenUsers()
	if err != nil {
		return fmt.Errorf("user not found")
	}
	
	_ = ctx.Cache.SetUser(&users);
	
	fu.resCtx.Time = fmt.Sprintf("%s", time.Since(ctx.Now))
	fu.resCtx.Data = users
	return nil
}
