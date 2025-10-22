package dto

type IntegrationUpdate struct {
	AccountID   int     `json:"account_id"`
	SecretKey   *string `json:"secret_key"`
	RedirectURL *string `json:"redirect_url"`
	AuthCode    *string `json:"auth_code"`
}
