package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/index.html web/app.js
var staticFiles embed.FS

func withStatic(next http.Handler) http.Handler {
	files, _ := fs.Sub(staticFiles, "web")
	staticHandler := http.FileServer(http.FS(files))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/app.js" {
			staticHandler.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
