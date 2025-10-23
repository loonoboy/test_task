package v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/integrations", handler.Integrations.CreateIntegration).Methods(http.MethodPost)
	api.HandleFunc("/integrations", handler.Integrations.ListIntegration).Methods(http.MethodGet)
	api.HandleFunc("/integrations/{id}", handler.Integrations.UpdateIntegration).Methods(http.MethodPatch)
	api.HandleFunc("/integrations/{id}", handler.Integrations.GetIntegration).Methods(http.MethodGet)
	api.HandleFunc("/integrations/{id}", handler.Integrations.DeleteIntegration).Methods(http.MethodDelete)

	api.HandleFunc("/account", handler.Account.CreateAccount).Methods(http.MethodGet)
	api.HandleFunc("/accounts", handler.Account.ListAccount).Methods(http.MethodGet)
	api.HandleFunc("/accounts/{id}", handler.Account.GetAccount).Methods(http.MethodGet)
	api.HandleFunc("/accounts/{id}", handler.Account.UpdateAccount).Methods(http.MethodPatch)
	api.HandleFunc("/accounts/{id}", handler.Account.DeleteAccount).Methods(http.MethodDelete)

	api.HandleFunc("/accounts/{id}/contacts", handler.Contacts.GetAllContacts).Methods(http.MethodGet)

	return r
}
