package workers

import (
"log"

"github.com/axellelanca/urlshortener/internal/models"
"github.com/axellelanca/urlshortener/internal/repository" // Nécessaire pour interagir avec le ClickRepository
)

// Démarre les workers
func StartClickWorkers(workerCount int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
log.Printf("Starting %d click worker(s)...", workerCount)
for i := 0; i < workerCount; i++ {
go clickWorker(clickEventsChan, clickRepo)
}
}

// Worker unique : lit les événements et les enregistre
func clickWorker(clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
for event := range clickEventsChan {
click := &models.Click{
LinkID:    event.LinkID,
Timestamp: event.Timestamp,
UserAgent: event.UserAgent,
IPAddress: event.IP,
}

err := clickRepo.CreateClick(click)
if err != nil {
log.Printf("ERROR: Failed to save click for LinkID %d (UserAgent: %s, IP: %s): %v",
event.LinkID, event.UserAgent, event.IP, err)
} else {
log.Printf("Click recorded successfully for LinkID %d", event.LinkID)
}
}
}