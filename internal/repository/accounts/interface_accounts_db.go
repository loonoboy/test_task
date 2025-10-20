package accounts

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type AccountRepository interface {
	CreateAccount(a *entity.Account) error
	GetAccount(id int) (*entity.Account, error)
	ListAccounts() ([]*entity.Account, error)
	UpdateAccount(id int, update dto.UpdateAccount) error
	DeleteAccount(id int) error
}
