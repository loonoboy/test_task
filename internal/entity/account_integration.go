package entity

import (
	"github.com/google/uuid"
)

type AccountIntegration struct {
	ClientID     uuid.UUID `json:"client_id"`
	AccountID    int       `json:"account_id"`
	ClientSecret string    `json:"client_secret"`
	RedirectURI  string    `json:"redirect_uri"`
	Code         string    `json:"code"`
}
