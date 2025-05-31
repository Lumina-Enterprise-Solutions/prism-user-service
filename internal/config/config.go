package config

import (
	"os"
	"time"

	commonConfig "github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/config"
)

type Config struct {
	Service  ServiceConfig               `mapstructure:"service"`
	Database commonConfig.DatabaseConfig `mapstructure:"database"`
	Redis    commonConfig.RedisConfig    `mapstructure:"redis"`
	JWT      commonConfig.JWTConfig      `mapstructure:"jwt"`
	Server   ServerConfig                `mapstructure:"server"`
	Log      LogConfig                   `mapstructure:"log"`
}

type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
	baseConfig, err := commonConfig.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Service: ServiceConfig{
			Name:        getEnvString("SERVICE_NAME", "prism-user-service"),
			Version:     getEnvString("SERVICE_VERSION", "v1.0.0"),
			Environment: getEnvString("SERVICE_ENVIRONMENT", "development"),
		},
		Database: baseConfig.Database,
		Redis:    baseConfig.Redis,
		JWT:      baseConfig.JWT,
		Server: ServerConfig{
			Host:         getEnvString("SERVER_HOST", "0.0.0.0"),
			Port:         baseConfig.Server.Port,
			ReadTimeout:  time.Duration(baseConfig.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(baseConfig.Server.WriteTimeout) * time.Second,
		},
		Log: LogConfig{
			Level:  getEnvString("LOG_LEVEL", "info"),
			Format: getEnvString("LOG_FORMAT", "json"),
		},
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
