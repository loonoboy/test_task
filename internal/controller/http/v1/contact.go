package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ContactHandler struct {
	usecase ContactUsecaseInterface
}

func NewContactHandler(usecase ContactUsecaseInterface) *ContactHandler {
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
