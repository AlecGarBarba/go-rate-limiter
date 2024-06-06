package config

import (
	"log"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
type Configuration struct {
	APIUrl *url.URL
	Redis  RedisConfig
}

func LoadConfig() (Configuration, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local" // default to local if no environment is set
	}

	// Initialize Viper
	viper.SetConfigName("config." + env)
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}

	apiUrl, ok := viper.Get("API_URL").(string)
	if !ok {
		apiUrl = "http://localhost:3000"
	}

	backendURL, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatalf("Error parsing API URL: %v", err)
	}

	return Configuration{
		APIUrl: backendURL,
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	}, nil
}
