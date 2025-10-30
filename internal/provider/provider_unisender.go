package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type UnisenderProvider struct {
	client *http.Client
}

func NewUnisenderProvider() *UnisenderProvider {
	return &UnisenderProvider{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}
func (p *UnisenderProvider) CreateOrGetList(apiKey, listName string) (int, error) {
	id := p.getListID(apiKey, listName)
	if id > 0 {
		return id, nil
	}
	return p.createList(apiKey, listName)
}

func (p *UnisenderProvider) ImportContacts(apiKey string, amoListID int, contacts []*entity.Contact) error {
	urlStr, values := buildImportContactsPayload(contacts, apiKey, amoListID)
	reqBody := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, urlStr, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request to UniSender failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("UniSender returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

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

func (p *UnisenderProvider) getListID(apiKey, listName string) int {
	reqURL := fmt.Sprintf("https://api.unisender.com/ru/api/getLists?format=json&api_key=%s", apiKey)
	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	resp, err := p.client.Do(req)
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

func (p *UnisenderProvider) createList(apiKey, listName string) (int, error) {
	reqURL := fmt.Sprintf("https://api.unisender.com/ru/api/createList?format=json&api_key=%s&title=%s", apiKey, listName)
	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	resp, err := p.client.Do(req)
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

func (p *UnisenderProvider) ExcludeContact(apiKey, email string, listID int) error {
	values := url.Values{}
	values.Set("format", "json")
	values.Set("api_key", apiKey)
	values.Set("contact", email)
	values.Set("contact_type", "email")
	values.Set("list_ids", strconv.Itoa(listID))

	reqURL := "https://api.unisender.com/ru/api/exclude"
	reqBody := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, reqURL, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request to UniSender failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("UniSender returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
