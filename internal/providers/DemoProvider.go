package providers

import (
	"go.uber.org/zap"
	"magic.pathao.com/parcel/prism/internal/config"
)

type DemoProvider interface {
	GetDemo() string
}
type demoProvider struct {
	cfg         *config.Config
	httpService *HttpProvider
	log         *zap.Logger
}

func NewDemoProvider(cfg *config.Config, httpService *HttpProvider, log *zap.Logger) DemoProvider {
	return &demoProvider{
		cfg:         cfg,
		httpService: httpService,
		log:         log,
	}
}

func (p *demoProvider) GetDemo() string {
	return "demo"
}
