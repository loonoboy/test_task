package account

import (
	"fmt"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type AccountUsecase struct {
	repo usecase.AccountRepository
}

func NewAccountUsecase(repo usecase.AccountRepository) *AccountUsecase {
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

func (s *AccountUsecase) CreateAccount(i entity.Account) error {
	if err := s.validateAccount(i); err != nil {
		return err
	}
	if _, err := s.repo.GetAccount(i.AccountID); err == nil {
		return fmt.Errorf("account with ID %v already exists", i.AccountID)
	}
	return s.repo.CreateAccount(&i)
}

func (s *AccountUsecase) GetAccount(id int) (*entity.Account, error) {
	account, err := s.repo.GetAccount(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountUsecase) ListAccounts() ([]*entity.Account, error) {
	accounts, err := s.repo.ListAccounts()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *AccountUsecase) UpdateAccount(id int, update dto.UpdateAccount) error {
	err := s.repo.UpdateAccount(id, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountUsecase) DeleteAccount(id int) error {
	err := s.repo.DeleteAccount(id)
	if err != nil {
		return err
	}
	return nil
}
