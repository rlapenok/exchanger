package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// main configuration
type Config struct {
	DB     DBConfig     `envPrefix:"DB_"`
	Redis  RedisConfig  `envPrefix:"REDIS_"`
	Logger LoggerConfig `envPrefix:"LOGGER_"`
}

// logger configuration
type LoggerConfig struct {
	Level  string `env:"LEVEL" envDefault:"debug"`
	Format string `env:"FORMAT" envDefault:"json"`
	Output string `env:"OUTPUT" envDefault:"stdout"`
}

// database configuration
type DBConfig struct {
	Host              string `env:"HOST,required,notEmpty"`
	Port              int    `env:"PORT,required,notEmpty"`
	User              string `env:"USER,required,notEmpty"`
	Password          string `env:"PASSWORD,required,notEmpty"`
	DBName            string `env:"NAME,required,notEmpty"`
	Schema            string `env:"SCHEMA" envDefault:"public"`
	SSLMode           string `env:"SSL_MODE" envDefault:"disable"`
	ConnectionTimeout int    `env:"CONNECTION_TIMEOUT" envDefault:"5"`
}

// redis configuration
type RedisConfig struct {
	Host     string `env:"HOST,required,notEmpty"`
	Port     int    `env:"PORT,required,notEmpty"`
	Password string `env:"PASSWORD,required,notEmpty"`
	DB       string `env:"DB,required,notEmpty"`
}

// load config from environment variables
func Load() (*Config, error) {
	_ = godotenv.Load()

	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
