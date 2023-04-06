package session

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Token   string
	Expires time.Time
}
