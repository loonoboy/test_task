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
	dto := struct {
		entity.AccountIntegration
		GrantType string `json:"grant_type"`
	}{
		*integration,
		"authorization_code",
	}
	log.Println("dto", "dto", dto)

	j, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	log.Println("json", "json", string(j))

	resp, err := c.DoRequest(http.MethodPost, oAuthPath, domain, dto.Code, bytes.NewBuffer(j))
	log.Println("resp", "resp", resp)
	if err != nil {
		log.Println("err", err)
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

	log.Println("get tokens", "account", account)

	return &account, nil
}
