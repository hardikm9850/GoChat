package service

import (
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	"github.com/hardikm9850/GoChat/internal/contacts/domain"
	"strings"
)

type ContactServiceImpl struct {
	userRepository repository.UserRepository
}

func (c *ContactServiceImpl) SyncContacts(contacts []string) ([]domain.ContactDTO, error) {
	contacts = normalizePhones(contacts)

	users, err := c.userRepository.FindByMobiles(contacts)
	if err != nil {
		return nil, err
	}

	contactDTOs := make([]domain.ContactDTO, 0, len(users))
	for _, user := range users {
		contactDTOs = append(contactDTOs, toContactDTO(user))
	}

	return contactDTOs, nil
}

func New(userRepository repository.UserRepository) ContactService {
	return &ContactServiceImpl{
		userRepository: userRepository,
	}
}

func normalizePhones(input []string) []string {
	set := map[string]struct{}{}

	for _, p := range input {
		cleaned := strings.TrimSpace(p)
		cleaned = strings.ReplaceAll(cleaned, " ", "")
		cleaned = strings.ReplaceAll(cleaned, "-", "")
		cleaned = strings.ReplaceAll(cleaned, "(", "")
		cleaned = strings.ReplaceAll(cleaned, ")", "")

		if cleaned != "" {
			set[cleaned] = struct{}{}
		}
	}

	phones := make([]string, 0, len(set))
	for p := range set {
		phones = append(phones, p)
	}

	return phones
}
