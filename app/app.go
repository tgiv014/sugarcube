package app

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	Config Config
	db     *gorm.DB
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
	a := &App{Config: config}

	db, err := gorm.Open(sqlite.Open(a.Config.DBPath))
	if err != nil {
		log.Fatal(err)
	}

	a.db = db

	return a
}

func (a *App) Run(ctx context.Context) error {
	// Start scheduler and set up sync
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait() // Don't leave without waiting for the scheduler to close

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Ask scheduler to stop before leaving

	go func(ctx context.Context) {
		defer wg.Done()
		a.runScheduler(ctx)
	}(ctx)

	// Start API
	if a.Config.Environment == Development {
		return a.runDev()
	}
	return a.runProd()
}

func (a *App) attachRoutes(r gin.IRouter) {
	api := r.Group("/api")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"a": "b",
		})
	})
}
