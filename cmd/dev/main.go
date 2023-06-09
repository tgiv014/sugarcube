package main

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/app"
	"github.com/tgiv014/sugarcube/internal/logger"
	"golang.org/x/time/rate"
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

		if c.Writer.Written() || c.Writer.Status() != http.StatusNotFound {
			return
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func startVite() (error, func()) {
	// Start vite's dev server
	cmd := exec.Command("npm", "run", "dev", "--", "--clearScreen=false", "--host")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	cmd.Dir = "web"
	return cmd.Start(), func() {
		cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func waitForVite(host string) error {
	limiter := rate.NewLimiter(rate.Every(time.Second), 1)
	start := time.Now()
	for time.Since(start) < time.Second*10 {
		ctx := context.Background()
		limiter.Wait(ctx)

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, host, nil)
		if err != nil {
			return err // Fail immediately
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Warn("couldn't reach vite", "err", err)
			continue
		}

		if response.StatusCode != http.StatusOK {
			log.Warn("Vite not yet ready", "status", response.StatusCode)
			continue
		}

		return nil
	}

	return errors.New("timed out waiting for vite to start")
}

func main() {
	a := app.New()

	err, cleanup := startVite()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Wait for vite to start serving requests
	err = waitForVite("http://localhost:5173")
	if err != nil {
		log.Fatal("timed out waiting for vite to start", "err", err)
	}

	// Start router on localhost:8080
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.RequestLogger)
	r.Use(reverseProxy("http://localhost:5173"))
	a.AttachRoutes(r)

	err = r.Run()
	if err != nil {
		log.Warn("gin server exited with error", "err", err)
	}
}
