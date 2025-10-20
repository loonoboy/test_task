package account_integrations

import (
	"fmt"
	"sync"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type IntegrationsRepository struct {
	mu                  sync.Mutex
	accountIntegrations map[int]*entity.AccountIntegration
}

func NewIntegrationsRepository() *IntegrationsRepository {
	return &IntegrationsRepository{
		accountIntegrations: make(map[int]*entity.AccountIntegration),
	}
}

func (repo *IntegrationsRepository) CreateIntegration(intg *entity.AccountIntegration) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.accountIntegrations[intg.ClientID]; exists {
		return fmt.Errorf("integration %v already exists", intg.ClientID)
	}
	repo.accountIntegrations[intg.ClientID] = intg
	return nil
}

func (repo *IntegrationsRepository) GetIntegration(id int) (*entity.AccountIntegration, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	intg, ok := repo.accountIntegrations[id]
	if !ok {
		return nil, fmt.Errorf("integration with id %v not found", id)
	}
	return intg, nil
}

func (repo *IntegrationsRepository) ListIntegrations() ([]*entity.AccountIntegration, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	integrations := make([]*entity.AccountIntegration, 0, len(repo.accountIntegrations))
	for _, val := range repo.accountIntegrations {
		integrations = append(integrations, val)
	}
	return integrations, nil
}

func (repo *IntegrationsRepository) UpdateIntegration(id int, update dto.IntegrationUpdate) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	intg, ok := repo.accountIntegrations[id]
	if !ok {
		return fmt.Errorf("integration with ClientID %d not found", id)
	}

	if update.SecretKey != nil {
		intg.SecretKey = *update.SecretKey
	}
	if update.RedirectURL != nil {
		intg.RedirectURL = *update.RedirectURL
	}
	if update.AuthCode != nil {
		intg.AuthCode = *update.AuthCode
	}

	return nil
}

func (repo *IntegrationsRepository) DeleteIntegration(id int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.accountIntegrations[id]; !ok {
		return fmt.Errorf("integration %v not found", id)
	}
	delete(repo.accountIntegrations, id)
	return nil
}
