package main

import (
	"github.com/go-chi/chi/v5"
	"os"
)

type Server struct {
	Environment string
	Router      *chi.Mux
}

func NewServer() (*Server, error) {
	environment := "local"
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		environment = env
	}
	router := chi.NewRouter()

	server := &Server{
		Environment: environment,
		Router:      router,
	}
	return server, nil
}
