package contact

import (
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
)

type ContactUsecaseInterface interface {
	GetAllContacts(accountID int) ([]entity.Contact, error)
}
