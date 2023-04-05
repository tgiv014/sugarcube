// package web provides the web frontend
package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed dist
var embedFS embed.FS

func Static() (gin.HandlerFunc, error) {
	// Prod mode
	dist, err := fs.Sub(embedFS, "dist")
	if err != nil {
		return nil, err
	}

	wwwroot := embeddedFS{http.FS(dist)}
	return static.Serve("/", wwwroot), nil
}
