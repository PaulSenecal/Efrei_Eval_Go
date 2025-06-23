//internal\repository\click_repository.go
package repository

import (
	"fmt"

	"urlshortenerGroupe8/internal/models"
	"gorm.io/gorm"
)

type ClickRepository interface {
    RecordNewClick(click *models.Click) error
}

type clickRepositoryImplementation struct {
    databaseConnection *gorm.DB
}

func NewClickRepository(databaseConnection *gorm.DB) ClickRepository {
    return &clickRepositoryImplementation{
        databaseConnection: databaseConnection,
    }
}

func (repository *clickRepositoryImplementation) RecordNewClick(click *models.Click) error {
    return repository.databaseConnection.Create(click).Error
}
