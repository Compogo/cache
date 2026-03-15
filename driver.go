package cache

import (
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/mapper"
	"github.com/eko/gocache/lib/v4/store"
)

var (
	// drivers stores all registered cache drivers by their string names.
	// Used for validation and listing available drivers.
	drivers = mapper.NewMapper[Driver]()

	// getters stores factory functions that create cache stores for each driver.
	// The getter receives a container to resolve dependencies (config, clients, etc.).
	getters = linker.NewLinker[Driver, Getter]()
)

// Registration registers a new cache driver with its factory function.
// This should be called from driver packages (e.g., redis, bigcache) during init().
//
// Example:
//
//	func init() {
//		cache.Registration("redis", NewRedisStore)
//	}
func Registration(d Driver, getter Getter) {
	drivers.Add(d)
	getters.Add(d, getter)
}

// Getter is a factory function that creates a cache store from a DI container.
// It receives the container to resolve any dependencies the store might need
// (configuration, connections, etc.) and returns a store.StoreInterface.
type Getter func(container container.Container) (store.StoreInterface, error)

// Driver is a type-safe identifier for cache backends.
// It implements fmt.Stringer for logging and display.
type Driver string

// String returns the driver name as a string.
func (d Driver) String() string {
	return string(d)
}
