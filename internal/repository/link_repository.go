package repository

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// Interface pour les opérations sur les liens
type LinkRepository interface {
	CreateLink(link *models.Link) error
	GetLinkByShortCode(shortCode string) (*models.Link, error)
	GetAllLinks() ([]models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}

// Implémentation GORM du LinkRepository
type GormLinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{db: db}
}

// Crée un lien court
func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	if err := r.db.Create(link).Error; err != nil {
		return fmt.Errorf("échec de la création du lien : %w", err)
	}
	return nil
}

// Récupère un lien par son code
func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	if err := r.db.Where("short_code = ?", shortCode).First(&link).Error; err != nil {
		return nil, fmt.Errorf("échec de récupération du lien %s : %w", shortCode, err)
	}
	return &link, nil
}

// Récupère tous les liens (utilisé par le moniteur)
func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	if err := r.db.Find(&links).Error; err != nil {
		return nil, fmt.Errorf("échec de récupération des liens : %w", err)
	}
	return links, nil
}

// Compte les clics associés à un lien
func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	if err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("échec du comptage des clics pour le lien %d : %w", linkID, err)
	}
	return int(count), nil
}