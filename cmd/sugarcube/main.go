package main

import (
	"context"
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
	ctx := context.Background()
	a := app.New(app.Config{
		Environment: getEnv(),
	})

	err := a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
