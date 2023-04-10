package app

import (
	"context"
	"os"

	"github.com/charmbracelet/log"

	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/events"
	"github.com/tgiv014/sugarcube/glucose"
	"github.com/tgiv014/sugarcube/session"
	"github.com/tgiv014/sugarcube/settings"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	Config   Config
	bus      *events.Bus
	db       *gorm.DB
	settings *settings.Service
	sessions *session.Service
	glucose  *glucose.Service
}

type Config struct {
	Environment Environment
	DBPath      string
}

type Environment string

var (
	Development Environment = "dev"
	Production  Environment = "prod"
)

func New(config Config) *App {
	bus := events.NewBus()

	log.Info("connecting to db")
	db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{
		// Shhh
		Logger: logger.New(log.New(os.Stdout), logger.Config{
			LogLevel: logger.Silent,
		}),
	})
	if err != nil {
		log.Fatal(err)
	}

	settingsService := settings.NewService(db)
	sessionService := session.NewService(db, settingsService)
	glucoseService := glucose.NewService(bus, db, settingsService)

	a := &App{
		Config:   config,
		bus:      bus,
		db:       db,
		settings: settingsService,
		sessions: sessionService,
		glucose:  glucoseService,
	}

	a.db = db

	return a
}

func (a *App) Run(ctx context.Context) error {
	// gin debug logs are pretty loud
	gin.SetMode(gin.ReleaseMode)
	// Start API
	if a.Config.Environment == Development {
		log.Info("Starting in dev mode")
		return a.runDev()
	}

	log.Info("Starting")
	return a.runProd()
}

func (a *App) attachRoutes(r gin.IRouter) {
	api := r.Group("/api")
	{
		api.POST("/login", a.login)
		api.POST("/signup", a.signup)
		api.GET("/status", a.status)
		api.GET("/settings", a.sessions.Authenticate, a.getSettings)
		api.PATCH("/settings", a.sessions.Authenticate, a.updateSettings)
		api.GET("/readings", a.sessions.Authenticate, a.getReadings)
		api.GET("/bus", a.bus.Handler)
	}
}
