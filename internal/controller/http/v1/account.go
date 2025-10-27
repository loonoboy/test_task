package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/google/uuid"
)

type AccountHandler struct {
	usecase AccountUsecaseInterface
}

func NewAccountHandler(usecase AccountUsecaseInterface) *AccountHandler {
	return &AccountHandler{usecase: usecase}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
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

	err = h.usecase.CreateAccount(authCode, domain, clientUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path, "/accounts/")
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.GetAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	ep := json.NewEncoder(w)
	ep.SetIndent("", "  ")
	ep.Encode(resp)
}

func (h *AccountHandler) ListAccount(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.ListAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}

func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path, "/accounts/")
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	var req dto.UpdateAccount
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(req)
	err = h.usecase.UpdateAccount(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path, "/accounts/")
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteAccount(id); err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
