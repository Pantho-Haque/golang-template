package beckon

import (
	"sync"

	// this package is necessary to read config from remote consul
	_ "github.com/spf13/viper/remote"
)

type statsdContext struct {
	url string
}

var mu sync.Mutex
var statsd *statsdContext

// SetStatsd ...
func SetStatsd(url string) {
	mu.Lock()
	defer mu.Unlock()
	if url != "" {
		statsd = &statsdContext{
			url: url,
		}
	} else {
		statsd = nil
	}
}

// GetStatsdURL ...
func GetStatsdURL() string {
	mu.Lock()
	defer mu.Unlock()
	if statsd != nil {
		return statsd.url
	}
	return ""
}
