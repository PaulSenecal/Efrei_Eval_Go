package cli

import (
"fmt"
"log"

cmd2 "github.com/axellelanca/urlshortener/cmd"
"github.com/axellelanca/urlshortener/internal/models"
"github.com/spf13/cobra"
"gorm.io/driver/sqlite" // Driver SQLite pour GORM
"gorm.io/gorm"
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
Use:   "migrate",
Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
Run: func(cmd *cobra.Command, args []string) {
cfg := cmd2.Cfg

log.Printf("Configuration chargée: Port=%d, DB=%s, Buffer=%d, Interval=%dmin",
cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
if err != nil {
log.Fatalf("FATAL: Échec de connexion à la base de données: %v", err)
}

sqlDB, err := db.DB()
if err != nil {
log.Fatalf("FATAL: Échec de l'obtention de la base SQL sous-jacente: %v", err)
}
defer sqlDB.Close()

if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
log.Fatalf("FATAL: Échec de l'exécution des migrations: %v", err)
}

fmt.Println("Migrations de la base de données exécutées avec succès.")
},
}

func init() {
cmd2.RootCmd.AddCommand(MigrateCmd)
}