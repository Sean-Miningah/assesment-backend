package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Environment    string `mapstructure:"ENVIRONMENT"`
	Address        string `mapstructure:"SERVER_ADDRESS"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	ServiceName    string `mapstructure:"TELEMETRY_SERVICE_NAME"`
	JaegerEndpoint string `mapstructure:"JAEGER_ENDPOINT"`
	PrometheusPort string `mapstructure:"PROMETHEUS_PORT"`

	//Google client
	GoogleClientID     string `mapstructure:"CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"CLIENT_SECRET"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"REDIRECT_URL"`

	// AT API
	ATAPIKey             string `mapstructure:"ATAPI_KEY"`
	NotificationUsername string `mapstructure:"NOTIFICATION_USERNAME"`
	GmailAppAPIKey       string `mapstructure:"GMAIL_APP_API_KEY"`
	ATAPIUrl             string `mapstructure:"ATAPI_URL"`
}

func Load(path string) *Config {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)

	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &config
}
