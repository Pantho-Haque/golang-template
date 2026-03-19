package core

type DoCtx struct {
	IsExit  bool
	NxtDoer Doer
}

type Doer interface {
	Do(ctx *Ctx, doCtx *DoCtx) error
}

type Doers []Doer

func (ds Doers) Do(ctx *Ctx, doCtx *DoCtx) error {
	for _, d := range ds {
		if doCtx != nil && doCtx.NxtDoer != nil && doCtx.NxtDoer != d {
			continue
		}
		doCtx.NxtDoer = nil
		err := d.Do(ctx, doCtx)
		if err != nil {
			return err
		}
		if doCtx != nil && doCtx.IsExit {
			break
		}
	}
	return nil
}
