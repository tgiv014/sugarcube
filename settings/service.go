package settings

import (
	"errors"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	err := db.AutoMigrate(&Settings{})
	if err != nil {
		panic(err)
	}

	return &Service{
		db: db,
	}
}

func (s *Service) Get() (*Settings, error) {
	// Retrieve latest settings object from DB
	var settings Settings
	result := s.db.Last(&settings)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result = s.db.Save(&settings)
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &settings, nil
}

func (s *Service) Save(settings *Settings) error {
	result := s.db.Save(settings)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
