package v1

type Handler struct {
	Account      *AccountHandler
	Integrations *AccountIntegrationHandler
	Contacts     *ContactHandler
	AmoClient    *AmoClientHandler
	Unisender    *UnisenderHandler
}

func NewHandler(account *AccountHandler, integrations *AccountIntegrationHandler, contacts *ContactHandler,
	amoClient *AmoClientHandler, unisender *UnisenderHandler) *Handler {
	return &Handler{
		Account:      account,
		Integrations: integrations,
		Contacts:     contacts,
		AmoClient:    amoClient,
		Unisender:    unisender,
	}
}
