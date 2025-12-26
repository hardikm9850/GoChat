package mysql

import (
	"errors"
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	"gorm.io/gorm"
	"time"
)

type RefreshTokenRepo struct {
	db *gorm.DB
}

func NewMySQLRefreshTokenRepo(db *gorm.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Create(token domain.RefreshToken) error {
	return r.db.Create(&token).Error
}

func (r *RefreshTokenRepo) FindByToken(token string) (domain.RefreshToken, error) {
	var rt domain.RefreshToken
	if err := r.db.First(&rt, "token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.RefreshToken{}, repository.ErrRefreshTokenNotFound
		}
		return domain.RefreshToken{}, err
	}
	return rt, nil
}

func (r *RefreshTokenRepo) FindByUserID(userId string) (domain.RefreshToken, error) {
	var rt domain.RefreshToken
	if err := r.db.First(&rt, "user_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.RefreshToken{}, repository.ErrRefreshTokenNotFound
		}
		return domain.RefreshToken{}, err
	}
	return rt, nil
}
func (r *RefreshTokenRepo) Delete(token string) error {
	return r.db.Delete(&domain.RefreshToken{}, "token = ?", token).Error
}

func (r *RefreshTokenRepo) DeleteByUserID(userID string) error {
	return r.db.Delete(&domain.RefreshToken{}, "user_id = ?", userID).Error
}

func (r *RefreshTokenRepo) UpdateByUserID(userID string, token string, expiresAt time.Time) error {
	return r.db.Model(&domain.RefreshToken{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
			"created_at": time.Now(),
		}).Error
}
