package contact

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/pkg/amocrm"
)

type ContactsService struct {
	accountsRepository usecase.AccountRepository
	client             *amocrm.AMOClient
}

func NewContactsService(client *amocrm.AMOClient, repository usecase.AccountRepository) *ContactsService {
	return &ContactsService{client: client, accountsRepository: repository}
}

func (s *ContactsService) GetAllContacts(accountID int) ([]entity.Contact, error) {
	account, err := s.accountsRepository.GetAccount(accountID)
	if err != nil {
		return nil, err
	}

	contacts, err := s.client.GetContacts(account)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}
