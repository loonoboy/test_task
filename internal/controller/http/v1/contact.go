package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/contact"
	"github.com/gorilla/mux"
)

type ContactHandler struct {
	usecase contact.ContactUsecaseInterface
}

func NewContactHandler(usecase contact.ContactUsecaseInterface) *ContactHandler {
	return &ContactHandler{usecase: usecase}
}

func (h *ContactHandler) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountIDStr := vars["id"]

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid account id"}`, http.StatusBadRequest)
		return
	}

	contacts, err := h.usecase.GetAllContacts(accountID)
	if err != nil {
		http.Error(w, `{"error":"invalid account id"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ep := json.NewEncoder(w)
	ep.SetIndent("", "  ")
	ep.Encode(contacts)
}
