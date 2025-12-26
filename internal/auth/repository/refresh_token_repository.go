package repository

import (
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	"time"
)

type RefreshTokenRepository interface {
	Create(token domain.RefreshToken) error
	FindByToken(token string) (domain.RefreshToken, error)
	FindByUserID(userID string) (domain.RefreshToken, error)
	Delete(token string) error
	DeleteByUserID(userID string) error
	UpdateByUserID(userID string, token string, time time.Time) error
}
