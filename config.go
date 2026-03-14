package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/Compogo/compogo/configurator"
)

const (
	DriverFieldName     = "cache.driver"
	ExpirationFieldName = "cache.expiration"

	ExpirationDefault = 5 * time.Minute
)

var DriverDefault = ""

type Config struct {
	driver     string
	Driver     Driver
	Expiration time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func Configuration(config *Config, configurator configurator.Configurator) (*Config, error) {
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
