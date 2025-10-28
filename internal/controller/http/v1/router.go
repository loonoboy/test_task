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
	api.HandleFunc("/integrations/{integrationID}", handler.Integrations.UpdateIntegration).Methods(http.MethodPatch)
	api.HandleFunc("/integrations/{integrationID}", handler.Integrations.GetIntegration).Methods(http.MethodGet)
	api.HandleFunc("/integrations/{integrationID}", handler.Integrations.DeleteIntegration).Methods(http.MethodDelete)

	api.HandleFunc("/accounts", handler.Account.CreateAccount).Methods(http.MethodPost)
	api.HandleFunc("/accounts", handler.Account.ListAccount).Methods(http.MethodGet)
	api.HandleFunc("/accounts/{accountID}", handler.Account.GetAccount).Methods(http.MethodGet)
	api.HandleFunc("/accounts/{accountID}", handler.Account.UpdateAccount).Methods(http.MethodPatch)
	api.HandleFunc("/accounts/{accountID}", handler.Account.DeleteAccount).Methods(http.MethodDelete)

	api.HandleFunc("/accounts/{accountID}/contacts", handler.Contacts.CreateContact).Methods(http.MethodPost)
	api.HandleFunc("/accounts/{accountID}/contacts", handler.Contacts.ListContact).Methods(http.MethodGet)
	api.HandleFunc("/accounts/{accountID}/contacts/{contactID}", handler.Contacts.GetContact).Methods(http.MethodGet)
	api.HandleFunc("/accounts/{accountID}/contacts/{contactID}", handler.Contacts.UpdateContact).Methods(http.MethodPatch)
	api.HandleFunc("/accounts/{accountID}/contacts/{contactID}", handler.Contacts.DeleteContact).Methods(http.MethodDelete)

	api.HandleFunc("/auth/callback", handler.AmoClient.SaveAccountInfo).Methods(http.MethodGet)

	api.HandleFunc("/", handler.Unisender.SaveUnisenderKey).Methods(http.MethodPost)

	return r
}
