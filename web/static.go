package web

import (
	"io/fs"
	"net/http"
)

func StaticHandler() http.Handler {
	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		panic("web: static assets not embedded: " + err.Error())
	}

	return http.FileServer(http.FS(staticFS))
}
