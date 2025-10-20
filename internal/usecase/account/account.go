package account

import (
	"fmt"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/accounts"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type AccountUsecase struct {
	repo accounts.AccountRepository
}

func NewAccountUsecase(repo accounts.AccountRepository) *AccountUsecase {
	return &AccountUsecase{repo: repo}
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

func (s *AccountUsecase) CreateAccount(account entity.Account) error {
	if err := s.validateAccount(account); err != nil {
		return err
	}
	if _, err := s.repo.GetAccount(account.AccountID); err == nil {
		return fmt.Errorf("account with ID %v already exists", account.AccountID)
	}
	return s.repo.CreateAccount(&account)
}

func (s *AccountUsecase) GetAccount(id int) (*entity.Account, error) {
	return s.repo.GetAccount(id)
}

func (s *AccountUsecase) ListAccounts() ([]*entity.Account, error) {
	return s.repo.ListAccounts()
}

func (s *AccountUsecase) UpdateAccount(id int, update dto.UpdateAccount) error {
	return s.repo.UpdateAccount(id, update)
}

func (s *AccountUsecase) DeleteAccount(id int) error {
	return s.repo.DeleteAccount(id)
}
