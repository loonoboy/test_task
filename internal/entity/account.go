package entity

type Account struct {
	AccountID    int                  `json:"account_id" validate:"required" gorm:"primaryKey;autoIncrement:false"`
	AccessToken  string               `json:"access_token" validate:"required" gorm:"type:text;not null"`
	RefreshToken string               `json:"refresh_token" validate:"required" gorm:"type:text;not null"`
	Expires      int                  `json:"expires" validate:"required" gorm:"not null"`
	Subdomain    string               `json:"subdomain" gorm:"type:varchar(255);not null"`
	Integrations []AccountIntegration `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AccountID;references:AccountID"`
	Contacts     []Contact            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AccountID;references:AccountID"`
}
