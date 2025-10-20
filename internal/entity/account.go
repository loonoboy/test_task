package entity

type Account struct {
	AccountID    int    `json:"account_id" validate:"required"`
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
	Expires      int    `json:"expires" validate:"required"`
}
