//internal\services\link_service.go
package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"urlshortenerGroupe8/internal/models"
	"urlshortenerGroupe8/internal/repository" // Importe le package repository
)

type LinkService interface {
    CreateShortLink(longURL string) (*models.Link, error)
    GetLinkByShortCode(shortCode string) (*models.Link, error)
    GetLinkStatistics(shortCode string) (*models.LinkStatistics, error)
    GetAllActiveLinks() ([]models.Link, error)
    UpdateLinkAccessibilityStatus(linkID uint, isAccessible bool) error
}

type linkServiceImplementation struct {
    linkRepository repository.LinkRepository
}

func NewLinkService(linkRepository repository.LinkRepository) LinkService {
    return &linkServiceImplementation{
        linkRepository: linkRepository,
    }
}

func (service *linkServiceImplementation) CreateShortLink(longURL string) (*models.Link, error) {
    if longURL == "" {
        return nil, errors.New("l'URL ne peut pas être vide")
    }
    
    uniqueShortCode, err := service.generateUniqueShortCode()
    if err != nil {
        return nil, err
    }
    
    newLink := &models.Link{
        ShortCode:    uniqueShortCode,
        LongURL:      longURL,
        IsAccessible: true,
    }
    
    err = service.linkRepository.CreateNewLink(newLink)
    if err != nil {
        return nil, err
    }
    
    return newLink, nil
}

func (service *linkServiceImplementation) GetLinkByShortCode(shortCode string) (*models.Link, error) {
    return service.linkRepository.GetLinkByShortCode(shortCode)
}

func (service *linkServiceImplementation) GetLinkStatistics(shortCode string) (*models.LinkStatistics, error) {
    link, err := service.linkRepository.GetLinkByShortCode(shortCode)
    if err != nil {
        return nil, err
    }
    
    return &models.LinkStatistics{
        ShortCode:   link.ShortCode,
        LongURL:     link.LongURL,
        TotalClicks: link.TotalClicks,
    }, nil
}

func (service *linkServiceImplementation) GetAllActiveLinks() ([]models.Link, error) {
    return service.linkRepository.GetAllActiveLinks()
}

func (service *linkServiceImplementation) UpdateLinkAccessibilityStatus(linkID uint, isAccessible bool) error {
    return service.linkRepository.UpdateLinkAccessibilityStatus(linkID, isAccessible)
}

func (service *linkServiceImplementation) generateUniqueShortCode() (string, error) {
    const maximumRetryAttempts = 10
    
    for attemptNumber := 0; attemptNumber < maximumRetryAttempts; attemptNumber++ {
        randomShortCode, err := generateRandomAlphanumericCode(6)
        if err != nil {
            return "", err
        }
        
        if !service.linkRepository.DoesShortCodeExist(randomShortCode) {
            return randomShortCode, nil
        }
    }
    
    return "", errors.New("impossible de générer un code unique après plusieurs tentatives")
}

func generateRandomAlphanumericCode(codeLength int) (string, error) {
    randomBytes := make([]byte, codeLength)
    _, err := rand.Read(randomBytes)
    if err != nil {
        return "", err
    }
    
    encodedCode := base64.URLEncoding.EncodeToString(randomBytes)
    return encodedCode[:codeLength], nil
}
