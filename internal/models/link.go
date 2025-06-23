//internal\models\link.go
package models
import (
    "time"
    "gorm.io/gorm"
)

import "time"

// Link représente un lien raccourci dans la base de données.
type Link struct {
	ID        uint      `gorm:"primaryKey"`          // Clé primaire
	ShortCode string    `gorm:"size:10;uniqueIndex"` // Code court unique, indexé pour recherche rapide
	LongURL   string    `gorm:"not null"`            // Lien original obligatoire
	CreatedAt time.Time // Géré automatiquement par GORM
}

