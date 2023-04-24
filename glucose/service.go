package glucose

import (
	"errors"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mattn/go-sqlite3"
	dexcomshare "github.com/tgiv014/dexcom-share"
	"github.com/tgiv014/sugarcube/events"
	"github.com/tgiv014/sugarcube/settings"
	"gorm.io/gorm"
)

type Service struct {
	bus      *events.Bus
	db       *gorm.DB
	settings *settings.Service
	fetchNow chan fetchParams
}

func NewService(bus *events.Bus, db *gorm.DB, settings *settings.Service) *Service {
	err := db.AutoMigrate(&GlucoseReading{})
	if err != nil {
		panic(err)
	}

	s := &Service{
		bus:      bus,
		db:       db,
		settings: settings,
	}

	go s.glucoseFetcher()

	return s
}

func (s *Service) glucoseFetcher() {
	fetchTimer := time.NewTicker(time.Minute)
	s.fetchNow = make(chan fetchParams)

	// Do a large initial fetch to try and backfill
	go func() {
		s.fetchNow <- fetchParams{
			minutes:  60 * 12,
			maxCount: 12 * 12,
		}

		updateEvents, close := s.bus.Subscribe("settingsUpdated")
		defer close()

		for range updateEvents {
			log.Info("Fetching glucose readings because of a settings update")
			s.fetchNow <- fetchParams{
				minutes:  60 * 12,
				maxCount: 12 * 12,
			}
		}
	}()

	for {
		select {
		case <-fetchTimer.C:
			log.Info("glucose fetch triggered by timer")
			s.fetchAndUpdateTicker(fetchTimer, fetchParams{
				minutes:  20,
				maxCount: 2,
			})
		case params := <-s.fetchNow:
			log.Info("glucose fetch triggered manually")
			s.fetchAndUpdateTicker(fetchTimer, params)
		}
	}
}

type fetchParams struct {
	minutes  int
	maxCount int
}

func (s *Service) fetchAndUpdateTicker(ticker *time.Ticker, params fetchParams) error {
	t, err := s.fetchGlucoseEntries(params)
	if err != nil {
		log.Warn("glucose fetch failed", "err", err)
		return err
	}
	ticker.Reset(t)
	log.Info("Next glucose fetch in", "t", t)
	return nil
}

// fetchGlucoseEntries pulls down glucose entries from dexcom share
// If successful, it returns a time representing when we expect the next value to appear
func (s *Service) fetchGlucoseEntries(params fetchParams) (time.Duration, error) {
	settings, err := s.settings.Get()
	if err != nil {
		return 0, err
	}

	if settings.DexcomUsername == "" {
		return 0, errors.New("dexcom username empty")
	}

	if settings.DexcomPassword == "" {
		return 0, errors.New("dexcom password empty")
	}

	client, err := dexcomshare.NewClient(settings.DexcomUsername, settings.DexcomPassword)
	if err != nil {
		return 0, err
	}

	entries, err := client.ReadGlucose(params.minutes, params.maxCount)
	if err != nil {
		return 0, err
	}

	newValues := 0
	for _, entry := range entries {
		reading, err := GlucoseReadingFromDexcomShare(entry)
		if err != nil {
			return 0, err
		}
		result := s.db.Create(reading)
		var sqliteErr sqlite3.Error
		if errors.As(result.Error, &sqliteErr) && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			continue
		}
		if result.Error != nil {
			log.Warn("error saving glucose reading", "err", result.Error)
			continue
		}
		newValues++
	}

	// No entries returned at all, check back in a minute
	// Sensor is likely in warmup
	if len(entries) == 0 {
		return time.Minute, nil
	}

	// Get last reading to determine next refresh time
	lastReading, err := GlucoseReadingFromDexcomShare(entries[0])
	if err != nil {
		return 0, err
	}

	// Start checking 5 minutes after the most recent reading
	nextFetch := lastReading.Timestamp.Add(5 * time.Minute)
	wait := time.Until(nextFetch)
	if wait < 0 {
		wait = 5 * time.Second
	}

	if newValues > 0 {
		log.Info("Newest glucose entry!", "value", lastReading.Value, "t", lastReading.Timestamp)
		s.bus.Emit(NewReadingsEvent{})
	}

	return wait.Round(time.Second), nil
}

type GetReadingsParams struct {
	StartTime time.Time
	EndTime   time.Time
}

func (s *Service) GetReadings(start, end time.Time) ([]GlucoseReading, error) {
	var readings []GlucoseReading
	result := s.db.Limit(1000).Where("timestamp >= ? AND timestamp < ?", start.UTC(), end.UTC()).Find(&readings)
	if result.Error != nil {
		return nil, result.Error
	}

	return readings, nil
}
