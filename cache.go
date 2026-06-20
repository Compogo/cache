package cache

import (
	"fmt"

	"github.com/Compogo/compogo"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/metrics"
)

// Cache — интерфейс кэша для хранения байтовых данных.
// Обёртка над gocache.CacheInterface[[]byte] с поддержкой метрик Prometheus.
//
// Используется для кэширования данных любого типа (сериализованных в JSON, Protobuf и т.д.).
//
// Пример:
//
//	var cacheInstance cache.Cache
//	container.Invoke(func(c cache.Cache) { cacheInstance = c })
//
//	// Сохранение
//	cacheInstance.Set(ctx, "user:123", []byte(`{"name":"John"}`), store.WithExpiration(time.Minute))
//
//	// Чтение
//	data, err := cacheInstance.Get(ctx, "user:123")
type Cache cache.CacheInterface[[]byte]

// NewCache создаёт новый экземпляр кэша с указанным драйвером.
// Автоматически добавляет метрики Prometheus для мониторинга.
//
// Процесс создания:
//  1. Определяет драйвер из конфигурации
//  2. Получает фабричную функцию (getter) для этого драйвера
//  3. Создаёт store через getter (использует DI-контейнер)
//  4. Оборачивает store в cache с метриками
func NewCache(config *Config, appConfig *compogo.Config, logger compogo.Logger, container compogo.Container) (Cache, error) {
	getter, err := getters.Get(config.Driver)
	if err != nil {
		return nil, fmt.Errorf("[cache] driver '%s' getter undefined: %w", config.Driver, err)
	}

	cacheStore, err := getter(container)
	if err != nil {
		return nil, fmt.Errorf("[cache] driver '%s' create failed: %w", config.Driver, err)
	}

	logger.GetLogger("cache").Infof("usage driver - '%s'", config.Driver)

	return cache.NewMetric[[]byte](
		metrics.NewPrometheus(appConfig.Name),
		cache.New[[]byte](cacheStore),
	), nil
}
