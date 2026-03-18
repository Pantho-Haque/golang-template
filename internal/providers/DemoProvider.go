package providers

import (
	"go.uber.org/zap"
	"pantho/golang/internal/config"
)

type DemoProvider interface {
	GetDemo() string
}
type demoProvider struct {
	cfg         *config.Config
	log         *zap.Logger
}

func NewDemoProvider(cfg *config.Config, log *zap.Logger) DemoProvider {
	return &demoProvider{
		cfg:         cfg,
		log:         log,
	}
}

func (p *demoProvider) GetDemo() string {
	return "demo"
}
