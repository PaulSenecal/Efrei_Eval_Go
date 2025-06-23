//cmd\cli\create.go
package cli

import (
	"fmt"
	"log"
	"net/url" // Pour valider le format de l'URL
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

var CreateCommand = &cobra.Command{
    Use:   "create",
    Short: "Crée une URL courte",
    Run:   executeCreateCommand,
}

var longURLToShorten string

func init() {
    CreateCommand.Flags().StringVar(&longURLToShorten, "url", "", "URL longue à raccourcir")
    CreateCommand.MarkFlagRequired("url")
}

func executeCreateCommand(cmd *cobra.Command, args []string) {
    applicationConfiguration := config.LoadApplicationConfiguration()
    
    databaseConnection, err := gorm.Open(sqlite.Open(applicationConfiguration.Database.Path), &gorm.Config{})
    if err != nil {
        log.Fatalf("Impossible de se connecter à la base de données: %v", err)
    }
    
    linkRepository := repository.NewLinkRepository(databaseConnection)
    linkService := services.NewLinkService(linkRepository)
    
    createdLink, err := linkService.CreateShortLink(longURLToShorten)
    if err != nil {
        log.Fatalf("Erreur lors de la création du lien court: %v", err)
    }
    
    fullShortURL := fmt.Sprintf("%s/%s", applicationConfiguration.ShortURL.BaseURL, createdLink.ShortCode)
    
    fmt.Println("URL courte créée avec succès:")
    fmt.Printf("Code: %s\n", createdLink.ShortCode)
    fmt.Printf("URL complète: %s\n", fullShortURL)
}