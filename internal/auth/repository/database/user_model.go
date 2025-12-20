package mysql

import "time"

type UserModel struct {
    ID           string `gorm:"primaryKey;size:36`
    PhoneNumber  string `gorm:"column:phone_number;uniqueIndex;size:15;not null"`
    Name         string `gorm:"column:name"`
    PasswordHash string `gorm:"not null"`
    CreatedAt    time.Time
}

func (UserModel) TableName() string {
    return "users"
}
