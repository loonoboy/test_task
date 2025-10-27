package accounts

import (
	"fmt"
	"sync"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type AccountsRepository struct {
	mu       sync.Mutex
	accounts map[int]*entity.Account
}

func NewAccountsRepository() *AccountsRepository {
	return &AccountsRepository{
		accounts: make(map[int]*entity.Account),
	}
}

func (repo *AccountsRepository) CreateAccount(account *entity.Account) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.accounts[account.AccountID]; exists {
		return fmt.Errorf("account %v already exists", account.AccountID)
	}
	repo.accounts[account.AccountID] = account
	return nil
}

func (repo *AccountsRepository) GetAccount(id int) (*entity.Account, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	account, ok := repo.accounts[id]
	if !ok {
		return nil, fmt.Errorf("account with id %v not found", id)
	}
	return account, nil
}

func (repo *AccountsRepository) ListAccounts() ([]*entity.Account, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	accounts := make([]*entity.Account, 0, len(repo.accounts))
	for _, val := range repo.accounts {
		accounts = append(accounts, val)
	}
	return accounts, nil
}

func (repo *AccountsRepository) UpdateAccount(id int, update dto.UpdateAccount) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	acc, ok := repo.accounts[id]
	if !ok {
		return fmt.Errorf("account %d not found", id)
	}

	if update.AccessToken != nil {
		acc.AccessToken = *update.AccessToken
	}
	if update.RefreshToken != nil {
		acc.RefreshToken = *update.RefreshToken
	}
	if update.Expires != nil {
		acc.Expires = *update.Expires
	}

	return nil
}

func (repo *AccountsRepository) DeleteAccount(id int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.accounts[id]; !ok {
		return fmt.Errorf("account %v not found", id)
	}
	delete(repo.accounts, id)
	return nil
}
