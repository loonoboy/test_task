package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/gorilla/mux"
)

type ContactHandler struct {
	usecase ContactsUsecaseInterface
}

func NewContactHandler(usecase ContactsUsecaseInterface) *ContactHandler {
	return &ContactHandler{usecase: usecase}
}

func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var req entity.Contact
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.usecase.CreateContact(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ContactHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contactID, err := strconv.Atoi(vars["contactID"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.GetContact(contactID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	ep := json.NewEncoder(w)
	ep.SetIndent("", "  ")
	ep.Encode(resp)
}

func (h *ContactHandler) ListContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contactID, err := strconv.Atoi(vars["contactID"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.ListContacts(contactID)
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

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contactID, err := strconv.Atoi(vars["contactID"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	var req dto.UpdateContact
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.usecase.UpdateContact(contactID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contactID, err := strconv.Atoi(vars["contactID"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteContact(contactID); err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
