package mysql

import (
	"github.com/hardikm9850/GoChat/internal/auth/domain"
)

func toDomainUser(model UserModel) domain.User {
	return domain.User{
		ID:           model.ID,
		PhoneNumber:  model.PhoneNumber,
		PasswordHash: model.PasswordHash,
		CreatedAt:    model.CreatedAt,
	}
}
