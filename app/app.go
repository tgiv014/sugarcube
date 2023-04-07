package app

import (
	"context"

	"github.com/charmbracelet/log"

	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/glucose"
	"github.com/tgiv014/sugarcube/session"
	"github.com/tgiv014/sugarcube/settings"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	Config   Config
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
	log.Info("connecting to db")
	db, err := gorm.Open(sqlite.Open(config.DBPath))
	if err != nil {
		log.Fatal(err)
	}

	settingsService := settings.NewService(db)
	sessionService := session.NewService(db, settingsService)
	glucoseService := glucose.NewService(db, settingsService)

	a := &App{
		Config:   config,
		db:       db,
		settings: settingsService,
		sessions: sessionService,
		glucose:  glucoseService,
	}

	a.db = db

	return a
}

func (a *App) Run(ctx context.Context) error {
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
	}
}
