package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/*
var assetsFS embed.FS

// serveAssets returns an http.Handler that serves embedded assets
func serveAssets() http.Handler {
	subFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(subFS))
}
