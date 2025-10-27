package entity

type Contact struct {
	ContactID int    `gorm:"primaryKey;autoIncrement:false"`
	AccountID int    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name      string `gorm:"type:varchar(255) CHARACTER SET utf8mb4;not null"`
	Email     string `gorm:"type:varchar(255) CHARACTER SET utf8mb4;not null"`
}

type ContactsResp struct {
	Embedded ContactsEmbedded `json:"_embedded"`
}

type ContactsEmbedded struct {
	Contacts []Contact `json:"contacts"`
}
