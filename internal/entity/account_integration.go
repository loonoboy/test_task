package entity

import (
	"github.com/google/uuid"
)

type AccountIntegration struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	ClientID     uuid.UUID `gorm:"type:char(36);not null" json:"client_id"`
	AccountID    int       `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"account_id"`
	ClientSecret string    `gorm:"type:text;not null" json:"client_secret"`
	RedirectURI  string    `gorm:"type:text;not null" json:"redirect_uri"`
	Code         string    `gorm:"type:text" json:"code"`
}
