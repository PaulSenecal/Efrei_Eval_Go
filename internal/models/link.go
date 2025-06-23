//internal\models\link.go
package models
import (
    "time"
)

type Link struct {
	ID        uint      `gorm:"primaryKey"`
	ShortCode string    `gorm:"unique;index;size:10"`
	LongURL   string    `gorm:"not null"`
	CreatedAt time.Time
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
