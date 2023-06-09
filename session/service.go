package session

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tgiv014/sugarcube/settings"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrNoPassword        = errors.New("no password configured")
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

	if len(settings.HashedPassword) == 0 {
		return nil, ErrNoPassword
	}

	err = settings.ComparePassword([]byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, ErrIncorrectPassword
	}
	if err != nil {
		return nil, err
	}

	return s.newSession()
}

var (
	ErrExpired = errors.New("session expired")
	ErrInvalid = errors.New("invalid session")
)

func (s *Service) GetSession(token string) (*Session, error) {
	session := Session{}
	result := s.db.Where("token = ?", token).First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrInvalid
	}
	if result.Error != nil {
		return nil, result.Error
	}
	if session.Expires.Before(time.Now()) {
		return nil, ErrExpired
	}

	return &session, nil
}

func (s *Service) Authenticate(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, err = s.GetSession(token)
	if errors.Is(err, ErrInvalid) || errors.Is(err, ErrExpired) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.Abort()
		return
	}
}
