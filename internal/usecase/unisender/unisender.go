package unisender

import (
	"fmt"
	"log"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/provider"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type UnisenderService struct {
	AccountRepo usecase.AccountRepository
	ContactRepo usecase.ContactRepository
	Provider    *provider.UnisenderProvider
}

func NewUnisenderService(accountRepo usecase.AccountRepository, contactRepo usecase.ContactRepository, prov *provider.UnisenderProvider) *UnisenderService {
	return &UnisenderService{
		AccountRepo: accountRepo,
		ContactRepo: contactRepo,
		Provider:    prov,
	}
}

func (s *UnisenderService) SaveExistingContacts(accountID int) error {
	account, err := s.AccountRepo.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("get account %d: %w", accountID, err)
	}

	contacts, err := s.ContactRepo.ListContacts(account.AccountID)
	if err != nil {
		return fmt.Errorf("list contacts for account %d: %w", account.AccountID, err)
	}

	amoListID, err := s.Provider.CreateOrGetList(account.UnisenderKey, "amoCRM")
	if err != nil {
		return fmt.Errorf("create or get UniSender list: %w", err)
	}

	if err := s.Provider.ImportContacts(account.UnisenderKey, amoListID, contacts); err != nil {
		return fmt.Errorf("import contacts: %w", err)
	}

	if err := s.markContactsSynced(account.AccountID, contacts); err != nil {
		return fmt.Errorf("mark contacts synced: %w", err)
	}

	return nil
}

func (s *UnisenderService) markContactsSynced(accountID int, contacts []*entity.Contact) error {
	for _, c := range contacts {
		upd := dto.UpdateContact{
			AccountID: &accountID,
			Email:     &c.Email,
			IsSynced:  true,
		}
		if err := s.ContactRepo.UpdateContact(c.ContactID, upd); err != nil {
			log.Printf("[Unisender] failed to update contact %d: %v", c.ContactID, err)
			return fmt.Errorf("update contact %d: %w", c.ContactID, err)
		}
	}
	return nil
}

func (u *UnisenderService) SaveUnisenderKey(id int, update dto.UpdateAccount) error {
	if err := u.AccountRepo.UpdateAccount(id, update); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *UnisenderService) MakeSyncContacts(id int) error {
	account, err := s.AccountRepo.GetAccount(id)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	apiKey := account.UnisenderKey
	if apiKey == "" {
		return fmt.Errorf("account %d has no UniSender API key", id)
	}

	contacts, err := s.ContactRepo.ListNotSyncedContacts(account.AccountID)
	if err != nil {
		return fmt.Errorf("list not synced contacts: %w", err)
	}
	if len(contacts) == 0 {
	}

	listName := fmt.Sprintf("amo_account_%d", account.AccountID)
	listID, err := s.Provider.CreateOrGetList(apiKey, listName)
	if err != nil {
		return fmt.Errorf("create or get UniSender list: %w", err)
	}

	if err := s.Provider.ImportContacts(apiKey, listID, contacts); err != nil {
		return fmt.Errorf("import contacts to UniSender: %w", err)
	}

	if err := s.markContactsSynced(account.AccountID, contacts); err != nil {
		return fmt.Errorf("mark contacts synced: %w", err)
	}

	return nil
}

func (s *UnisenderService) DeleteContact(email string, id int) error {
	account, err := s.AccountRepo.GetAccount(id)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	apiKey := account.UnisenderKey
	if apiKey == "" {
		return fmt.Errorf("account %d has no UniSender API key", id)
	}

	listName := fmt.Sprintf("amo_account_%d", account.AccountID)
	listID, err := s.Provider.CreateOrGetList(apiKey, listName)
	if err != nil {
		return fmt.Errorf("create or get UniSender list: %w", err)
	}

	if err := s.Provider.ExcludeContact(apiKey, email, listID); err != nil {
		return fmt.Errorf("exclude contact: %w", err)
	}
	return nil
}
