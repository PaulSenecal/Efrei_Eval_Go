//internal\repository\link_repository.go
package repository

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

type LinkRepository interface {
    CreateNewLink(link *models.Link) error
    GetLinkByShortCode(shortCode string) (*models.Link, error)
    DoesShortCodeExist(shortCode string) bool
    IncrementClickCount(shortCode string) error
    GetAllActiveLinks() ([]models.Link, error)
    UpdateLinkAccessibilityStatus(linkID uint, isAccessible bool) error
}

type linkRepositoryImplementation struct {
    databaseConnection *gorm.DB
}

func NewLinkRepository(databaseConnection *gorm.DB) LinkRepository {
    return &linkRepositoryImplementation{
        databaseConnection: databaseConnection,
    }
}

func (repository *linkRepositoryImplementation) CreateNewLink(link *models.Link) error {
    return repository.databaseConnection.Create(link).Error
}

func (repository *linkRepositoryImplementation) GetLinkByShortCode(shortCode string) (*models.Link, error) {
    var link models.Link
    result := repository.databaseConnection.Where("short_code = ?", shortCode).First(&link)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, errors.New("lien non trouvé")
        }
        return nil, result.Error
    }
    return &link, nil
}

func (repository *linkRepositoryImplementation) DoesShortCodeExist(shortCode string) bool {
    var count int64
    repository.databaseConnection.Model(&models.Link{}).Where("short_code = ?", shortCode).Count(&count)
    return count > 0
}

func (repository *linkRepositoryImplementation) IncrementClickCount(shortCode string) error {
    return repository.databaseConnection.Model(&models.Link{}).
        Where("short_code = ?", shortCode).
        Update("total_clicks", gorm.Expr("total_clicks + ?", 1)).Error
}

func (repository *linkRepositoryImplementation) GetAllActiveLinks() ([]models.Link, error) {
    var activeLinks []models.Link
    result := repository.databaseConnection.Find(&activeLinks)
    return activeLinks, result.Error
}

func (repository *linkRepositoryImplementation) UpdateLinkAccessibilityStatus(linkID uint, isAccessible bool) error {
    return repository.databaseConnection.Model(&models.Link{}).
        Where("id = ?", linkID).
        Update("is_accessible", isAccessible).Error
}