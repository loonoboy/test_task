package v1

type Handler struct {
	Account      *AccountHandler
	Integrations *AccountIntegrationHandler
	Contacts     *ContactHandler
	AmoClient    *AmoClientHandler
}

func NewHandler(account *AccountHandler, integrations *AccountIntegrationHandler, contacts *ContactHandler, amoClient *AmoClientHandler) *Handler {
	return &Handler{
		Account:      account,
		Integrations: integrations,
		Contacts:     contacts,
		AmoClient:    amoClient,
	}
}
