package memory

import "sync"

type RefreshRepo struct {
	mu     sync.RWMutex
	tokens map[string]string // userID -> refreshToken
}

func NewRefreshRepo() *RefreshRepo {
	return &RefreshRepo{
		tokens: make(map[string]string),
	}
}

func (r *RefreshRepo) Save(userID, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[userID] = token
	return nil
}

func (r *RefreshRepo) Validate(userID, token string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stored, ok := r.tokens[userID]
	return ok && stored == token
}

func (r *RefreshRepo) Revoke(userID, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if stored, ok := r.tokens[userID]; ok && stored == token {
		delete(r.tokens, userID)
	}
	return nil
}
