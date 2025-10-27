package amo_client

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"git.amocrm.ru/study_group/in_memory_database/pkg/amocrm"
	"github.com/google/uuid"
)

type AmoClientService struct {
	accRepo     usecase.AccountRepository
	intgrRepo   usecase.IntegrationRepository
	contactRepo usecase.ContactRepository
	client      *amocrm.AMOClient
}

func NewAmoClientServiceService(client *amocrm.AMOClient, repository usecase.AccountRepository) *AmoClientService {
	return &AmoClientService{client: client, accRepo: repository}
}

func (s *AmoClientService) SaveAccountInfo(authCode, subdomain string, clientID uuid.UUID) error {
	integration, err := s.intgrRepo.GetIntegration(clientID)
	if err != nil {
		return err
	}
	integration.Code = authCode

	account, err := s.client.GetTokens(integration, subdomain)
	if err != nil {
		return err
	}
	account.Subdomain = subdomain

	accountID, err := s.client.GetAccountID(account)
	if err != nil {
		return err
	}

	account.AccountID = accountID

	err = s.accRepo.CreateAccount(account)
	if err != nil {
		return err
	}
	integrationUpdates := &dto.IntegrationUpdate{
		AccountID: accountID,
	}

	err = s.intgrRepo.UpdateIntegration(clientID, *integrationUpdates)
	if err != nil {
		return err
	}

	account, err = s.accRepo.GetAccount(accountID)
	if err != nil {
		return err
	}

	contacts, err := s.client.GetContacts(account)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		if err = s.contactRepo.CreateContact(&contact); err != nil {
			return err
		}
	}
	return nil
}
