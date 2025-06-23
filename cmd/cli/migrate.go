//cmd\cli\migrate.go
package cli

import (
	"fmt"
	"log"

	cmd2 "urlshortenerGroupe8/cmd"
	"urlshortenerGroupe8/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

var MigrateCommand = &cobra.Command{
    Use:   "migrate",
    Short: "Exécute les migrations de la base de données",
    Run:   executeMigrateCommand,
}

func executeMigrateCommand(cmd *cobra.Command, args []string) {
    applicationConfiguration := config.LoadApplicationConfiguration()
    
    databaseConnection, err := gorm.Open(sqlite.Open(applicationConfiguration.Database.Path), &gorm.Config{})
    if err != nil {
        log.Fatalf("Impossible de se connecter à la base de données: %v", err)
    }
    
    err = databaseConnection.AutoMigrate(&models.Link{}, &models.Click{})
    if err != nil {
        log.Fatalf("Erreur lors de la migration: %v", err)
    }
    
    fmt.Println("Migration de la base de données réussie")
}

