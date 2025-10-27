package accounts

import (
	"log"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"gorm.io/gorm"
)

type AccountRepoMySQL struct {
	db *gorm.DB
}

func NewAccountRepoMySQL(repo *gorm.DB) *AccountRepoMySQL {
	err := repo.AutoMigrate(&entity.Account{})
	if err != nil {
		log.Fatal(err)
	}
	return &AccountRepoMySQL{db: repo}
}

func (d *AccountRepoMySQL) CreateAccount(a *entity.Account) error {
	return d.db.Create(a).Error
}

func (d *AccountRepoMySQL) GetAccount(id int) (*entity.Account, error) {
	var account entity.Account
	if err := d.db.Where("account_id = ?", id).
		First(&account).
		Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (d *AccountRepoMySQL) ListAccounts() ([]*entity.Account, error) {
	var accounts []*entity.Account
	if err := d.db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (d *AccountRepoMySQL) UpdateAccount(id int, update dto.UpdateAccount) error {
	return d.db.Model(&entity.Account{}).
		Where("account_id = ?", id).
		Updates(update).
		Error
}

func (d *AccountRepoMySQL) DeleteAccount(id int) error {
	return d.db.Where("account_id = ?", id).Delete(entity.Account{}).Error
}
