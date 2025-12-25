package service

import "github.com/hardikm9850/GoChat/internal/contacts/domain"

type ContactService interface {
	// SyncContacts takes a list of phone numbers from the client
	// and returns registered users as DTOs
	SyncContacts(contacts []string) ([]domain.ContactDTO, error)
}
