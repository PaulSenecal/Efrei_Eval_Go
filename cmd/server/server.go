//cmd\server\server.go
package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "urlshortenerGroupe8/cmd"
	"urlshortenerGroupe8/internal/api"
	"urlshortenerGroupe8/internal/models"
	"urlshortenerGroupe8/internal/monitor"
	"urlshortenerGroupe8/internal/repository"
	"urlshortenerGroupe8/internal/services"
	"urlshortenerGroupe8/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)


var ServerCommand = &cobra.Command{
    Use:   "run-server",
    Short: "Lance le serveur API, les workers et le moniteur",
    Run:   executeServerCommand,
}

func executeServerCommand(cmd *cobra.Command, args []string) {
    applicationConfiguration := config.LoadApplicationConfiguration()
    
    databaseConnection, err := initializeDatabaseConnection(applicationConfiguration.Database.Path)
    if err != nil {
        log.Fatalf("Impossible de se connecter à la base de données: %v", err)
    }

    linkRepository := repository.NewLinkRepository(databaseConnection)
    clickRepository := repository.NewClickRepository(databaseConnection)
    
    linkService := services.NewLinkService(linkRepository)
    clickService := services.NewClickService(clickRepository)
    
    clickEventChannel := make(chan models.ClickEvent, applicationConfiguration.Analytics.BufferSize)
    
    clickWorkerManager := workers.NewClickWorkerManager(
        clickEventChannel,
        clickService,
        applicationConfiguration.Analytics.WorkerCount,
    )
    clickWorkerManager.StartAllWorkers()
    
    urlMonitor := monitor.NewURLMonitor(
        linkService,
        time.Duration(applicationConfiguration.Monitoring.CheckIntervalMinutes)*time.Minute,
    )
    urlMonitor.StartMonitoring()
    
    httpServer := createAndConfigureHTTPServer(
        applicationConfiguration,
        linkService,
        clickEventChannel,
    )
    
    go func() {
        serverAddress := fmt.Sprintf("%s:%d", applicationConfiguration.Server.Host, applicationConfiguration.Server.Port)
        log.Printf("Serveur HTTP démarré sur %s", serverAddress)
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
        }
    }()
    
    waitForShutdownSignal()
    
    shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelShutdown()
    
    if err := httpServer.Shutdown(shutdownContext); err != nil {
        log.Printf("Erreur lors de l'arrêt du serveur: %v", err)
    }
    
    clickWorkerManager.StopAllWorkers()
    urlMonitor.StopMonitoring()
    close(clickEventChannel)
    
    log.Println("Serveur arrêté proprement")
}

func initializeDatabaseConnection(databasePath string) (*gorm.DB, error) {
    return gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
}

func createAndConfigureHTTPServer(
    applicationConfiguration *config.ApplicationConfiguration,
    linkService services.LinkService,
    clickEventChannel chan models.ClickEvent,
) *http.Server {
    ginRouter := gin.Default()
    
    apiHandlers := api.NewAPIHandlers(linkService, clickEventChannel)
    
    ginRouter.GET("/health", apiHandlers.HandleHealthCheck)
    ginRouter.POST("/api/v1/links", apiHandlers.HandleCreateShortLink)
    ginRouter.GET("/:shortCode", apiHandlers.HandleRedirectToLongURL)
    ginRouter.GET("/api/v1/links/:shortCode/stats", apiHandlers.HandleGetLinkStatistics)
    
    return &http.Server{
        Addr:    fmt.Sprintf("%s:%d", applicationConfiguration.Server.Host, applicationConfiguration.Server.Port),
        Handler: ginRouter,
    }
}

func waitForShutdownSignal() {
    shutdownSignalChannel := make(chan os.Signal, 1)
    signal.Notify(shutdownSignalChannel, os.Interrupt, syscall.SIGTERM)
    <-shutdownSignalChannel
    log.Println("Signal d'arrêt reçu")
}
/*
// RunServerCmd représente la commande 'run-server' de Cobra.
// C'est le point d'entrée pour lancer le serveur de l'application.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : Charger la configuration chargée globalement via cmd.cfg
		// Ne pas oublier la gestion d'erreur (si nil ?), si erreur, faire un log.Fataf

		// TODO : Initialiser la connexion à la base de données SQLite avec GORM.
		// Utilisez le nom de la base de données depuis la configuration (cfg.Database.Name).

		// TODO : Initialiser les repositories.
		// Créez des instances de GormLinkRepository et GormClickRepository.

		// Laissez le log
		log.Println("Repositories initialisés.")

		// TODO : Initialiser les services métiers.
		// Créez des instances de LinkService et ClickService, en leur passant les repositories nécessaires.

		// Laissez le log
		log.Println("Services métiers initialisés.")

		// TODO : Initialiser le channel ClickEventsChannel (api/handlers) des événements de clic et lancer les workers (StartClickWorkers).
		// Le channel est bufferisé avec la taille configurée.
		// Passez le channel et le clickRepo aux workers.

		// TODO : Remplacer les XXX par les bonnes variables
		log.Printf("Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
			XXX, XXX)

		// TODO : Initialiser et lancer le moniteur d'URLs.
		// Utilisez l'intervalle configuré (cfg.Monitor.IntervalMinutes).
		// Lancez le moniteur dans sa propre goroutine.
		monitorInterval := time.Duration(XXX) * time.Minute
		urlMonitor := monitor.NewUrlMonitor() // Le moniteur a besoin du linkRepo et de l'interval
		go urlMonitor.Start()
		log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)

		// TODO : Configurer le routeur Gin et les handlers API.
		// Passez les services nécessaires aux fonctions de configuration des routes.

		// Pas toucher au log
		log.Println("Routes API configurées.")

		// Créer le serveur HTTP Gin
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// TODO : Démarrer le serveur Gin dans une goroutine anonyme pour ne pas bloquer.
		// Pensez à logger des ptites informations...

		// Gére l'arrêt propre du serveur (graceful shutdown).
		// Créez un channel pour les signaux OS (SIGINT, SIGTERM).
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Attendre Ctrl+C ou signal d'arrêt

		// Bloquer jusqu'à ce qu'un signal d'arrêt soit reçu.
		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")

		// Arrêt propre du serveur HTTP avec un timeout.
		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	// TODO : ajouter la commande
}
*/
