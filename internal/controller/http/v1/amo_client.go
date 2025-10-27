package v1

import (
	"net/http"

	"github.com/google/uuid"
)

type AmoClientHandler struct {
	usecase AmoClientUsecaseInterface
}

func NewAmoClientHandler(usecase AmoClientUsecaseInterface) *AmoClientHandler {
	return &AmoClientHandler{usecase: usecase}
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
