# Compogo Cache 📦

**Compogo Cache** — это гибкая и расширяемая система кэширования для Go, построенная поверх [gocache](https://github.com/eko/gocache). Поддерживает множество драйверов (Redis, BigCache, Ristretto и др.) через простую систему регистрации, полностью интегрируется с Compogo и автоматически добавляет Prometheus-метрики.

## 🚀 Установка

```bash
go get github.com/Compogo/cache
```

### 📦 Быстрый старт

```go
package main

import (
    "github.com/Compogo/compogo"
    "github.com/Compogo/cache"
    "github.com/Compogo/cache/redis"  // подключаем Redis-драйвер
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithOsSignalCloser(),
        cache.Component,           // базовый компонент кэша
        redis.Component,           // регистрируем Redis-драйвер
        compogo.WithComponents(
            userServiceComponent,
        ),
    )

    if err := app.Serve(); err != nil {
        panic(err)
    }
}

// Использование в сервисе
var userServiceComponent = &component.Component{
    Dependencies: component.Components{cache.Component},
    Execute: component.StepFunc(func(c container.Container) error {
        return c.Invoke(func(cache cache.CacheInterface[[]byte]) {
            service := &UserService{cache: cache}
            service.Register()
        })
    }),
}

type UserService struct {
    cache cache.CacheInterface[[]byte]
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Пытаемся достать из кэша
    data, err := s.cache.Get(ctx, fmt.Sprintf("user:%d", id))
    if err == nil {
        var user User
        json.Unmarshal(data, &user)
        return &user, nil
    }
    
    // Нет в кэше — грузим из БД
    user, err := s.db.LoadUser(id)
    if err != nil {
        return nil, err
    }
    
    // Кладём в кэш
    data, _ = json.Marshal(user)
    s.cache.Set(ctx, fmt.Sprintf("user:%d", id), data)
    
    return user, nil
}
```

### ✨ Возможности

#### 🎯 Множество драйверов

Кэш поддерживает любые бэкенды через систему регистрации:

```go
// Регистрация нового драйвера (в пакете драйвера)
func init() {
    cache.Registration("redis", NewRedisStore)
    cache.Registration("bigcache", NewBigCacheStore)
    cache.Registration("ristretto", NewRistrettoStore)
}
```

#### ⚙️ Конфигурация через флаги

```bash
./myapp \
    --cache.driver=redis \
    --cache.expiration=10m
```

Если зарегистрирован только один драйвер, он автоматически становится значением по умолчанию.

#### 🔧 Создание своего драйвера

```go
package mydriver

import (
    "github.com/Compogo/compogo/container"
    "github.com/Compogo/cache"
    "github.com/eko/gocache/lib/v4/store"
)

func init() {
    cache.Registration("mydriver", NewMyDriverStore)
}

func NewMyDriverStore(container container.Container) (store.StoreInterface, error) {
    // Достаём зависимости из контейнера
    var config *Config
    container.Invoke(func(cfg *Config) {
        config = cfg
    })
    
    // Создаём и возвращаем store
    return mydriver.NewStore(config), nil
}
```
