package main

import (
	"github.com/mborroni/dreamlab-challenge/internal/application"
	"net/http"
)

func main() {
	server, err := newServer()
	if err != nil {
		panic(err)
	}
	engine, err := application.Build()
	if err != nil {
		panic(err)
	}
	routes(server, engine)
	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}
