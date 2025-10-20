package dto

type UpdateAccount struct {
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
	Expires      *int    `json:"expires"`
}
