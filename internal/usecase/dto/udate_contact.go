package dto

type UpdateContact struct {
	ContactID int     `json:"contact_id" validate:"required"`
	AccountID *int    `json:"account_id" validate:"required"`
	Name      *string `json:"name" validate:"required"`
	Email     *string `json:"email" validate:"required"`
	IsSynced  bool    `json:"is_synced"`
}
