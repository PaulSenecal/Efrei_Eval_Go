package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port    int    `mapstructure:"port"`
		Host    string `mapstructure:"host"`
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"server"`
	
	Database struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"database"`
	
	Analytics struct {
		BufferSize  int `mapstructure:"buffer_size"`
		WorkerCount int `mapstructure:"worker_count"`
	} `mapstructure:"analytics"`
	
	Monitor struct {
		IntervalMinutes int `mapstructure:"interval_minutes"`
	} `mapstructure:"monitor"`
	
	ShortURL struct {
		CodeLength int    `mapstructure:"code_length"`
		BaseURL    string `mapstructure:"base_url"`
	} `mapstructure:"short_url"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// Valeurs par défaut
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.base_url", "http://localhost:8080")
	viper.SetDefault("database.name", "url_shortener.db")
	viper.SetDefault("analytics.buffer_size", 1000)
	viper.SetDefault("analytics.worker_count", 5)
	viper.SetDefault("monitor.interval_minutes", 5)
	viper.SetDefault("short_url.code_length", 6)
	viper.SetDefault("short_url.base_url", "http://localhost:8080")
	
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Impossible de lire le fichier de configuration: %v", err)
	}
	
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	
	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)
	
	return &cfg, nil
}