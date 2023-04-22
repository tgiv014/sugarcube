package settings

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/tgiv014/sugarcube/events"

	"gorm.io/gorm"
)

type Service struct {
	bus *events.Bus
	db  *gorm.DB
}

func NewService(bus *events.Bus, db *gorm.DB) *Service {
	err := db.AutoMigrate(&Settings{})
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		bus: bus,
		db:  db,
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

	s.bus.Emit(SettingsUpdatedEvent{})
	return nil
}
