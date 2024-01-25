package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	API      APIConfig
	Server   ServerConfig
	Postgres PostgresConfig
}

type APIConfig struct {
	AgeURL         string
	GenderURL      string
	NationalityURL string
}

type ServerConfig struct {
	Host        string
	Port        int
	RequestTime time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func Init() (*Config, error) {
	viper.SetConfigFile("../../.env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		APIConfig{
			AgeURL:         viper.GetString("AGE_URL"),
			GenderURL:      viper.GetString("GENDER_URL"),
			NationalityURL: viper.GetString("NATIONALITY_URL"),
		},
		ServerConfig{
			Host:        viper.GetString("SERVER_HOST"),
			Port:        viper.GetInt("SERVER_PORT"),
			RequestTime: viper.GetDuration("REQUEST_TIME"),
		},
		PostgresConfig{
			Host:     viper.GetString("PG_HOST"),
			Port:     viper.GetInt("PG_PORT"),
			Username: viper.GetString("PG_USER"),
			Password: viper.GetString("PG_PASS"),
			DBName:   viper.GetString("PG_BASE"),
			SSLMode:  viper.GetString("PG_SSL_MODE"),
		},
	}, nil
}
