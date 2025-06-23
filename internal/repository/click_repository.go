package repository

import (
"fmt"

"github.com/axellelanca/urlshortener/internal/models"
"gorm.io/gorm"
)

// ClickRepository définit les méthodes d'accès aux clics
type ClickRepository interface {
CreateClick(click *models.Click) error
CountClicksByLinkID(linkID uint) (int, error)
}

// Implémentation avec GORM
type GormClickRepository struct {
db *gorm.DB
}

func NewClickRepository(db *gorm.DB) *GormClickRepository {
return &GormClickRepository{db: db}
}

// Insère un nouveau clic
func (r *GormClickRepository) CreateClick(click *models.Click) error {
if err := r.db.Create(click).Error; err != nil {
return fmt.Errorf("échec de la création du clic : %w", err)
}
return nil
}

// Compte les clics pour un lien donné
func (r *GormClickRepository) CountClicksByLinkID(linkID uint) (int, error) {
var count int64
if err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error; err != nil {
return 0, fmt.Errorf("échec du comptage des clics : %w", err)
}
return int(count), nil
}