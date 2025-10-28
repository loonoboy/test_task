package unisender

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type UnisenderService struct {
	AccountRepo usecase.AccountRepository
	ContactRepo usecase.ContactRepository
	HTTPClient  *http.Client
}

func NewUnisenderService(accountRepo usecase.AccountRepository, contactRepo usecase.ContactRepository) *UnisenderService {
	return &UnisenderService{
		AccountRepo: accountRepo,
		ContactRepo: contactRepo,
		HTTPClient:  &http.Client{Timeout: 15 * time.Second},
	}
}

func ptrBool(b bool) *bool       { return &b }
func ptrInt(i int) *int          { return &i }
func ptrString(s string) *string { return &s }

func buildImportContactsPayload(contacts []*entity.Contact, apiKey string, amoListID int) (string, url.Values) {
	values := url.Values{}
	values.Set("format", "json")
	values.Set("api_key", apiKey)
	values.Add("field_names[0]", "email")
	values.Add("field_names[1]", "email_list_ids")

	for i, c := range contacts {
		values.Add(fmt.Sprintf("data[%d][0]", i), c.Email)
		values.Add(fmt.Sprintf("data[%d][1]", i), strconv.Itoa(amoListID))
	}

	return "https://api.unisender.com/ru/api/importContacts", values
}

func (s *UnisenderService) markContactsSynced(accountID int, contacts []*entity.Contact) error {
	for _, c := range contacts {
		upd := dto.UpdateContact{
			AccountID: ptrInt(accountID),
			Email:     ptrString(c.Email),
			IsSynced:  ptrBool(true),
		}
		if err := s.ContactRepo.UpdateContact(c.ContactID, upd); err != nil {
			log.Printf("[Unisender] failed to update contact %d: %v", c.ContactID, err)
			fmt.Errorf("update contact %d: %w", c.ContactID, err)
		}
	}
	return nil
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

	amoListID, err := s.createOrGetList(account.UnisenderKey, "amoCRM")
	if err != nil {
		return fmt.Errorf("create or get UniSender list: %w", err)
	}

	urlStr, values := buildImportContactsPayload(contacts, account.UnisenderKey, amoListID)

	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequest(http.MethodPost, urlStr, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to UniSender failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("UniSender returned status %d: %s", resp.StatusCode, string(body))
	}

	if err := s.markContactsSynced(account.AccountID, contacts); err != nil {
		return fmt.Errorf("mark contacts synced: %w", err)
	}

	return nil
}

func (s *UnisenderService) createOrGetList(apiKey, listName string) (int, error) {
	listID := s.getListID(apiKey, listName)
	if listID > 0 {
		return listID, nil
	}
	return s.createList(apiKey, listName)
}

func (s *UnisenderService) getListID(apiKey, listName string) int {
	reqURL := fmt.Sprintf("https://api.unisender.com/ru/api/getLists?format=json&api_key=%s", apiKey)
	resp, err := s.HTTPClient.Get(reqURL)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	var res dto.GetListResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return 0
	}

	for _, v := range res.Result {
		if v.Title == listName {
			return v.ID
		}
	}
	return 0
}

func (s *UnisenderService) createList(apiKey, listName string) (int, error) {
	reqURL := fmt.Sprintf("https://api.unisender.com/ru/api/createList?format=json&api_key=%s&title=%s", apiKey, listName)
	resp, err := s.HTTPClient.Get(reqURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res dto.GetListResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return 0, err
	}

	if len(res.Result) > 0 {
		return res.Result[0].ID, nil
	}
	return 0, fmt.Errorf("failed to create UniSender list")
}

func (u *UnisenderService) SaveUnisenderKey(id int, update dto.UpdateAccount) error {
	err := u.AccountRepo.UpdateAccount(id, update)
	if err != nil {
		log.Println(err)
	}
	return nil
}
