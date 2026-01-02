package mysql

import "time"

type UserModel struct {
	ID              string `gorm:"primaryKey;size:36`
	PhoneNumber     string `gorm:"column:phone_number;uniqueIndex;size:15;not null"`
	Name            string `gorm:"column:name"`
	PasswordHash    string `gorm:"not null"`
	PhoneNumberHash string `gorm:"not null"`
	CountryCode     string `gorm:"column:country_code;not null"`
	CreatedAt       time.Time
}

func (UserModel) TableName() string {
	return "users"
}
