package service

import (
	"errors"
	"fmt"
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	"log"
	"time"

	"github.com/google/uuid"
	autherror "github.com/hardikm9850/GoChat/internal/auth"
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	authkitjwt "github.com/hardikm9850/authkit/jwt"
	authkitpassword "github.com/hardikm9850/authkit/password"
)

type authService struct {
	userRepo repository.UserRepository
	//refresh repository.RefreshTokenRepository
	jwtManager authkitjwt.Manager
}

func New(
	userRepo repository.UserRepository,
	//refresh repository.RefreshTokenRepository,
	jwtManager authkitjwt.Manager,
) AuthService {
	return &authService{
		userRepo: userRepo,
		//	refresh:    refresh,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(phone, password, name string) error {
	log.Println("Service 👉 /auth/register hit")
	if _, err := s.userRepo.FindByMobile(phone); err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return autherror.ErrUserAlreadyExists
		}
	}

	hashedPassword, err := authkitpassword.HashPassword(password)
	if err != nil {
		return err
	}

	user := domain.User{
		ID:           uuid.NewString(),
		Name:         name,
		PhoneNumber:  phone,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(phone, password string) (Tokens, error) {
	user, err := s.userRepo.FindByMobile(phone)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return Tokens{}, autherror.ErrInvalidCredentials
		}
		return Tokens{}, err
	}

	if err := authkitpassword.VerifyPassword(password, user.PasswordHash); err != nil {
		return Tokens{}, autherror.ErrInvalidCredentials
	}

	claims := authkitjwt.Claims{
		UserID: user.ID,
	}

	accessToken, err := s.jwtManager.Generate(claims)
	if err != nil {
		return Tokens{}, fmt.Errorf("generate jwt: %w", err)
	}

	return Tokens{
		AccessToken: accessToken,
	}, nil
}
