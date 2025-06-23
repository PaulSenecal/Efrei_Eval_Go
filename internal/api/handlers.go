//internal\api\handlers.go
package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Pour gérer gorm.ErrRecordNotFound
)

type APIHandlers struct {
    linkService       services.LinkService
    clickEventChannel chan models.ClickEvent
}

type CreateLinkRequest struct {
    LongURL string `json:"long_url" binding:"required"`
}

type CreateLinkResponse struct {
    ShortCode string `json:"short_code"`
    LongURL   string `json:"long_url"`
}

type LinkStatisticsResponse struct {
    ShortCode   string `json:"short_code"`
    LongURL     string `json:"long_url"`
    TotalClicks int    `json:"total_clicks"`
}

func NewAPIHandlers(linkService services.LinkService, clickEventChannel chan models.ClickEvent) *APIHandlers {
    return &APIHandlers{
        linkService:       linkService,
        clickEventChannel: clickEventChannel,
    }
}

func (handlers *APIHandlers) HandleHealthCheck(ginContext *gin.Context) {
    ginContext.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (handlers *APIHandlers) HandleCreateShortLink(ginContext *gin.Context) {
    var createLinkRequest CreateLinkRequest
    
    if err := ginContext.ShouldBindJSON(&createLinkRequest); err != nil {
        ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Format de requête invalide"})
        return
    }
    
    createdLink, err := handlers.linkService.CreateShortLink(createLinkRequest.LongURL)
    if err != nil {
        ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    createLinkResponse := CreateLinkResponse{
        ShortCode: createdLink.ShortCode,
        LongURL:   createdLink.LongURL,
    }
    
    ginContext.JSON(http.StatusCreated, createLinkResponse)
}

func (handlers *APIHandlers) HandleRedirectToLongURL(ginContext *gin.Context) {
    shortCode := ginContext.Param("shortCode")
    
    link, err := handlers.linkService.GetLinkByShortCode(shortCode)
    if err != nil {
        ginContext.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
        return
    }
    
    clickEvent := models.ClickEvent{
        ShortCode: shortCode,
        UserAgent: ginContext.GetHeader("User-Agent"),
        IPAddress: ginContext.ClientIP(),
        ClickedAt: time.Now(),
    }
    
    select {
    case handlers.clickEventChannel <- clickEvent:
    default:
    }
    
    ginContext.Redirect(http.StatusFound, link.LongURL)
}

func (handlers *APIHandlers) HandleGetLinkStatistics(ginContext *gin.Context) {
    shortCode := ginContext.Param("shortCode")
    
    linkStatistics, err := handlers.linkService.GetLinkStatistics(shortCode)
    if err != nil {
        ginContext.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
        return
    }
    
    statisticsResponse := LinkStatisticsResponse{
        ShortCode:   linkStatistics.ShortCode,
        LongURL:     linkStatistics.LongURL,
        TotalClicks: linkStatistics.TotalClicks,
    }
    
    ginContext.JSON(http.StatusOK, statisticsResponse)
}