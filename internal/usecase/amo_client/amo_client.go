package amo_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"git.amocrm.ru/study_group/in_memory_database/pkg/amocrm"
	"github.com/google/uuid"
)

type AmoClientService struct {
	accRepo     usecase.AccountRepository
	intgrRepo   usecase.IntegrationRepository
	contactRepo usecase.ContactRepository
	client      *amocrm.AMOClient
}

func NewAmoClientServiceService(client *amocrm.AMOClient, repository usecase.AccountRepository,
	intgrRepo usecase.IntegrationRepository, contactRepo usecase.ContactRepository) *AmoClientService {
	return &AmoClientService{
		client:      client,
		accRepo:     repository,
		intgrRepo:   intgrRepo,
		contactRepo: contactRepo,
	}
}

func (s *AmoClientService) SaveAccountInfo(authCode, subdomain string, clientID uuid.UUID) error {
	integration, err := s.intgrRepo.GetIntegration(clientID)
	if err != nil {
		return err
	}
	integration.Code = authCode

	account, err := s.client.GetTokens(integration, subdomain)
	if err != nil {
		return err
	}
	account.Subdomain = subdomain

	accountID, err := s.client.GetAccountID(account)
	if err != nil {
		return err
	}

	account.AccountID = accountID

	err = s.accRepo.CreateAccount(account)
	if err != nil {
		return err
	}
	integrationUpdates := &dto.IntegrationUpdate{
		AccountID: accountID,
	}

	err = s.intgrRepo.UpdateIntegration(clientID, *integrationUpdates)
	if err != nil {
		return err
	}

	account, err = s.accRepo.GetAccount(accountID)
	if err != nil {
		return err
	}

	contacts, err := s.client.GetContacts(account)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		if err = s.contactRepo.CreateContact(&contact); err != nil {
			return err
		}
	}
	return nil
}

func (s *AmoClientService) RegisterWebHook(id int, webHookURL, subdomain string) error {
	var set []string
	set = append(set, "add_contact", "update_contact", "delete_contact")
	req := &dto.RegisterWebHookRequest{
		Destination: webHookURL,
		Settings:    set,
		Sort:        10,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	accInfo, err := s.accRepo.GetAccount(id)
	if err != nil {
		return err
	}
	reqUrl := fmt.Sprintf("https://%v/api/v4/webhooks", subdomain)
	httpReq, err := http.NewRequest("POST", reqUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accInfo.AccessToken))
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
