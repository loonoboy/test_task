package http

import (
	"encoding/json"
	"net/http"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account_integration"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type AccountIntegrationHandler struct {
	usecase account_integration.AccountIntegrationUsecaseInterface
}

func NewAccountIntegrationHandler(usecase account_integration.AccountIntegrationUsecaseInterface) *AccountIntegrationHandler {
	return &AccountIntegrationHandler{usecase: usecase}
}

func (h *AccountIntegrationHandler) CreateIntegration(w http.ResponseWriter, r *http.Request) {
	var req entity.AccountIntegration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.usecase.CreateIntegration(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AccountIntegrationHandler) GetIntegration(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path, "/integrations/")
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.GetIntegration(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	ep := json.NewEncoder(w)
	ep.SetIndent("", "  ")
	ep.Encode(resp)
}

func (h *AccountIntegrationHandler) ListIntegration(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.ListIntegrations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ep := json.NewEncoder(w)
	ep.SetIndent("", "  ")
	if err := ep.Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *AccountIntegrationHandler) UpdateIntegration(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path, "/integrations/")
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	var req dto.IntegrationUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.usecase.UpdateIntegration(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AccountIntegrationHandler) DeleteIntegration(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path, "/integrations/")
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteIntegration(id); err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
