package v1

type Handler struct {
	Account      *AccountHandler
	Integrations *AccountIntegrationHandler
	Contacts     *ContactHandler
}

func NewHandler(account *AccountHandler, integrations *AccountIntegrationHandler, contact *ContactHandler) *Handler {
	return &Handler{
		Account:      account,
		Integrations: integrations,
		Contacts:     contact,
	}
}
