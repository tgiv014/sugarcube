package glucose

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	dexcomshare "github.com/tgiv014/dexcom-share"
	"gorm.io/gorm"
)

var (
	reDate = regexp.MustCompile(`Date\(([0-9]+)\)`)
)

type GlucoseReading struct {
	gorm.Model
	Value     int
	Timestamp time.Time
}

func GlucoseReadingFromDexcomShare(ds dexcomshare.GlucoseEntry) (*GlucoseReading, error) {
	matches := reDate.FindStringSubmatch(ds.WT)
	if len(matches) != 2 {
		return nil, errors.New("unexpected timestamp format")
	}
	unixtimeMillis, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, err
	}
	reading := &GlucoseReading{
		Timestamp: time.UnixMilli(unixtimeMillis),
		Value:     ds.Value,
	}

	// Readings should be unique by timestamp
	reading.ID = uint(unixtimeMillis / 1000)
	return reading, nil
}
