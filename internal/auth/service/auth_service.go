package service

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// AuthService
// ========= Auth Contract =============
type AuthService interface {
	Register(phone, password, name string) error
	Login(phone, password string) (Tokens, error)
	RefreshAccessToken(refreshToken string) (Tokens, error)
	Logout(id string) error
}
