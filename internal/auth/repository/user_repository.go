package repository

import (
	"github.com/hardikm9850/GoChat/internal/auth/domain"
)

type UserRepository interface {
	Create(user domain.User) error
	FindByID(id string) (domain.User, error)
	FindByMobile(mobile, countryCode string) (domain.User, error)
	FindByMobiles(mobile []string) ([]domain.User, error)
}
