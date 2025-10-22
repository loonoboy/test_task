package account_integration

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/google/uuid"
)

type AccountIntegrationUsecaseInterface interface {
	CreateIntegration(i entity.AccountIntegration) error
	GetIntegration(id uuid.UUID) (*entity.AccountIntegration, error)
	ListIntegrations() ([]*entity.AccountIntegration, error)
	UpdateIntegration(id uuid.UUID, update dto.IntegrationUpdate) error
	DeleteIntegration(id uuid.UUID) error
}
