package amocrm

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"github.com/google/uuid"
)

type response struct {
	entity.AccountIntegration
	GrantType string `json:"grant_type"`
}

type AmoCloseRequest struct {
	AccountId  int       `json:"account_id"`
	ClientUUID uuid.UUID `json:"client_uuid"`
}
