package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashService defines the interface for hashing passwords
type HashService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) bool
}

// BcryptService implementation
type BcryptService struct {
	Cost int
}

func NewHashService() HashService {
	return &BcryptService{Cost: bcrypt.DefaultCost}
}

func (s *BcryptService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.Cost)
	return string(bytes), err
}

func (s *BcryptService) Compare(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
