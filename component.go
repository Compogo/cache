package cache

import (
	"fmt"
	"strings"

	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
)

// Component — компонент кэша для Compogo.
// Регистрирует конфигурацию и экземпляр кэша в DI-контейнере.
//
// Поддерживает различные драйверы хранения (Redis, Memcached, in-memory и т.д.),
// которые регистрируются через Registration().
//
// Пример подключения:
//
//	app.AddComponents(&cache.Component)
//
// Получение кэша в компоненте:
//
//	var c cache.Cache
//	container.Invoke(func(cache cache.Cache) { c = cache })
//
//	c.Set(ctx, "key", []byte("value"), store.WithExpiration(time.Minute))
var Component = compogo.Component{
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(
			NewConfig,
			NewCache,
		)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			allDrivers := drivers.Keys()
			if len(allDrivers) == 1 {
				DriverDefault = allDrivers[0]
			}

			flagSet.StringVar(&config.driver, DriverFieldName, DriverDefault, fmt.Sprintf("cache driver. Available drivers: [%s]", strings.Join(allDrivers, ",")))
			flagSet.DurationVar(&config.Expiration, ExpirationFieldName, ExpirationDefault, "default data retention time")
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
}
