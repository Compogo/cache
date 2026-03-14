package cache

import (
	"fmt"

	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/logger"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/metrics"
)

// NewCache creates a new cache instance with the selected driver.
// It:
//   - Looks up the getter function for the configured driver
//   - Creates the underlying store using the getter
//   - Wraps it with Prometheus metrics (using the app name as a label)
//   - Returns a cache.CacheInterface that can be used throughout the application
//
// The returned cache is generic over []byte, but gocache's cache.CacheInterface
// can work with any serializable type through marshaling/unmarshaling.
func NewCache(config *Config, appConfig *compogo.Config, logger logger.Logger, container container.Container) (cache.CacheInterface[[]byte], error) {
	getter, err := getters.Get(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("[cache] driver '%s' getter undefined: %w", config.Driver, err)
	}

	cacheStore, err := getter(container)
	if err != nil {
		return nil, fmt.Errorf("[cache] driver '%s' create failed: %w", config.Driver, err)
	}

	logger.Infof("[cache] usage driver - '%s'", config.Driver)

	return cache.NewMetric[[]byte](
		metrics.NewPrometheus(appConfig.Name),
		cache.New[[]byte](cacheStore),
	), nil
}
