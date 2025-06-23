//internal\monitor\url_monitor.go
package monitor

import (
	"log"
	"net/http"
	"sync" // Pour protéger l'accès concurrentiel à knownStates
	"time"

	_ "urlshortenerGroupe8/internal/models"   // Importe les modèles de liens
	"urlshortenerGroupe8/internal/repository" // Importe le repository de liens
)
type URLMonitor struct {
    linkService           services.LinkService
    checkIntervalDuration time.Duration
    stopMonitoringChannel chan bool
}

func NewURLMonitor(linkService services.LinkService, checkInterval time.Duration) *URLMonitor {
    return &URLMonitor{
        linkService:           linkService,
        checkIntervalDuration: checkInterval,
        stopMonitoringChannel: make(chan bool),
    }
}

func (monitor *URLMonitor) StartMonitoring() {
    go monitor.runMonitoringLoop()
    log.Printf("Moniteur d'URLs démarré avec un intervalle de %v", monitor.checkIntervalDuration)
}

func (monitor *URLMonitor) StopMonitoring() {
    close(monitor.stopMonitoringChannel)
    log.Println("Moniteur d'URLs arrêté")
}

func (monitor *URLMonitor) runMonitoringLoop() {
    monitoringTicker := time.NewTicker(monitor.checkIntervalDuration)
    defer monitoringTicker.Stop()
    
    for {
        select {
        case <-monitoringTicker.C:
            monitor.checkAllLinksAccessibility()
            
        case <-monitor.stopMonitoringChannel:
            return
        }
    }
}

func (monitor *URLMonitor) checkAllLinksAccessibility() {
    activeLinks, err := monitor.linkService.GetAllActiveLinks()
    if err != nil {
        log.Printf("Erreur lors de la récupération des liens: %v", err)
        return
    }
    
    for _, link := range activeLinks {
        currentAccessibilityStatus := monitor.isURLAccessible(link.LongURL)
        
        if currentAccessibilityStatus != link.IsAccessible {
            monitor.handleAccessibilityStatusChange(link, currentAccessibilityStatus)
        }
    }
}

func (monitor *URLMonitor) isURLAccessible(urlToCheck string) bool {
    httpClient := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    response, err := httpClient.Get(urlToCheck)
    if err != nil {
        return false
    }
    defer response.Body.Close()
    
    return response.StatusCode >= 200 && response.StatusCode < 400
}

func (monitor *URLMonitor) handleAccessibilityStatusChange(link services.Link, newAccessibilityStatus bool) {
    err := monitor.linkService.UpdateLinkAccessibilityStatus(link.ID, newAccessibilityStatus)
    if err != nil {
        log.Printf("Erreur lors de la mise à jour du statut: %v", err)
        return
    }
    
    previousStatus := "ACCESSIBLE"
    currentStatus := "INACCESSIBLE"
    
    if newAccessibilityStatus {
        previousStatus = "INACCESSIBLE"
        currentStatus = "ACCESSIBLE"
    }
    
    notificationMessage := fmt.Sprintf(
        "[NOTIFICATION] Le lien %s (%s) est passé de %s à %s !",
        link.ShortCode,
        link.LongURL,
        previousStatus,
        currentStatus,
    )
    
    log.Println(notificationMessage)
}
