package web

import "net/http"

type embeddedFS struct {
	http.FileSystem
}

func (e embeddedFS) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}
