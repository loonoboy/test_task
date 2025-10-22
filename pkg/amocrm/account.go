package amocrm

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
)

const accountPath = "api/v4/account"

func (c *AMOClient) GetAccountID(account *entity.Account) (int, error) {
	resp, err := c.DoRequest(http.MethodGet, accountPath, account.Subdomain, account.AccessToken, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return 0, errors.New("unauthorized")
	}

	var dto struct {
		ID int `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&dto)
	if err != nil {
		return 0, err
	}

	return dto.ID, nil
}
