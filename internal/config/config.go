//internal\config\config.go
package config

import (
	"fmt"
	"log" // Pour logger les informations ou erreurs de chargement de config

	"github.com/spf13/viper" // La bibliothèque pour la gestion de configuration
)

type ApplicationConfiguration struct {
    Server     ServerConfiguration
    Database   DatabaseConfiguration
    Monitoring MonitoringConfiguration
    Analytics  AnalyticsConfiguration
    ShortURL   ShortURLConfiguration
}

type ServerConfiguration struct {
    Port int
    Host string
}

type DatabaseConfiguration struct {
    Path string
}

type MonitoringConfiguration struct {
    CheckIntervalMinutes int
}

type AnalyticsConfiguration struct {
    BufferSize   int
    WorkerCount  int
}

type ShortURLConfiguration struct {
    CodeLength int
    BaseURL    string
}

func LoadApplicationConfiguration() *ApplicationConfiguration {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    viper.AddConfigPath(".")
    
    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Impossible de lire le fichier de configuration: %v", err)
        return getDefaultConfiguration()
    }
    
    var applicationConfiguration ApplicationConfiguration
    if err := viper.Unmarshal(&applicationConfiguration); err != nil {
        log.Printf("Impossible de parser la configuration: %v", err)
        return getDefaultConfiguration()
    }
    
    return &applicationConfiguration
}

func getDefaultConfiguration() *ApplicationConfiguration {
    return &ApplicationConfiguration{
        Server: ServerConfiguration{
            Port: 8080,
            Host: "localhost",
        },
        Database: DatabaseConfiguration{
            Path: "url_shortener.db",
        },
        Monitoring: MonitoringConfiguration{
            CheckIntervalMinutes: 5,
        },
        Analytics: AnalyticsConfiguration{
            BufferSize:  1000,
            WorkerCount: 5,
        },
        ShortURL: ShortURLConfiguration{
            CodeLength: 6,
            BaseURL:    "http://localhost:8080",
        },
    }
}