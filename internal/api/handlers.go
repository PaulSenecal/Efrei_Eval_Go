package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Channel global des événements de clics (injecté dans RegisterRoutes)
var ClickEventsChannel chan models.ClickEvent

// Enregistrement des routes de l'API
func RegisterRoutes(router *gin.Engine, linkService *services.LinkService, clickChan chan models.ClickEvent) {
	ClickEventsChannel = clickChan

	router.GET("/health", HealthCheckHandler)

	api := router.Group("/api/v1")
	{
		api.POST("/links", CreateShortLinkHandler(linkService))
		api.GET("/links/:shortCode/stats", GetLinkStatsHandler(linkService))
	}

	router.GET("/:shortCode", RedirectHandler(linkService))
}

// Vérification de l’état du service
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Structure de la requête pour POST /links
type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
}

// Création d’un lien court
func CreateShortLinkHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide ou URL manquante."})
			return
		}

		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			log.Printf("Erreur création lien court: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne lors de la création du lien."})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.ShortCode,
			"long_url":       link.LongURL,
			"full_short_url": "http://localhost:8080/" + link.ShortCode, // Remplace par cfg.Server.BaseURL si nécessaire
		})
	}
}

// Redirection avec enregistrement asynchrone du clic
func RedirectHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
				return
			}
			log.Printf("Erreur récupération lien: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne"})
			return
		}

		clickEvent := models.ClickEvent{
			LinkID:    link.ID,
			Timestamp: time.Now(),
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		}

		select {
		case ClickEventsChannel <- clickEvent:
		default:
			log.Printf("Warning: ClickEventsChannel is full, dropping click event for %s.", shortCode)
		}

		c.Redirect(http.StatusFound, link.LongURL)
	}
}

// Récupération des statistiques d’un lien
func GetLinkStatsHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		link, totalClicks, err := linkService.GetLinkStats(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
				return
			}
			log.Printf("Erreur stats: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.ShortCode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
	}
}
