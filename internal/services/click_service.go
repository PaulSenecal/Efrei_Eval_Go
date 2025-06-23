
package services

import (
    "time"
	"urlshortenerGroupe8/internal/models"
	"urlshortenerGroupe8/internal/repository" // Importe le package repository
)

type ClickService interface {
    RecordClickEvent(clickEvent models.ClickEvent, linkID uint) error
}

type clickServiceImplementation struct {
    clickRepository repository.ClickRepository
    linkRepository  repository.LinkRepository
}

func NewClickService(clickRepository repository.ClickRepository) ClickService {
    return &clickServiceImplementation{
        clickRepository: clickRepository,
    }
}

func (service *clickServiceImplementation) RecordClickEvent(clickEvent models.ClickEvent, linkID uint) error {
    clickRecord := &models.Click{
        LinkID:    linkID,
        IPAddress: clickEvent.IPAddress,
        UserAgent: clickEvent.UserAgent,
        ClickedAt: time.Now(),
    }
    
    return service.clickRepository.RecordNewClick(clickRecord)
}
