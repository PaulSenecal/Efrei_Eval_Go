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


var longURLFlag string // Variable globale pour stocker le flag --url

var CreateCmd = &cobra.Command{
Use:   "create",
Short: "Crée une URL courte à partir d'une URL longue.",
Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
Run: func(cmd *cobra.Command, args []string) {
// Vérification que le flag --url est fourni
if longURLFlag == "" {
fmt.Println("Erreur : vous devez fournir une URL avec le flag --url")
os.Exit(1)
}

// Validation du format de l'URL
if _, err := url.ParseRequestURI(longURLFlag); err != nil {
fmt.Printf("Erreur : URL invalide : %v\n", err)
os.Exit(1)
}

// Chargement de la configuration globale
cfg := cmd2.Cfg

// Connexion à la base de données SQLite
db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
if err != nil {
log.Fatalf("FATAL: impossible de se connecter à la base de données: %v", err)
}

sqlDB, err := db.DB()
if err != nil {
log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
}
defer sqlDB.Close() // Fermeture automatique de la base

// Initialisation des repositories et services
linkRepo := repository.NewLinkRepository(db)
linkService := services.NewLinkService(linkRepo)

// Création du lien court
link, err := linkService.CreateLink(longURLFlag)
if err != nil {
log.Printf("Erreur lors de la création du lien : %v", err)
os.Exit(1)
}

fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
fmt.Printf("URL courte créée avec succès:\n")
fmt.Printf("Code: %s\n", link.ShortCode)
fmt.Printf("URL complète: %s\n", fullShortURL)
},
}

func init() {
// Définir et lier le flag --url
CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")

// Rendre le flag requis
err := CreateCmd.MarkFlagRequired("url")
if err != nil {
log.Fatalf("Erreur lors du marquage du flag comme requis: %v", err)
}

// Ajouter la commande à la commande racine
cmd2.RootCmd.AddCommand(CreateCmd)
