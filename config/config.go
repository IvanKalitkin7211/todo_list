package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Database    DatabaseConfig
	Server      ServersConfig
	Logger      LoggerConfig
	Redis       RedisConfig
	RateLimiter RateLimiterConfig
	JWTSecret   string `mapstructure:"jwt_secret"`
}
type ServersConfig struct {
	HTTP HTTPConfig `mapstructure:"http"`
}

type HTTPConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}

type LoggerConfig struct {
	Mode     string `mapstructure:"mode"`
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"filepath"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RateLimiterConfig struct {
	Limit        int `mapstructure:"limit"`
	WindowSec    int `mapstructure:"windowSeconds"`
	Window       time.Duration
	Enabled      bool   `mapstructure:"enabled"`
	ErrorMessage string `mapstructure:"errorMessage"`
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] No .env file found. Using system environment variables.")
	}

	// Mapping
	_ = viper.BindEnv("server.http.host", "TODO_HTTP_HOST")
	_ = viper.BindEnv("server.http.port", "TODO_HTTP_PORT")
	_ = viper.BindEnv("database.host", "TODO_DATABASE_HOST")
	_ = viper.BindEnv("database.port", "TODO_DATABASE_PORT")
	_ = viper.BindEnv("database.user", "TODO_DATABASE_USER")
	_ = viper.BindEnv("database.password", "TODO_DATABASE_PASSWORD")
	_ = viper.BindEnv("database.dbname", "TODO_DATABASE_DBNAME")
	_ = viper.BindEnv("database.sslmode", "TODO_DATABASE_SSLMODE")
	_ = viper.BindEnv("database.maxidleconns", "TODO_DATABASE_MAXIDLECONNS")
	_ = viper.BindEnv("database.maxopenconns", "TODO_DATABASE_MAXOPENCONNS")
	_ = viper.BindEnv("database.connmaxlifetime", "TODO_DATABASE_CONNMAXLIFETIME")

	_ = viper.BindEnv("redis.address", "TODO_REDIS_ADDRESS")
	_ = viper.BindEnv("redis.password", "TODO_REDIS_PASSWORD")
	_ = viper.BindEnv("redis.db", "TODO_REDIS_DB")
	_ = viper.BindEnv("rate_limiter.enabled", "TODO_RATELIMIT_ENABLED")
	_ = viper.BindEnv("rate_limiter.limit", "TODO_RATELIMIT_LIMIT")
	_ = viper.BindEnv("rate_limiter.windowSeconds", "TODO_RATELIMIT_WINDOW_SECONDS")
	_ = viper.BindEnv("rate_limiter.errorMessage", "TODO_RATELIMIT_ERROR_MESSAGE")
	_ = viper.BindEnv("jwt_secret", "TODO_JWT_SECRET")

	// Cfg file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	// Read cfg file
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println("[DEBUG] Config file not found, fallback to .env")
		} else {
			log.Printf("[WARN] Failed to merge config file: %v", err)
		}

	}

	// Decode
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("[ERROR] Failed to load config: %v", err)
	}

	// time.Duration для Rate Limiter
	cfg.RateLimiter.Window = time.Duration(cfg.RateLimiter.WindowSec) * time.Second

	return cfg
}
