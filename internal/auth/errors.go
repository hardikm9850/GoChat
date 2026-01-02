package auth

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserDoesNotExists  = errors.New("user does not exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
	ErrTokenReuseDetected = errors.New("refresh token is being reused")
)
