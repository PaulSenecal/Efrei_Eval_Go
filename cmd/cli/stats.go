//cmd\cli\stats.go
package cli

import (
"fmt"
"log"
"os"

cmd2 "github.com/axellelanca/urlshortener/cmd"
"github.com/axellelanca/urlshortener/internal/repository"
"github.com/axellelanca/urlshortener/internal/services"
"github.com/spf13/cobra"

"gorm.io/driver/sqlite" // Driver SQLite pour GORM
"gorm.io/gorm"
)


var shortCodeFlag string // variable pour le flag --code

var StatsCmd = &cobra.Command{
Use:   "stats",
Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
Run: func(cmd *cobra.Command, args []string) {
// Vérifier que le flag est bien fourni
if shortCodeFlag == "" {
fmt.Println("Erreur : vous devez fournir un code court avec --code")
os.Exit(1)
}

// Charger la configuration globale
cfg := cmd2.Cfg

// Connexion à la base de données
db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
if err != nil {
log.Fatalf("FATAL: échec de la connexion à la base de données: %v", err)
}

sqlDB, err := db.DB()
if err != nil {
log.Fatalf("FATAL: échec de l'obtention de la base de données SQL sous-jacente: %v", err)
}
defer sqlDB.Close()

// Initialiser repository et service
linkRepo := repository.NewLinkRepository(db)
linkService := services.NewLinkService(linkRepo)

// Appel à GetLinkStats
link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
if err != nil {
if err == gorm.ErrRecordNotFound {
fmt.Printf("Erreur : Aucun lien trouvé pour le code '%s'\n", shortCodeFlag)
} else {
fmt.Printf("Erreur lors de la récupération des statistiques : %v\n", err)
}
os.Exit(1)
}

fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
fmt.Printf("URL longue: %s\n", link.LongURL)
fmt.Printf("Total de clics: %d\n", totalClicks)
},
}

func init() {
// Définir le flag --code
StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Code court pour lequel afficher les statistiques")

// Marquer le flag comme requis
err := StatsCmd.MarkFlagRequired("code")
if err != nil {
log.Fatalf("Erreur lors du marquage du flag --code comme requis: %v", err)

// Ajouter la commande à RootCmd
cmd2.RootCmd.AddCommand(StatsCmd)
}