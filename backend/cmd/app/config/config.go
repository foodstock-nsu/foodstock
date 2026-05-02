package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type adminSeed struct {
	Login    string
	Password string
}

type Config struct {
	// Postgres
	DbHost     string `env:"DB_HOST,required"`
	DbPort     int    `env:"DB_PORT" envDefault:"5432"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DbSSLMode  string `env:"DB_SSL_MODE" envDefault:"prefer"`

	DbMaxConn         int           `env:"DB_MAX_CONNECTIONS" envDefault:"30"`
	DbMinConn         int           `env:"DB_MIN_CONNECTIONS" envDefault:"10"`
	DbMaxConnLifeTime time.Duration `env:"DB_MAX_CONNECTION_LIFETIME" envDefault:"10m"`
	DbMaxConnIdleTime time.Duration `env:"DB_MAX_CONNECTION_IDLETIME" envDefault:"5m"`

	// Auth constants
	AuthSecret string        `env:"AUTH_SECRET,required"`
	AuthTTL    time.Duration `env:"AUTH_TTL,required"`
	AdminSeeds []string      `env:"ADMIN_SEEDS,required" envSeparator:","`

	// Yookassa secrets
	YookassaShopID  string        `env:"YOOKASSA_SHOP_ID,required"`
	YookassaAPIKey  string        `env:"YOOKASSA_API_KEY,required"`
	YookassaTimeout time.Duration `env:"YOOKASSA_TIMEOUT" envDefault:"1m"`

	// Password hasher
	PasswordCost int `env:"PASSWORD_COST" envDefault:"10"`

	// QR-code generator
	QRCodeBaseURL string `env:"QR_CODE_BASE_URL,required"`
	QRCodeSize    int    `env:"QR_CODE_SIZE" envDefault:"512"`

	// Service
	HttpPort    int    `env:"HTTP_PORT" envDefault:"8080"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
}

func (c *Config) GetAdminSeeds() []adminSeed {
	var seeds []adminSeed
	for _, s := range c.AdminSeeds {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) == 2 {
			seeds = append(seeds, adminSeed{Login: parts[0], Password: parts[1]})
		}
	}
	return seeds
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   Environment: %s\n", cfg.Environment)
	fmt.Printf("   Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("   Postgres Host: %s\n", cfg.DbHost)
	fmt.Printf("   HTTP Port: %d\n", cfg.HttpPort)

	return cfg, nil
}

type TestConfig struct {
	// Postgres
	DbHost     string `env:"TEST_DB_HOST,required"`
	DbPort     int    `env:"TEST_DB_PORT" envDefault:"5433"`
	DbUser     string `env:"TEST_DB_USER,required"`
	DbPassword string `env:"TEST_DB_PASSWORD,required"`
	DBName     string `env:"TEST_DB_NAME,required"`
	DbSSLMode  string `env:"TEST_DB_SSL_MODE" envDefault:"prefer"`

	DbMaxConn         int           `env:"TEST_DB_MAX_CONNECTIONS" envDefault:"30"`
	DbMinConn         int           `env:"TEST_DB_MIN_CONNECTIONS" envDefault:"10"`
	DbMaxConnLifeTime time.Duration `env:"TEST_DB_MAX_CONNECTION_LIFETIME" envDefault:"10m"`
	DbMaxConnIdleTime time.Duration `env:"TEST_DB_MAX_CONNECTION_IDLETIME" envDefault:"5m"`
}

func LoadTest() (*TestConfig, error) {
	cfg := &TestConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load test config: %v", err)
	}

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   Postgres Host: %s\n", cfg.DbHost)

	return cfg, nil
}
