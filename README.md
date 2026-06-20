# Compogo Cache

[![Go Reference](https://pkg.go.dev/badge/github.com/Compogo/cache.svg)](https://pkg.go.dev/github.com/Compogo/cache)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Плагинная система кэширования для фреймворка [Compogo](https://github.com/Compogo/compogo).

Построена на основе [gocache](https://github.com/eko/gocache) и предоставляет:

* Единый интерфейс для работы с кэшем
* Поддержку различных драйверов (Redis, Memcached, in-memory)
* Автоматические метрики Prometheus
* Конфигурацию через флаги

## Установка

```shell
go get github.com/Compogo/cache
```

## Быстрый старт

```go
package main

import (
    "context"
    "time"

    "github.com/Compogo/cache"
    "github.com/Compogo/compogo"
    "github.com/eko/gocache/lib/v4/store"
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithComponents(&cache.Component),
    )

    app.AddComponents(&compogo.Component{
        Name: "my_service",
        Init: compogo.StepFunc(func(container compogo.Container) error {
            return container.Invoke(func(c cache.Cache) error {
                ctx := context.Background()
                
                // Сохранение
                err := c.Set(ctx, "key", []byte("value"), store.WithExpiration(time.Minute))
                if err != nil {
                    return err
                }
                
                // Чтение
                data, err := c.Get(ctx, "key")
                if err != nil {
                    return err
                }
                
                return nil
            })
        }),
    })

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

## Конфигурация

### Флаги командной строки

```shell
# Драйвер кэша
--cache.driver=redis

# Время жизни по умолчанию
--cache.expiration=5m
```

### Регистрация драйверов

```go
import "github.com/Compogo/cache"

// Регистрация кастомного драйвера
func init() {
    cache.Registration(cache.Driver("custom"), func(container compogo.Container) (store.StoreInterface, error) {
        var client *CustomClient
        container.Invoke(func(c *CustomClient) { client = c })
        return custom.NewStore(client), nil
    })
}
```

## Метрики Prometheus

Кэш автоматически собирает метрики:

```plantuml
# Метрики gocache
cache_hit_total{result="hit", cache_type="...", app="myapp"}
cache_miss_total{result="miss", cache_type="...", app="myapp"}
cache_set_success_total{result="success", cache_type="...", app="myapp"}
cache_set_error_total{result="error", cache_type="...", app="myapp"}
cache_delete_success_total{result="success", cache_type="...", app="myapp"}
cache_delete_error_total{result="error", cache_type="...", app="myapp"}
```

## Зависимости

* [Compogo](https://github.com/Compogo/compogo) — основной фреймворк
* [gocache](https://github.com/eko/gocache) — библиотека кэширования
* [Prometheus](https://github.com/prometheus/client_golang) — метрики

## Лицензия

```plantuml
MIT License

Copyright (c) 2026 Compogo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
