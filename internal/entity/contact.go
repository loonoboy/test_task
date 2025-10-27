package entity

type Contact struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	FirstName         string        `json:"first_name"`
	LastName          string        `json:"last_name"`
	ResponsibleUserID int           `json:"responsible_user_id"`
	GroupID           int           `json:"group_id"`
	CreatedBy         int           `json:"created_by"`
	UpdatedBy         int           `json:"updated_by"`
	CreatedAt         int           `json:"created_at"`
	UpdatedAt         int           `json:"updated_at"`
	IsDeleted         bool          `json:"is_deleted"`
	IsUnsorted        bool          `json:"is_unsorted"`
	AccountID         int           `json:"account_id"`
	CustomFields      []CustomField `json:"custom_fields_values"`
}

type ContactsResp struct {
	Embedded struct {
		Contacts []Contact `json:"contacts"`
	} `json:"_embedded"`
}
