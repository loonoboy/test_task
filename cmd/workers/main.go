package workers

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/provider"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/accounts"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/contacts"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/init_mysql"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/contact"

	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/unisender"
	"github.com/beanstalkd/go-beanstalk"
)

func main() {
	reqFlag := flag.String("tube", "", "Tube to use")
	flag.Parse()
	tubeName := *reqFlag
	if tubeName == "" {
		log.Fatal("tube flag is required")
	}

	db := init_mysql.NewConnectMySQL()

	accountsRepo := accounts.NewAccountRepoMySQL(db)
	contactsRepo := contacts.NewContactRepoMySQL(db)
	unisenderProvider := provider.NewUnisenderProvider()

	contactService := contact.NewContactUsecase(contactsRepo)
	unisenderService := unisender.NewUnisenderService(accountsRepo, contactsRepo, unisenderProvider)

	conn, err := beanstalk.Dial("tcp", "ddev-beanstalkd:11300")
	if err != nil {
		log.Fatalf("error dialing beanstalk: %v", err)
	}
	defer conn.Close()

	tube := beanstalk.NewTubeSet(conn, tubeName)
	for {
		id, body, err := tube.Reserve(5 * time.Second)
		if err != nil {
			if errors.Is(err, beanstalk.ErrTimeout) {
				continue
			}
			log.Printf("[error] beanstalk reserve: %v", err)
			continue
		}

		switch tubeName {
		case "create_contact":
			err = handleCreate(contactService, unisenderService, body)
		case "update_contact":
			err = handleUpdate(contactService, unisenderService, body)
		case "delete_contact":
			err = handleDelete(contactService, unisenderService, body)
		default:
			log.Printf("unknown tube: %s", tubeName)
			conn.Delete(id)
			continue
		}

		if err != nil {
			log.Printf("job %d failed: %v", id, err)
			conn.Release(id, 1, 10*time.Second)
		} else {
			conn.Delete(id)
		}
	}
}

func handleCreate(contactService *contact.ContactsUsecase, unisenderService *unisender.UnisenderService,
	body []byte) error {
	log.Println("Creating contact...")
	var req entity.Contact
	if err := json.Unmarshal(body, &req); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	if req.Email == "" {
		return errors.New("email is empty")
	}
	if err := contactService.CreateContact(req); err != nil {
		return fmt.Errorf("create contact: %w", err)
	}
	return unisenderService.MakeSyncContacts(req.AccountID)
}

func handleUpdate(contactService *contact.ContactsUsecase, unisenderService *unisender.UnisenderService,
	body []byte) error {
	log.Println("Contact update...")
	var req dto.UpdateContact
	if err := json.Unmarshal(body, &req); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	old, err := contactService.GetContact(req.ContactID)
	if err != nil {
		if req.Email == "" {
			return nil
		}
		if err := contactService.CreateContact(entity.Contact{
			Email:     req.Email,
			Name:      req.Name,
			AccountID: req.AccountID,
			ContactID: req.ContactID,
		}); err != nil {
			return fmt.Errorf("create contact: %w", err)
		}
		return unisenderService.MakeSyncContacts(req.AccountID)
	}

	oldEmail := old.Email
	if req.Email == "" {
		req.IsSynced = false
		if err := contactService.DeleteContact(req.ContactID); err != nil {
			return fmt.Errorf("delete contact: %w", err)
		}
		return unisenderService.DeleteContact(oldEmail, req.ContactID)
	}

	req.IsSynced = false
	if err := contactService.UpdateContact(req.ContactID, req); err != nil {
		return fmt.Errorf("update contact: %w", err)
	}
	if err := unisenderService.DeleteContact(oldEmail, req.ContactID); err != nil {
		log.Printf("[warn] failed to delete old contact in UniSender: %v", err)
	}
	return unisenderService.MakeSyncContacts(req.AccountID)
}

func handleDelete(contactService *contact.ContactsUsecase, unisenderService *unisender.UnisenderService,
	body []byte) error {
	log.Println("Contact delete...")
	var req dto.DeleteContactRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	contact, err := contactService.GetContact(req.ID)
	if err != nil {
		return fmt.Errorf("get contact: %w", err)
	}
	if err := contactService.DeleteContact(req.ID); err != nil {
		return fmt.Errorf("delete contact: %w", err)
	}
	return unisenderService.DeleteContact(contact.Email, contact.ContactID)
}
