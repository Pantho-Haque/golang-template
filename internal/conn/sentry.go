package conn

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"magic.pathao.com/parcel/prism/internal/config"
)

type Sentry struct {
	Client *sentry.Client
}

func ConnectSentry(cfg *config.Config) error {
	conf := cfg.Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         conf.Key,
		Environment: conf.Environment,
		Debug:       conf.Debug,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to sentry: %w", err)
	}
	fmt.Println("Connected to sentry")
	return nil
}
