package contacts

import (
	"log"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"gorm.io/gorm"
)

type ContactRepoMySQL struct {
	db *gorm.DB
}

func NewContactRepoMySQL(repo *gorm.DB) *ContactRepoMySQL {
	err := repo.AutoMigrate(&entity.Contact{})
	if err != nil {
		log.Fatal(err)
	}
	return &ContactRepoMySQL{db: repo}
}

func (d *ContactRepoMySQL) CreateContact(a *entity.Contact) error {
	return d.db.Create(a).Error
}

func (d *ContactRepoMySQL) GetContact(id int) (*entity.Contact, error) {
	var contact entity.Contact
	if err := d.db.Where("contact_id = ?", id).
		First(&contact).
		Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (d *ContactRepoMySQL) ListContacts(id int) ([]*entity.Contact, error) {
	var contacts []*entity.Contact
	if err := d.db.Where("account_id = ?", id).
		Find(&contacts).
		Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (d *ContactRepoMySQL) UpdateContact(id int, update dto.UpdateContact) error {
	return d.db.Model(&entity.Contact{}).
		Where("contact_id = ?", id).
		Updates(update).
		Error
}

func (d *ContactRepoMySQL) DeleteContact(id int) error {
	return d.db.Where("contact_id = ?", id).Delete(entity.Contact{}).Error
}
