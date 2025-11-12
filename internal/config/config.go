package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env                string        `mapstructure:"env"`
	Port               int           `mapstructure:"port"`
	RequestPerMinute   int           `mapstructure:"requests_per_minute"`
	AccessTokenSecret  string        `mapstructure:"access_token_secret"`
	RefreshTokenSecret string        `mapstructure:"refresh_token_secret"`
	Timeout            time.Duration `mapstructure:"timeout"`
	Event              Event         `mapstructure:"event"`
	Auth               Auth          `mapstructure:"auth"`
}

type Event struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type Auth struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func New() *Config {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	BindEnv()

	viper.AutomaticEnv()
	viper.SetDefault("env", env)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}

func BindEnv() {
	viper.BindEnv("event.port", "EVENT_PORT")
	viper.BindEnv("event.host", "EVENT_HOST")

	viper.BindEnv("auth.port", "AUTH_PORT")
	viper.BindEnv("auth.host", "AUTH_HOST")

	viper.BindEnv("access_token_secret", "ACCESS_TOKEN_SECRET")
	viper.BindEnv("refresh_token_secret", "REFRESH_TOKEN_SECRET")
}
