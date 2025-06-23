//cmd\root.go
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/axellelanca/urlshortener/internal/config"
	"github.com/spf13/cobra"
)

// cfg est la variable globale qui contiendra la configuration chargée.
// Elle sera accessible à toutes les commandes Cobra.
var Cfg *config.Config

// TODO : Créer la RootCmd avec Cobra
// Utiliser ces descriptions :
// "Un service de raccourcissement d'URLs avec API REST et CLI"
// `
//'url-shortener' est une application complète pour gérer des URLs courtes.
//Elle inclut un serveur API pour le raccourcissement et la redirection,
//ainsi qu'une interface en ligne de commande pour l'administration.
//
//Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`

// rootCmd représente la commande de base lorsque l'on appelle l'application sans sous-commande.
// C'est le point d'entrée principal pour Cobra.

// Execute est le point d'entrée principal pour l'application Cobra.
// Il est appelé depuis 'main.go'.


var rootCommand = &cobra.Command{
    Use:   "url-shortener",
    Short: "Service de raccourcissement d'URLs",
}

func Execute() {
    if err := rootCommand.Execute(); err != nil {
        //fmt.Println(err)
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de la commande: %v\n", err)
        os.Exit(1)
    }
}

func init() {
    rootCommand.AddCommand(server.ServerCommand)
    rootCommand.AddCommand(cli.CreateCommand)
    rootCommand.AddCommand(cli.StatsCommand)
    rootCommand.AddCommand(cli.MigrateCommand)
}

func initConfig() {
	var err error
	Cfg, err = config.LoadConfig()
	if err != nil {
		log.Printf("Attention: Problème lors du chargement de la configuration: %v. Utilisation des valeurs par défaut.", err)
	}
}
