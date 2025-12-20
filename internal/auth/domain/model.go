package domain

import "time"

type User struct {
    ID           string
    PhoneNumber  string
    Name         string
    PasswordHash string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
