package accounts

type UpdateInput struct {
	AccessToken  *string
	RefreshToken *string
	Expires      *int
}
