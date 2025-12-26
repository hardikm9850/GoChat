package repository

import "errors"

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrUserNotFound         = errors.New("user not found")
)
