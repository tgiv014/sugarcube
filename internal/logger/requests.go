package logger

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/charmbracelet/log"
)

var RequestLogger gin.HandlerFunc = func(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	// raw := c.Request.URL.RawQuery

	c.Next()

	// Not logging forwarded requests and 404s for the frontend
	if !strings.HasPrefix(path, "/api") {
		return
	}

	log.With("latency", time.Since(start),
		"ip", c.ClientIP(),
		"method", c.Request.Method).Infof("%d %s", c.Writer.Status(), path)
}
