package account

import (
	"fmt"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"git.amocrm.ru/study_group/in_memory_database/pkg/amocrm"
	"github.com/google/uuid"
)

type AccountUsecase struct {
	accRepo   usecase.AccountRepository
	intgrRepo usecase.IntegrationRepository
	client    *amocrm.AMOClient
}

func NewAccountUsecase(accRepo usecase.AccountRepository,
	intgrRepo usecase.IntegrationRepository, client *amocrm.AMOClient) *AccountUsecase {
	return &AccountUsecase{accRepo: accRepo,
		intgrRepo: intgrRepo,
		client:    client}
}

func (s *AccountUsecase) validateAccount(account entity.Account) error {
	if account.AccessToken == "" {
		return fmt.Errorf("access_token is empty")
	}
	if account.RefreshToken == "" {
		return fmt.Errorf("refresh_token is empty")
	}
	if account.Expires == 0 {
		return fmt.Errorf("expires is empty")
	}
	return nil
}

func (s *AccountUsecase) CreateAccount(authCode, subdomain string, clientID uuid.UUID) error {
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
	return nil
}

func (s *AccountUsecase) GetAccount(id int) (*entity.Account, error) {
	account, err := s.accRepo.GetAccount(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountUsecase) ListAccounts() ([]*entity.Account, error) {
	accounts, err := s.accRepo.ListAccounts()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *AccountUsecase) UpdateAccount(id int, update dto.UpdateAccount) error {
	err := s.accRepo.UpdateAccount(id, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountUsecase) DeleteAccount(id int) error {
	err := s.accRepo.DeleteAccount(id)
	if err != nil {
		return err
	}
	return nil
}
