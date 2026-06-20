package cache

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/mapper"
	"github.com/eko/gocache/lib/v4/store"
)

var (
	// drivers — маппер для хранения доступных драйверов кэша.
	// Используется для отображения имени драйвера на его тип.
	drivers = mapper.NewMapper[Driver]()

	// getters — линкер для хранения фабричных функций драйверов.
	// Используется для создания store по типу драйвера.
	getters = linker.NewLinker[Driver, Getter]()
)

// Registration регистрирует драйвер кэша и его фабричную функцию.
// Должна вызываться в init() каждого пакета драйвера.
//
// Пример регистрации Redis-драйвера:
//
//	func init() {
//	    cache.Registration(cache.Driver("redis"), func(container compogo.Container) (store.StoreInterface, error) {
//	        var redisClient *redis.Client
//	        container.Invoke(func(r *redis.Client) { redisClient = r })
//	        return redis.NewStore(redisClient), nil
//	    })
//	}
func Registration(d Driver, getter Getter) {
	drivers.Add(d)
	getters.Add(d, getter)
}

// Getter — фабричная функция для создания store кэша.
// Принимает DI-контейнер для получения зависимостей драйвера.
type Getter func(container compogo.Container) (store.StoreInterface, error)

// Driver — тип драйвера кэша (например, "redis", "memcached", "memory").
type Driver string

// String возвращает строковое представление драйвера.
func (d Driver) String() string {
	return string(d)
}
