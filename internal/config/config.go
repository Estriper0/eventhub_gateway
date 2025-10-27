package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Env     string   `mapstructure:"env"`
	Port    string   `mapstructure:"port"`
	Timeout int      `mapstructure:"timeout"`
	DB      Database `mapstructure:"database"`
}

type Database struct {
	DbHost     string `mapstructure:"dbhost"`
	DbPort     string `mapstructure:"dbport"`
	DbUser     string `mapstructure:"dbuser"`
	DbPassword string `mapstructure:"dbpassword"`
	DbName     string `mapstructure:"dbname"`
	SSLMode    string `mapstructure:"sslmode"`
}

func New() *Config {
	_, file, _, _ := runtime.Caller(0)

	path := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(file))), "configs")

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	BindEnv()

	viper.SetDefault("port", 8080)
	viper.SetDefault("timeout", 30)
	viper.SetDefault("database.sslmode", "disable")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	config.Env = env

	if err := Validate(&config); err != nil {
		panic(err)
	}

	return &config
}

func BindEnv() {
	viper.BindEnv("database.dbhost", "DB_HOST")
	viper.BindEnv("database.dbport", "DB_PORT")
	viper.BindEnv("database.dbuser", "DB_USER")
	viper.BindEnv("database.dbpassword", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
}

func Validate(config *Config) error {
	if config.DB.DbHost == "" {
		return errors.New("DB_HOST not exist")
	} else if config.DB.DbPort == "" {
		return errors.New("DB_PORT not exist")
	} else if config.DB.DbUser == "" {
		return errors.New("DB_USER not exist")
	} else if config.DB.DbPassword == "" {
		return errors.New("DB_PASSWORD not exist")
	} else if config.DB.DbName == "" {
		return errors.New("DB_NAME not exist")
	} else if config.DB.SSLMode == "" {
		return errors.New("DB_SSLMODE not exist")
	}
	return nil
}
