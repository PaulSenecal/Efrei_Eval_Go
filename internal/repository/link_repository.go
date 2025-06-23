package repository

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/axellelanca/urlshortener/internal/models"
	//"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Struct LinkService
type LinkService struct {
	linkRepo repository.LinkRepository
}

// Constructeur
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// Génère un shortCode sécurisé et aléatoire
func (s *LinkService) GenerateShortCode(length int) (string, error) {
	code := make([]byte, length)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("erreur lors de la génération du shortCode : %w", err)
		}
		code[i] = charset[num.Int64()]
	}
	return string(code), nil
}

// Création d’un lien court avec gestion des collisions
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	var shortCode string
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		code, err := s.GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("échec de la génération du code : %w", err)
		}

		_, err = s.linkRepo.GetLinkByShortCode(code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code
				break
			}
			return nil, fmt.Errorf("erreur lors de la vérification du code : %w", err)
		}

		log.Printf("Short code '%s' déjà existant, tentative %d/%d", code, i+1, maxRetries)
	}

	if shortCode == "" {
		return nil, errors.New("échec : impossible de générer un code unique après plusieurs tentatives")
	}

	link := &models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("échec de la création du lien : %w", err)
	}

	return link, nil
}

// Récupération d’un lien via son code
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, fmt.Errorf("échec de la récupération du lien : %w", err)
	}
	return link, nil
}

// Récupération des statistiques
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("échec de la récupération du lien : %w", err)
	}

	count, err := s.linkRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return link, 0, fmt.Errorf("échec du comptage des clics : %w", err)
	}

	return link, count, nil
}
