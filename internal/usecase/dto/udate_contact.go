package dto

type UpdateContact struct {
	AccountID *int    `json:"account_id" validate:"required"`
	Name      *string `json:"name" validate:"required"`
	Email     *string `json:"email" validate:"required"`
}
