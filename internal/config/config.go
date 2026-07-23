package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config interface {
	Database() DatabaseConfig
	HandlerType() string
	Pagination() PaginationConfig
	PriorityPolicy() PriorityPolicyConfig
	RequestTimeout() RequestTimeoutConfig
}

type EnvironmentConfig struct{}

type PaginationConfig struct {
	DefaultLimit string
	MaxLimit     string
}

type PriorityPolicyConfig struct {
	UsePolicy           string
	NegativeStockFactor string
	LeadTimeFactor      string
	ZeroSalesFactor     string
}

type DatabaseConfig struct {
	Driver           string
	ConnectionString string
}

type RequestTimeoutConfig struct {
	Enabled  string
	Duration string
}

func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env not found, using system environment variables")
	}

	return EnvironmentConfig{}
}

func (EnvironmentConfig) Database() DatabaseConfig {
	return DatabaseConfig{
		Driver:           os.Getenv("DATABASE_DRIVER"),
		ConnectionString: os.Getenv("DATABASE_URL"),
	}
}

func (EnvironmentConfig) HandlerType() string {
	return os.Getenv("HANDLER_TYPE")
}

func (EnvironmentConfig) Pagination() PaginationConfig {
	return PaginationConfig{
		DefaultLimit: os.Getenv("PAGINATION_DEFAULT_LIMIT"),
		MaxLimit:     os.Getenv("PAGINATION_MAX_LIMIT"),
	}
}

func (EnvironmentConfig) PriorityPolicy() PriorityPolicyConfig {
	return PriorityPolicyConfig{
		UsePolicy:           os.Getenv("PRIORITY_USE_POLICY"),
		NegativeStockFactor: os.Getenv("PRIORITY_NEGATIVE_STOCK_FACTOR"),
		LeadTimeFactor:      os.Getenv("PRIORITY_LEAD_TIME_FACTOR"),
		ZeroSalesFactor:     os.Getenv("PRIORITY_ZERO_SALES_FACTOR"),
	}
}

func (EnvironmentConfig) RequestTimeout() RequestTimeoutConfig {
	enabled := os.Getenv("REQUEST_TIMEOUT_ENABLED")
	if enabled == "" {
		enabled = "true"
	}

	duration := os.Getenv("REQUEST_TIMEOUT")
	if duration == "" {
		duration = "5s"
	}

	return RequestTimeoutConfig{Enabled: enabled, Duration: duration}
}
