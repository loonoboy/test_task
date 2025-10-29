package controller

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/google/uuid"
)

type AccountUsecaseInterface interface {
	CreateAccount(i entity.Account) error
	GetAccount(id int) (*entity.Account, error)
	ListAccounts() ([]*entity.Account, error)
	UpdateAccount(id int, account dto.UpdateAccount) error
	DeleteAccount(id int) error
}

type AccountIntegrationUsecaseInterface interface {
	CreateIntegration(i entity.AccountIntegration) error
	GetIntegration(id uuid.UUID) (*entity.AccountIntegration, error)
	ListIntegrations() ([]*entity.AccountIntegration, error)
	UpdateIntegration(id uuid.UUID, update dto.IntegrationUpdate) error
	DeleteIntegration(id uuid.UUID) error
}

type AmoClientUsecaseInterface interface {
	SaveAccountInfo(authCode, subdomain string, clientID uuid.UUID) error
}

type ContactsUsecaseInterface interface {
	CreateContact(i entity.Contact) error
	GetContact(id int) (*entity.Contact, error)
	ListContacts(id int) ([]*entity.Contact, error)
	UpdateContact(id int, contact dto.UpdateContact) error
	DeleteContact(id int) error
}

type UnisenderUsecaseInterface interface {
	SaveUnisenderKey(id int, update dto.UpdateAccount) error
	SaveExistingContacts(accountID int) error
}
