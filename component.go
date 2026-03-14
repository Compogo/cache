package cache

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
)

// Component is a ready-to-use Compogo component that provides a configurable cache.
// It automatically:
//   - Registers Config and Cache factory in the DI container
//   - Adds command-line flags for driver selection and TTL
//   - Validates the selected driver during Configuration phase
//   - Sets the default driver automatically if only one is registered
//
// Usage:
//
//	compogo.WithComponents(
//	    cache.Component,
//	    // ... driver components (redis, bigcache, etc.)
//	)
//
// Then in your service:
//
//	type Service struct {
//	    cache cache.CacheInterface[[]byte]
//	}
var Component = &component.Component{
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			NewConfig,
			NewCache,
		)
	}),
	BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
		return container.Invoke(func(config *Config) {
			allDrivers := drivers.Keys()
			if len(allDrivers) == 1 {
				DriverDefault = allDrivers[0]
			}

			flagSet.StringVar(&config.driver, DriverFieldName, DriverDefault, fmt.Sprintf("cache driver. Available drivers: [%s]", strings.Join(allDrivers, ",")))
			flagSet.DurationVar(&config.Expiration, ExpirationFieldName, ExpirationDefault, "default data retention time")
		})
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
}
