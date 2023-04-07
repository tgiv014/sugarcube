package glucose

import (
	"errors"
	"time"

	"github.com/charmbracelet/log"
	dexcomshare "github.com/tgiv014/dexcom-share"
	"github.com/tgiv014/sugarcube/settings"
	"gorm.io/gorm"
)

type Service struct {
	db       *gorm.DB
	settings *settings.Service
	fetchNow chan fetchParams
}

func NewService(db *gorm.DB, settings *settings.Service) *Service {
	err := db.AutoMigrate(&GlucoseReading{})
	if err != nil {
		panic(err)
	}

	s := &Service{
		db:       db,
		settings: settings,
	}

	go s.glucoseFetcher()

	return s
}

func (s *Service) glucoseFetcher() {
	lastFetch := time.Time{}
	timer := time.NewTimer(time.Minute)
	s.fetchNow = make(chan fetchParams)

	go func() {
		s.fetchNow <- fetchParams{
			minutes:  60,
			maxCount: 60,
		}
	}()

	for {
		select {
		case <-timer.C:
			log.Info("glucose fetch triggered by timer")
			params := fetchParams{
				minutes:  2,
				maxCount: 2,
			}
			if time.Since(lastFetch) > time.Hour {
				params = fetchParams{
					minutes:  120,
					maxCount: 120,
				}
			}
			err := s.fetchGlucoseEntries(params)
			if err != nil {
				log.Warn("glucose fetch failed", "err", err)
				continue
			}
			lastFetch = time.Now()
			continue
		case params := <-s.fetchNow:
			log.Info("glucose fetch triggered manually")
			err := s.fetchGlucoseEntries(params)
			if err != nil {
				log.Warn("manual glucose fetch failed", "err", err)
				continue
			}
			lastFetch = time.Now()
			continue
		}
	}
}

type fetchParams struct {
	minutes  int
	maxCount int
}

// fetchGlucoseEntries pulls down glucose entries from dexcom share
// If successful, it returns a time representing when we expect the next value to appear
func (s *Service) fetchGlucoseEntries(params fetchParams) error {
	settings, err := s.settings.Get()
	if err != nil {
		return err
	}

	if settings.DexcomUsername == "" {
		return errors.New("dexcom username empty")
	}

	if settings.DexcomPassword == "" {
		return errors.New("dexcom password empty")
	}

	client, err := dexcomshare.NewClient(settings.DexcomUsername, settings.DexcomPassword)
	if err != nil {
		return err
	}

	entries, err := client.ReadGlucose(params.minutes, params.maxCount)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		reading, err := GlucoseReadingFromDexcomShare(entry)
		if err != nil {
			return err
		}
		result := s.db.Save(reading)
		if result.Error != nil {
			log.Warn("error saving glucose reading", "err", result.Error)
			return result.Error
		}
	}

	return nil
}

type GetReadingsParams struct {
	StartTime time.Time
	EndTime   time.Time
}

func (s *Service) GetReadings(start, end time.Time) ([]GlucoseReading, error) {
	var readings []GlucoseReading
	result := s.db.Where("timestamp >= ? AND timestamp < ?", start, end).Find(&readings)
	if result.Error != nil {
		return nil, result.Error
	}

	return readings, nil
}
