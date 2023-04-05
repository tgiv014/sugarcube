package main

import (
	"log"
	"os"

	"github.com/tgiv014/sugarcube/app"
)

func getEnv() app.Environment {
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		return app.Development
	default:
		return app.Production
	}
}

func main() {
	a := app.New(app.Config{
		Environment: getEnv(),
	})

	err := a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
