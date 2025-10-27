package usecase

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/google/uuid"
)

type AccountRepository interface {
	CreateAccount(a *entity.Account) error
	GetAccount(id int) (*entity.Account, error)
	ListAccounts() ([]*entity.Account, error)
	UpdateAccount(id int, update dto.UpdateAccount) error
	DeleteAccount(id int) error
}

type IntegrationRepository interface {
	CreateIntegration(i *entity.AccountIntegration) error
	GetIntegration(id uuid.UUID) (*entity.AccountIntegration, error)
	ListIntegrations() ([]*entity.AccountIntegration, error)
	UpdateIntegration(id uuid.UUID, update dto.IntegrationUpdate) error
	DeleteIntegration(id uuid.UUID) error
}
