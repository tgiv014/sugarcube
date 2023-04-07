package app

import (
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func logger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	// raw := c.Request.URL.RawQuery

	c.Next()

	if !strings.HasPrefix(path, "/api") {
		return
	}

	log.With("latency", time.Since(start),
		"ip", c.ClientIP(),
		"method", c.Request.Method).Infof("%d %s", c.Writer.Status(), path)
}
