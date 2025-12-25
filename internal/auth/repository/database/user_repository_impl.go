package mysql

import (
	"errors"
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user domain.User) error {
	log.Println("DB 👉 /auth/register hit")
	model := UserModel{
		ID:           user.ID,
		Name:         user.Name,
		PhoneNumber:  user.PhoneNumber,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}
	err := r.db.Create(&model).Error
	if err != nil && isDuplicateKeyError(err) {
		return repository.ErrUserAlreadyExists
	}
	return nil
}

func (r *UserRepository) FindByID(id string) (domain.User, error) {
	var model UserModel

	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, repository.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return toDomainUser(model), nil
}
func (r *UserRepository) FindByMobile(mobile string) (domain.User, error) {
	var model UserModel

	err := r.db.Where("phone_number = ?", mobile).First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, repository.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return toDomainUser(model), nil
}

func (r *UserRepository) FindByMobiles(mobiles []string) ([]domain.User, error) {

	if len(mobiles) == 0 {
		return []domain.User{}, nil
	}

	var models []UserModel

	err := r.db.
		Where("phone_number IN ?", mobiles).
		Find(&models).
		Error

	if err != nil {
		return nil, err
	}

	// map DB models to domain entities
	users := make([]domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, domain.User{
			ID:          m.ID,
			PhoneNumber: m.PhoneNumber,
			Name:        m.Name,
		})
	}

	return users, nil
}
