package app

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

// reverseProxy provides middleware to proxy unhandled requests to vite
func reverseProxy(target string) gin.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		c.Next()
		fmt.Println(c.Writer.Status())
		if c.Writer.Written() || c.Writer.Status() != http.StatusNotFound {
			return
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (a *App) runDev() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Start vite's dev server
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
		r.Use(reverseProxy("http://localhost:5173"))
		a.attachRoutes(r)

		err := r.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	return nil
}
