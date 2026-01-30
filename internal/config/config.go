// Package config содержит структуры конфигурации сервиса.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config - корневая структура конфигурации.
type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Log      LogConfig      `yaml:"log"`
	Metrics  MetricsConfig  `yaml:"metrics"`
}

// AppConfig - настройки приложения.
type AppConfig struct {
	Name  string `yaml:"name"`
	Env   string `yaml:"env"`
	Debug bool   `yaml:"debug"`
}

// ServerConfig - настройки серверов.
type ServerConfig struct {
	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`
}

// GRPCConfig - настройки gRPC сервера.
type GRPCConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Addr возвращает адрес gRPC сервера.
func (c GRPCConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// HTTPConfig - настройки HTTP сервера.
type HTTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Addr возвращает адрес HTTP сервера.
func (c HTTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DatabaseConfig - настройки подключения к БД.
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Name            string        `yaml:"name"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	SSLMode         string        `yaml:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// DSN возвращает строку подключения к PostgreSQL.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SSLMode,
	)
}

// LogConfig - настройки логирования.
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// MetricsConfig - настройки метрик.
type MetricsConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

// Load загружает конфигурацию из файла.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	// Подставляем переменные окружения
	data = []byte(os.ExpandEnv(string(data)))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}