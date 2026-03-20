// internal/cache/keys.go
package cache

import "fmt"

const (
	prefix = "template:"

	// user keys
	keyUser = "user:top10" // usage: prefix + keyUser formatted with userId
)


func (cs *CacheStore) keyBuilder(pattern string, values ...interface{}) string {
	return fmt.Sprintf(prefix+pattern, values...)
}
