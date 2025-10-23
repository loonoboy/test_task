package amocrm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
)

const contactsPath = "api/v4/contacts"

func (c *AMOClient) GetContacts(account *entity.Account) ([]entity.Contact, error) {
	resp, err := c.DoRequest(http.MethodGet, contactsPath, account.Subdomain, account.AccessToken, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dto entity.ContactsResp
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, err
	}

	return dto.Embedded.Contacts, nil
}
