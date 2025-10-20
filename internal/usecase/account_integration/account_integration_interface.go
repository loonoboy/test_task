package account_integration

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type AccountIntegrationUsecaseInterface interface {
	CreateIntegration(i entity.AccountIntegration) error
	GetIntegration(id int) (*entity.AccountIntegration, error)
	ListIntegrations() ([]*entity.AccountIntegration, error)
	UpdateIntegration(id int, update dto.IntegrationUpdate) error
	DeleteIntegration(id int) error
}
