package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/google/uuid"
)

type AmoClientHandler struct {
	usecase controller.AmoClientUsecaseInterface
	queue   controller.QueueInterface
}

func NewAmoClientHandler(usecase controller.AmoClientUsecaseInterface, queue controller.QueueInterface) *AmoClientHandler {
	return &AmoClientHandler{
		usecase: usecase,
		queue:   queue,
	}
}

func (h *AmoClientHandler) SaveAccountInfo(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	domain := r.URL.Query().Get("referer")
	errorParam := r.URL.Query().Get("error")
	clientID := r.URL.Query().Get("client_id")

	if errorParam == "access_denied" {
		http.Error(w, "Access denied for this request", http.StatusForbidden)
		return
	}

	if authCode == "" {
		http.Error(w, "authorization authCode is missing", http.StatusBadRequest)
		return
	}

	if domain == "" {
		http.Error(w, "domain is missing", http.StatusBadRequest)
		return
	}

	clientUUID, err := uuid.Parse(clientID)
	if err != nil {
		http.Error(w, "client_id is invalid", http.StatusBadRequest)
		return
	}

	err = h.usecase.SaveAccountInfo(authCode, domain, clientUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AmoClientHandler) infoFromWebHook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	eventType := h.getIventType(values)

	switch eventType {
	case "add":
		h.jobCreate(values, eventType)
	case "update":
		h.jobUpdate(values, eventType)
	case "delete":
		h.jobDelete(values, eventType)
	}
}

func (h *AmoClientHandler) getIventType(values url.Values) string {
	for v := range values {
		if strings.HasPrefix(v, "contacts[add]") {
			return "add"
		}
		if strings.HasPrefix(v, "contacts[update]") {
			return "update"
		}
		if strings.HasPrefix(v, "contacts[delete]") {
			return "delete"
		}
	}
	return ""
}

func (h *AmoClientHandler) jobCreate(values url.Values, eventType string) {
	accountId, _ := strconv.Atoi(values.Get(fmt.Sprintf("account[id]")))
	id, _ := strconv.Atoi(values.Get(fmt.Sprintf("contacts[%s][0][id]", eventType)))
	name := values.Get(fmt.Sprintf("contacts[%s][0][name]", eventType))
	email := values.Get(fmt.Sprintf("contacts[%s][0][custom_fields][0][values][0][value]", eventType))
	var req dto.UpdateContact
	req.AccountID = &accountId
	req.ContactID = id
	req.Name = &name
	req.Email = &email

	body, _ := json.Marshal(req)
	log.Println("adding a job to create contact for account -  ", accountId)
	if h.queue.AddJob(body, "contact") != nil {
		log.Fatal("error adding job")
	}
}
func (h *AmoClientHandler) jobUpdate(values url.Values, eventType string) {
	accountId, _ := strconv.Atoi(values.Get(fmt.Sprintf("account[id]")))
	id, _ := strconv.Atoi(values.Get(fmt.Sprintf("contacts[%s][0][id]", eventType)))
	name := values.Get(fmt.Sprintf("contacts[%s][0][name]", eventType))
	email := values.Get(fmt.Sprintf("contacts[%s][0][custom_fields][0][values][0][value]", eventType))
	var req dto.UpdateContact
	req.AccountID = &accountId
	req.ContactID = id
	req.Name = &name
	req.Email = &email

	body, _ := json.Marshal(req)
	log.Println("adding a job to update contact for account -  ", accountId)
	if h.queue.AddJob(body, "contact") != nil {
		log.Fatal("error adding job")
	}
}
func (h *AmoClientHandler) jobDelete(values url.Values, eventType string) {
	id, _ := strconv.Atoi(values.Get(fmt.Sprintf("contacts[%s][0][id]", eventType)))

	body, _ := json.Marshal(map[string]int{"id": id})
	log.Println("adding a job to delete contact -  ", id)
	if h.queue.AddJob(body, "contact") != nil {
		log.Fatal("error adding job")
	}
}
