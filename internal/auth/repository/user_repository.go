package repository

import (
	"errors"
	"github.com/hardikm9850/GoChat/internal/auth/domain"
)

type UserRepository interface {
	Create(user domain.User) error
	FindByID(id string) (domain.User, error)
	FindByMobile(mobile string) (domain.User, error)
	FindByMobiles(mobile []string) ([]domain.User, error)
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
