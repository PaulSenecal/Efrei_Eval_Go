package cmd

import (
	"fmt"
	"log"
	"os"

	"urlshortenerGroupe8/cmd/cli"
	"urlshortenerGroupe8/cmd/server"
	"urlshortenerGroupe8/internal/config"
	"github.com/spf13/cobra"
)

var Cfg *config.Config

var RootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long: `'url-shortener' est une application complète pour gérer des URLs courtes.
Elle inclut un serveur API pour le raccourcissement et la redirection,
ainsi qu'une interface en ligne de commande pour l'administration.

Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de la commande: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	
	RootCmd.AddCommand(server.RunServerCmd)
	RootCmd.AddCommand(cli.CreateCmd)
	RootCmd.AddCommand(cli.StatsCmd)
	RootCmd.AddCommand(cli.MigrateCmd)
}

func initConfig() {
	var err error
	Cfg, err = config.LoadConfig()
	if err != nil {
		log.Printf("Attention: Problème lors du chargement de la configuration: %v. Utilisation des valeurs par défaut.", err)
	}
}