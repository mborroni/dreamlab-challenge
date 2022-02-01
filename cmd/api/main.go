package main

import (
	"github.com/mborroni/dreamlab-challenge/internal/application"
)

func main() {
	server, err := NewServer()
	if err != nil {
		panic(err)
	}

	engine, _ := application.Build()

	routes(server.Router, engine)
}
