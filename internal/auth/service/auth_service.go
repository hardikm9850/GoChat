package service

type Tokens struct {
	AccessToken string
	//RefreshToken string
}

// Service ========= Auth Contract =============
type AuthService interface {
	Register(phone, password, name string) error
	Login(phone, password string) (Tokens, error)
	//Refresh(email string) (Tokens, error)
}
