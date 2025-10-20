package account

type UpdateInput struct {
	AccessToken  *string
	RefreshToken *string
	Expires      *int
}
