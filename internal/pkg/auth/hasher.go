package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
}

func (h *Hasher) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {

		return "", err

	}

	return string(hashBytes), nil
}

func (h *Hasher) Verify(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewHasher() *Hasher {
	return &Hasher{}
}
