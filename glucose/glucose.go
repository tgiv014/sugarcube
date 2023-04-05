package glucose

import (
	"time"

	"gorm.io/gorm"
)

type GlucoseReading struct {
	gorm.Model
	Value     int
	Timestamp time.Time
}
