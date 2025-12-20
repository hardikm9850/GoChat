package repository

type RefreshTokenRepository interface {
	Save(userID, token string) error
	Validate(userID, token string) bool
	Revoke(userID, token string) error
}
