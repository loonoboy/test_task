package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller"
	"git.amocrm.ru/study_group/in_memory_database/internal/entity"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
	"github.com/gorilla/mux"
)

type AccountHandler struct {
	usecase controller.AccountUsecaseInterface
}

func NewAccountHandler(usecase controller.AccountUsecaseInterface) *AccountHandler {
	return &AccountHandler{usecase: usecase}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req entity.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.usecase.CreateAccount(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID, err := strconv.Atoi(vars["accountID"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.GetAccount(accountID)
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
	vars := mux.Vars(r)

	accountID, err := strconv.Atoi(vars["accountID"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	var req dto.UpdateAccount
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.usecase.UpdateAccount(accountID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID, err := strconv.Atoi(vars["accountID"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteAccount(accountID); err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
