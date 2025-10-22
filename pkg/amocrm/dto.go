package amocrm

import "git.amocrm.ru/study_group/in_memory_database/internal/entity"

type response struct {
	entity.AccountIntegration
	GrantType string `json:"grant_type"`
}
