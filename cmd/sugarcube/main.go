package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/app"
	"github.com/tgiv014/sugarcube/internal/logger"
	"github.com/tgiv014/sugarcube/web"
)

func main() {
	a := app.New(app.Config{
		DBPath: "db.sqlite",
	})

	static, err := web.Static()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.RequestLogger)
	r.Use(static)
	a.AttachRoutes(r)
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet &&
			!strings.ContainsRune(c.Request.URL.Path, '.') &&
			!strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Request.URL.Path = "/"
			static(c)
		}
	})
	r.Run()
}
