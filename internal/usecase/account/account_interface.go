package account

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/google/uuid"
)

type AccountUsecaseInterface interface {
	CreateAccount(authCode, subdomain string, clientID uuid.UUID) error
	GetAccount(id int) (*entity.Account, error)
	ListAccounts() ([]*entity.Account, error)
	UpdateAccount(id int, account dto.UpdateAccount) error
	DeleteAccount(id int) error
}
