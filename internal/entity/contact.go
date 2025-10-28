package entity

type Contact struct {
	ContactID int    `json:"contact_id" validate:"required" gorm:"primaryKey;autoIncrement:false"`
	AccountID int    `json:"account_id" validate:"required" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name      string `json:"name" validate:"required" gorm:"type:varchar(255) CHARACTER SET utf8mb4;not null"`
	Email     string `json:"email" validate:"required" gorm:"type:varchar(255) CHARACTER SET utf8mb4;not null;unique"`
}

type ContactsResp struct {
	Embedded ContactsEmbedded `json:"_embedded"`
}

type ContactsEmbedded struct {
	Contacts []Contact `json:"contact"`
}
