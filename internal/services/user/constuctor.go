package userService

import "pantho/golang/internal/core"

type user struct {
	ctx *ResponseCtx
}

func New(ctx *ResponseCtx) core.Doer {
	return &user{
		ctx: ctx,
	}
}

func (s *user) Do(ctx *core.Ctx, doCtx *core.DoCtx) error {
	doers := core.Doers{
		&FetchUser{resCtx: s.ctx},
	}

	return doers.Do(ctx, doCtx)
}