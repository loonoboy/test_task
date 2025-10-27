package account_integrations

import (
	"log"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IntegrationRepoMySQL struct {
	db *gorm.DB
}

func NewIntegrationRepoMySQL(repo *gorm.DB) *IntegrationRepoMySQL {
	err := repo.AutoMigrate(&entity.AccountIntegration{})
	if err != nil {
		log.Fatal(err)
	}
	return &IntegrationRepoMySQL{db: repo}
}

func (d *IntegrationRepoMySQL) CreateIntegration(a *entity.AccountIntegration) error {
	return d.db.Create(a).Error
}

func (d *IntegrationRepoMySQL) GetIntegration(id uuid.UUID) (*entity.AccountIntegration, error) {
	var inter entity.AccountIntegration
	if err := d.db.Where("client_id = ?", id).
		First(&inter).
		Error; err != nil {
		return nil, err
	}
	return &inter, nil
}

func (d *IntegrationRepoMySQL) ListIntegrations() ([]*entity.AccountIntegration, error) {
	var inters []*entity.AccountIntegration
	if err := d.db.Find(&inters).Error; err != nil {
		return nil, err
	}
	return inters, nil
}

func (d *IntegrationRepoMySQL) UpdateIntegration(id uuid.UUID, update dto.IntegrationUpdate) error {
	return d.db.Model(&entity.AccountIntegration{}).
		Where("client_id = ?", id).
		Updates(update).
		Error
}

func (d *IntegrationRepoMySQL) DeleteIntegration(id uuid.UUID) error {
	return d.db.Where("client_id = ?", id).Delete(entity.AccountIntegration{}).Error
}
