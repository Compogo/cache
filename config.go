package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/Compogo/compogo"
)

const (
	// DriverFieldName имя драйвера кэша
	DriverFieldName = "cache.driver"

	// ExpirationFieldName время жизни данных по умолчанию
	ExpirationFieldName = "cache.expiration"
)

var (
	// DriverDefault драйвер по умолчанию (определяется при старте)
	DriverDefault = ""

	// ExpirationDefault время жизни по умолчанию
	ExpirationDefault = 5 * time.Minute
)

// Config содержит конфигурацию кэша.
type Config struct {
	driver     string
	Driver     Driver
	Expiration time.Duration
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
// Проверяет, что драйвер указан и зарегистрирован.
func Configuration(config *Config, configurator compogo.Configurator) (*Config, error) {
	if config.driver == "" || config.driver == DriverDefault {
		configurator.SetDefault(DriverFieldName, DriverDefault)
		config.driver = configurator.GetString(DriverFieldName)
	}

	if config.driver == "" {
		return nil, errors.New("[cache] driver is not set")
	}

	driver, err := drivers.Get(config.driver)
	if err != nil {
		return nil, fmt.Errorf("[cache] driver '%s' get failed: %w", config.driver, err)
	}

	config.Driver = driver

	if config.Expiration == 0 || config.Expiration == ExpirationDefault {
		configurator.SetDefault(ExpirationFieldName, ExpirationDefault)
		config.Expiration = configurator.GetDuration(ExpirationFieldName)
	}

	return config, nil
}
