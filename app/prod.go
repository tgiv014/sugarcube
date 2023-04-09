package app

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/gin-gonic/gin"
	"github.com/tgiv014/sugarcube/web"
)

func (a *App) runProd() error {
	static, err := web.Static()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger)
	r.Use(static)
	a.attachRoutes(r)
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet &&
			!strings.ContainsRune(c.Request.URL.Path, '.') &&
			!strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Request.URL.Path = "/"
			static(c)
		}
	})
	return r.Run()
}
