package account_integration

import (
	"fmt"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/account_integrations"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type AccountInegrationUsecase struct {
	repo account_integrations.IntegrationRepository
}

func NewAccountInegrationUsecase(repo account_integrations.IntegrationRepository) *AccountInegrationUsecase {
	return &AccountInegrationUsecase{repo: repo}
}

func (s *AccountInegrationUsecase) validateIntegration(i entity.AccountIntegration) error {
	if i.AuthCode == "" {
		return fmt.Errorf("auth_code is empty")
	}
	if i.SecretKey == "" {
		return fmt.Errorf("secret_key is empty")
	}
	if i.RedirectURL == "" {
		return fmt.Errorf("redirec_url is empty")
	}
	return nil
}

func (s *AccountInegrationUsecase) CreateIntegration(i entity.AccountIntegration) error {
	if err := s.validateIntegration(i); err != nil {
		return err
	}
	if _, err := s.repo.GetIntegration(i.ClientID); err == nil {
		return fmt.Errorf("account with ID %v already exists", i.ClientID)
	}
	return s.repo.CreateIntegration(&i)
}

func (s *AccountInegrationUsecase) GetIntegration(id int) (*entity.AccountIntegration, error) {
	return s.repo.GetIntegration(id)
}

func (s *AccountInegrationUsecase) ListIntegrations() ([]*entity.AccountIntegration, error) {
	return s.repo.ListIntegrations()
}

func (s *AccountInegrationUsecase) UpdateIntegration(id int, update dto.IntegrationUpdate) error {
	return s.repo.UpdateIntegration(id, update)
}

func (s *AccountInegrationUsecase) DeleteIntegration(id int) error {
	return s.repo.DeleteIntegration(id)
}
