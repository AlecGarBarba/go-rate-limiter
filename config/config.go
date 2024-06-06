package config

import (
	"log"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	APIUrl *url.URL
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
	}, nil
}
