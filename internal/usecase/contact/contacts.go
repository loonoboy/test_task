package contact

import (
	"fmt"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type ContactsUsecase struct {
	repo usecase.ContactRepository
}

func NewContactUsecase(repo usecase.ContactRepository) *ContactsUsecase {
	return &ContactsUsecase{repo: repo}
}

func (s *ContactsUsecase) validateContact(contact entity.Contact) error {
	if contact.Email == "" {
		return fmt.Errorf("email is empty")
	}
	if contact.Name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}

func (s *ContactsUsecase) CreateContact(i entity.Contact) error {
	if err := s.validateContact(i); err != nil {
		return err
	}
	if _, err := s.repo.GetContact(i.AccountID); err == nil {
		return fmt.Errorf("contact with ID %v already exists", i.AccountID)
	}
	return s.repo.CreateContact(&i)
}

func (s *ContactsUsecase) GetContact(id int) (*entity.Contact, error) {
	account, err := s.repo.GetContact(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *ContactsUsecase) ListContacts(id int) ([]*entity.Contact, error) {
	contacts, err := s.repo.ListContacts(id)
	if err != nil {
		return nil, err
	}
	var result []*entity.Contact
	for _, contact := range contacts {
		if s.validateContact(*contact) != nil {
			result = append(result, contact)
		}
	}
	return result, nil
}

func (s *ContactsUsecase) UpdateContact(id int, update dto.UpdateContact) error {
	err := s.repo.UpdateContact(id, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *ContactsUsecase) DeleteContact(id int) error {
	err := s.repo.DeleteContact(id)
	if err != nil {
		return err
	}
	return nil
}
