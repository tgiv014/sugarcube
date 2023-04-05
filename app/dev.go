package app

import (
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

func (a *App) runDev() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Start vite's dev server, which is configured to proxy `/api` requests to localhost:8080
	go func() {
		defer wg.Done()

		cmd := exec.Command("npm", "run", "dev", "--", "--clearScreen=false", "--host")
		cmd.Dir = "web"
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Start router on localhost:8080
	go func() {
		defer wg.Done()

		r := gin.Default()

		a.attachRoutes(r)

		err := r.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	return nil
}
