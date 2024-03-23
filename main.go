package main

import (
	"fmt"
	"log/slog"

	"github.com/17xande/go-htmx-ts/internal/server"
)

func main() {
	s := server.NewServer()

	slog.Info("Starting server")

	if err := s.Http.ListenAndServe(); err != nil {
		fmt.Printf("%v", err)
	}
}
