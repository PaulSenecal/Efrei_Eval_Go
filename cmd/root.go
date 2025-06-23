package cmd

import (
"fmt"
"log"
"os"

"github.com/axellelanca/urlshortener/internal/config"
"github.com/spf13/cobra"
)

// Cfg est la variable globale contenant la configuration
var Cfg *config.Config

// Déclaration de la commande racine
var RootCmd = &cobra.Command{
Use:   "url-shortener",
Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
Long: `'url-shortener' est une application complète pour gérer des URLs courtes.
Elle inclut un serveur API pour le raccourcissement et la redirection,
ainsi qu'une interface en ligne de commande pour l'administration.

Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`,
}

// Point d’entrée principal de l’application
func Execute() {
if err := RootCmd.Execute(); err != nil {
fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de la commande: %v\n", err)
os.Exit(1)
}
}

func init() {
// Initialise la config au démarrage de Cobra
cobra.OnInitialize(initConfig)
// Les sous-commandes seront enregistrées via leur propre init()
}

// Chargement de la configuration Viper (appelé automatiquement par Cobra)
func initConfig() {
var err error
Cfg, err = config.LoadConfig()
if err != nil {
log.Printf("Attention: Problème lors du chargement de la configuration: %v. Utilisation des valeurs par défaut.", err)
}
}