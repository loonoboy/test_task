package entity

type AccountIntegration struct {
	ClientID    int    `json:"client_id" validate:"required"`
	SecretKey   string `json:"secret_key" validate:"required"`
	RedirectURL string `json:"redirect_url" validate:"required"`
	AuthCode    string `json:"auth_code" validate:"required"`
}
