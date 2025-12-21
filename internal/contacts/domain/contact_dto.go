package domain

type ContactDTO struct {
	ID    string `json:"user_id"`
	Phone string `json:"phone"`
	Name  string `json:"name"`
}
