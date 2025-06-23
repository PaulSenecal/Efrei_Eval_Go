package services

import (
"fmt"

"github.com/axellelanca/urlshortener/internal/models"
"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// ClickService gère la logique métier liée aux clics
type ClickService struct {
clickRepo repository.ClickRepository
}

// Constructeur
func NewClickService(clickRepo repository.ClickRepository) *ClickService {
return &ClickService{
clickRepo: clickRepo,
}
}

// Enregistre un clic
func (s *ClickService) RecordClick(click *models.Click) error {
if err := s.clickRepo.CreateClick(click); err != nil {
return fmt.Errorf("erreur lors de l'enregistrement du clic : %w", err)
}
return nil
}

// Compte les clics pour un lien donné
func (s *ClickService) GetClicksCountByLinkID(linkID uint) (int, error) {
count, err := s.clickRepo.CountClicksByLinkID(linkID)
if err != nil {
return 0, fmt.Errorf("erreur lors du comptage des clics : %w", err)
}
return count, nil
}