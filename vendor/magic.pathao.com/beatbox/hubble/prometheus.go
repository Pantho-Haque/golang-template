package hubble

import (
	"magic.pathao.com/beatbox/hubble/context"
)

func NewAPM(app string) *context.Prom {
	prom := context.NewProm(app)
	return prom
}
