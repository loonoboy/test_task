package dto

type IntegrationUpdate struct {
	SecretKey   *string `json:"secret_key"`
	RedirectURL *string `json:"redirect_url"`
	AuthCode    *string `json:"auth_code"`
}
