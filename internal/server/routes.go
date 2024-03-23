package server

import (
	"html/template"
	"net/http"
)

func addRoutes(mux *http.ServeMux, layout *template.Template) {
	mux.Handle("GET /", index(layout))
	mux.Handle("GET /about", about(layout))
}
