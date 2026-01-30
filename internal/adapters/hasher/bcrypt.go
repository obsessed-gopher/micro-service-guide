// Package hasher содержит реализации хэширования паролей.
package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher - реализация хэширования через bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher создаёт новый bcrypt hasher.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

// Hash хэширует пароль.
func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare сравнивает хэш с паролем.
func (h *BcryptHasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}