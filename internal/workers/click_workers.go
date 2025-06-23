//internal\workers\click_workers.go
package workers

import (
	"log"

	"urlshortenerGroupe8/internal/models"
	"urlshortenerGroupe8/internal/repository" // Nécessaire pour interagir avec le ClickRepository
)

type ClickWorkerManager struct {
    clickEventChannel chan models.ClickEvent
    clickService      services.ClickService
    numberOfWorkers   int
    waitGroup         sync.WaitGroup
    shouldStop        chan bool
}

func NewClickWorkerManager(
    clickEventChannel chan models.ClickEvent,
    clickService services.ClickService,
    numberOfWorkers int,
) *ClickWorkerManager {
    return &ClickWorkerManager{
        clickEventChannel: clickEventChannel,
        clickService:      clickService,
        numberOfWorkers:   numberOfWorkers,
        shouldStop:        make(chan bool),
    }
}

func (manager *ClickWorkerManager) StartAllWorkers() {
    for workerID := 1; workerID <= manager.numberOfWorkers; workerID++ {
        manager.waitGroup.Add(1)
        go manager.runSingleWorker(workerID)
    }
    log.Printf("Démarrage de %d workers pour l'enregistrement des clics", manager.numberOfWorkers)
}

func (manager *ClickWorkerManager) StopAllWorkers() {
    close(manager.shouldStop)
    manager.waitGroup.Wait()
    log.Println("Tous les workers de clics ont été arrêtés")
}

func (manager *ClickWorkerManager) runSingleWorker(workerID int) {
    defer manager.waitGroup.Done()
    
    for {
        select {
        case clickEvent, channelIsOpen := <-manager.clickEventChannel:
            if !channelIsOpen {
                return
            }
            manager.processClickEvent(clickEvent, workerID)
            
        case <-manager.shouldStop:
            return
        }
    }
}

func (manager *ClickWorkerManager) processClickEvent(clickEvent models.ClickEvent, workerID int) {
    log.Printf("Worker %d traite un clic pour le code: %s", workerID, clickEvent.ShortCode)
    
    err := manager.clickService.RecordClickEvent(clickEvent, 0)
    if err != nil {
        log.Printf("Worker %d - Erreur lors de l'enregistrement du clic: %v", workerID, err)
    }
}
