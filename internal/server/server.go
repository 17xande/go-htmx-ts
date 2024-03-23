package server

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
)

//go:embed templates/*
var embTemplates embed.FS

//go:embed assets/*
var embAssets embed.FS

type Server struct {
	Http *http.Server
}

func NewServer() *Server {
	mux := http.ServeMux{}
	s := &http.Server{
		Addr:    ":8080",
		Handler: &mux,
	}

	serv := &Server{
		Http: s,
	}

	// TODO: handle request not found.

	mux.HandleFunc("GET /", serv.index)
	mux.HandleFunc("GET /about", serv.about)
	return serv
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	// TODO: rather call Must() in the server constructor to prevent panics at runtime.
	layout := template.Must(template.ParseFS(embTemplates, "templates/layouts/layout.gohtml"))
	t := template.Must(layout.ParseFS(embTemplates, "templates/index.gohtml"))
	if err := t.ExecuteTemplate(w, "layout.gohtml", nil); err != nil {
		slog.Error(err.Error())
	}
}

func (s *Server) about(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFS(embTemplates, "templates/*.gohtml"))
	t.ExecuteTemplate(w, "about.gohtml", r)
}
