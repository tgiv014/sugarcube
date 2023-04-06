package session

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tgiv014/sugarcube/settings"
	"gorm.io/gorm"
)

type Service struct {
	db       *gorm.DB
	settings *settings.Service
}

func NewService(db *gorm.DB, settingsService *settings.Service) *Service {
	err := db.AutoMigrate(&Session{})
	if err != nil {
		panic(err)
	}
	return &Service{
		db:       db,
		settings: settingsService,
	}
}

func (s *Service) newSession() (*Session, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	session := Session{
		Token:   token.String(),
		Expires: time.Now().Add(24 * time.Hour),
	}
	result := s.db.Save(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

func (s *Service) Login(password string) (*Session, error) {
	settings, err := s.settings.Get()
	if err != nil {
		return nil, err
	}

	err = settings.ComparePassword([]byte(password))
	if err != nil {
		return nil, err
	}

	return s.newSession()
}

func (s *Service) Authenticate(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	session := Session{}
	result := s.db.Where("token = ?", token).First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if result.Error != nil {
		c.Abort()
		return
	}
	if session.Expires.Before(time.Now()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
