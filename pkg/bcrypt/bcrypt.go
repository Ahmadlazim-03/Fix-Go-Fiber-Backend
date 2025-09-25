package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptUtil struct {
	cost int
}

func NewBcryptUtil(cost int) *BcryptUtil {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptUtil{
		cost: cost,
	}
}

func (h *BcryptUtil) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

func (h *BcryptUtil) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Legacy methods for backward compatibility
type BcryptHelper interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) error
}

type bcryptHelper struct {
	*BcryptUtil
}

func NewBcryptHelper(cost int) BcryptHelper {
	return &bcryptHelper{
		BcryptUtil: NewBcryptUtil(cost),
	}
}

func (h *bcryptHelper) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}