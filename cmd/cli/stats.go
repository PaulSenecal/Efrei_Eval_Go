//cmd\cli\stats.go
package cli

import (
	"fmt"
	"log"
	"os"

	cmd2 "urlshortenerGroupe8/cmd"
	"urlshortenerGroupe8/internal/repository"
	"urlshortenerGroupe8/internal/services"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

var StatsCommand = &cobra.Command{
    Use:   "stats",
    Short: "Affiche les statistiques d'un lien",
    Run:   executeStatsCommand,
}

var shortCodeToAnalyze string

func init() {
    StatsCommand.Flags().StringVar(&shortCodeToAnalyze, "code", "", "Code court du lien")
    StatsCommand.MarkFlagRequired("code")
}

func executeStatsCommand(cmd *cobra.Command, args []string) {
    applicationConfiguration := config.LoadApplicationConfiguration()
    
    databaseConnection, err := gorm.Open(sqlite.Open(applicationConfiguration.Database.Path), &gorm.Config{})
    if err != nil {
        log.Fatalf("Impossible de se connecter à la base de données: %v", err)
    }
    
    linkRepository := repository.NewLinkRepository(databaseConnection)
    linkService := services.NewLinkService(linkRepository)
    
    linkStatistics, err := linkService.GetLinkStatistics(shortCodeToAnalyze)
    if err != nil {
        log.Fatalf("Erreur lors de la récupération des statistiques: %v", err)
    }
    
    fmt.Printf("Statistiques pour le code court: %s\n", shortCodeToAnalyze)
    fmt.Printf("URL longue: %s\n", linkStatistics.LongURL)
    fmt.Printf("Total de clics: %d\n", linkStatistics.TotalClicks)
}
