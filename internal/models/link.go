//internal\models\link.go
package models
import (
    "time"
    "gorm.io/gorm"
)

type Link struct {
    ID              uint      `gorm:"primaryKey"`
    ShortCode       string    `gorm:"uniqueIndex;not null"`
    LongURL         string    `gorm:"not null"`
    TotalClicks     int       `gorm:"default:0"`
    IsAccessible    bool      `gorm:"default:true"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    LastCheckedAt   *time.Time
}

type LinkStatistics struct {
    ShortCode   string
    LongURL     string
    TotalClicks int
}

type ClickEvent struct {
    ShortCode    string
    UserAgent    string
    IPAddress    string
    ClickedAt    time.Time
}