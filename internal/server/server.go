package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

//go:embed templates/*
var embTemplates embed.FS

//go:embed assets/*
var embAssets embed.FS

func NewServer(ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	layout := template.Must(template.ParseFS(embTemplates, "templates/layouts/layout.gohtml"))
	addRoutes(mux, layout)
	var handler http.Handler = mux

	// Middleware.
	handler = logger(handler)

	return handler

}

// logger middleware to log requests.
func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	return v, nil
}

func index(layout *template.Template) http.Handler {
	layout = template.Must(layout.Clone())
	t := template.Must(layout.ParseFS(embTemplates, "templates/index.gohtml"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if err := t.ExecuteTemplate(w, "layout.gohtml", nil); err != nil {
			log.Printf("%v\n", err)
		}
	})
}

func about(layout *template.Template) http.Handler {
	layout = template.Must(layout.Clone())
	t := template.Must(layout.ParseFS(embTemplates, "templates/about.gohtml"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "layout.gohtml", r); err != nil {
			slog.Error(err.Error())
		}
	})
}
