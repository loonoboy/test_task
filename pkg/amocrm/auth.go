package amocrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
)

const (
	oAuthPath = "oauth2/access_token"
)

func (c *AMOClient) GetTokens(integration *entity.AccountIntegration, domain string) (*entity.Account, error) {
	dto := &response{
		*integration,
		"authorization_code",
	}

	j, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRequest(http.MethodPost, oAuthPath, domain, dto.Code, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("unexpected status code", "code", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var account entity.Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
