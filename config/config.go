package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	App
	DB
}

type App struct {
	Port  string `env:"APP_PORT,required"`
	Url   string `env:"APP_URL,required"`
	Debug bool   `env:"APP_DEBUG_MODE" envDefault:"false"`
}

type DB struct {
	Driver   string `env:"DB_DRIVER,required"`
	User     string `env:"DB_USERNAME,required"`
	Password string `env:"DB_PASS,required"`
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	Name     string `env:"DB_NAME,required"`
}

func Get() *Config {

	viper.AutomaticEnv()

	return &Config{
		App: App{
			Port:  getEnvAsString("APP_PORT", "", true),
			Url:   getEnvAsString("APP_URL", "", true),
			Debug: getEnvAsBool("APP_DEBUG_MODE", false, false),
		},
		DB: DB{
			Driver:   getEnvAsString("DB_DRIVER", "", true),
			User:     getEnvAsString("DB_USERNAME", "", true),
			Password: getEnvAsString("DB_PASS", "", true),
			Host:     getEnvAsString("DB_HOST", "", true),
			Port:     getEnvAsString("DB_PORT", "", true),
			Name:     getEnvAsString("DB_NAME", "", true),
		},
	}
}

func getEnvAsString(name string, defaultVal string, required bool) string {
	rawEnv := viper.Get(name)

	if rawEnv == nil {
		if required {
			log.Printf("[FATAL]: env var with name %s is not specified\n", name)
			os.Exit(1)
		}
		return defaultVal
	}

	val, ok := rawEnv.(string)
	if !ok {
		log.Printf("[FATAL]: env var with name %s cannot be converted to necessary type string\n", name)
		os.Exit(1)
	}

	return val
}

func getEnvAsBool(name string, defaultVal, required bool) bool {
	rawEnv := viper.Get(name)

	if rawEnv == nil {
		if required {
			log.Printf("[FATAL]: env var with name %s is not specified\n", name)
			os.Exit(1)
		}
		return defaultVal
	}

	val, err := strconv.ParseBool(rawEnv.(string))
	if err != nil {
		log.Printf("[FATAL]: env var with name %s cannot be converted to necessary type bool\n", name)
		os.Exit(1)
	}

	return val
}
