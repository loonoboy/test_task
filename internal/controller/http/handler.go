package http

import "net/http"

type Handler struct {
	Account      *AccountHandler
	Integrations *AccountIntegrationHandler
}

func NewHandler(account *AccountHandler, integrations *AccountIntegrationHandler) *Handler {
	return &Handler{
		Account:      account,
		Integrations: integrations,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Account.CreateAccount(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Account.ListAccount(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Account.GetAccount(w, r)
		case http.MethodPut:
			h.Account.UpdateAccount(w, r)
		case http.MethodDelete:
			h.Account.DeleteAccount(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/integrations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Integrations.ListIntegration(w, r)
		case http.MethodPost:
			h.Integrations.CreateIntegration(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/integrations/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Integrations.GetIntegration(w, r)
		case http.MethodPut:
			h.Integrations.UpdateIntegration(w, r)
		case http.MethodDelete:
			h.Integrations.DeleteIntegration(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
