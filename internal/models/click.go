package models

import "time"

// Click représente un événement de clic sur un lien raccourci.
type Click struct {
	ID          uint      `gorm:"primaryKey"`
	LinkID      uint      `gorm:"not null"`
	IPAddress   string
	UserAgent   string
	ClickedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ClickEvent représente un clic brut, transmis dans un channel (non stocké directement par GORM)
type ClickEvent struct {
	LinkID    uint
	Timestamp time.Time
	UserAgent string
	IP        string
}
