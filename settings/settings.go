package settings

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	HashedPassword []byte `json:"-"`
	DexcomUsername string `json:"dexcomUsername"`
	DexcomPassword string `json:"-"`
}

func (s Settings) ComparePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(s.HashedPassword, password)
}

func (s Settings) Completed() bool {
	if len(s.HashedPassword) == 0 {
		return false
	}

	return true
}
