package config

import (
"fmt"
"log" // Pour logger les informations ou erreurs de chargement de config

"github.com/spf13/viper" // La bibliothèque pour la gestion de configuration
)

type Config struct {
Server struct {
Port    int    `mapstructure:"port"`
BaseURL string `mapstructure:"base_url"`
} `mapstructure:"server"`

Database struct {
DSN  string `mapstructure:"dsn"`
Name string `mapstructure:"name"`
} `mapstructure:"database"`

Analytics struct {
BufferSize  int `mapstructure:"buffer_size"`
WorkerCount int `mapstructure:"worker_count"`
} `mapstructure:"analytics"`

Monitor struct {
IntervalMinutes int `mapstructure:"interval_minutes"`
} `mapstructure:"monitor"`
}

func LoadConfig() (*Config, error) {
viper.AddConfigPath("configs") // Dossier relatif
viper.SetConfigName("config")  // Nom du fichier sans extension
viper.SetConfigType("yaml")    // Format

// Valeurs par défaut si le fichier est absent ou incomplet
viper.SetDefault("server.port", 8080)
viper.SetDefault("server.base_url", "http://localhost:8080")
viper.SetDefault("database.dsn", "file::memory:?cache=shared")
viper.SetDefault("database.name", "url_shortener.db")
viper.SetDefault("analytics.buffer_size", 100)
viper.SetDefault("analytics.worker_count", 5)
viper.SetDefault("monitor.interval_minutes", 5)

// Lecture du fichier (silencieuse si le fichier est absent, car on a des defaults)
if err := viper.ReadInConfig(); err != nil {
fmt.Println("Fichier de configuration non trouvé, utilisation des valeurs par défaut.")
}

var cfg Config
if err := viper.Unmarshal(&cfg); err != nil {
return nil, fmt.Errorf("erreur lors du démappage de la configuration : %w", err)
}

log.Printf("Configuration chargée: Port=%d, DB=%s, Buffer=%d, Interval=%dmin",
cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

return &cfg, nil
}