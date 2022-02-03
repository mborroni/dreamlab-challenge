package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mborroni/dreamlab-challenge/cmd/api/handlers"
	"github.com/mborroni/dreamlab-challenge/cmd/api/middleware"
	"github.com/mborroni/dreamlab-challenge/internal/application"
	"net/http"
	"strings"
)

func routes(router *chi.Mux, engine *application.Engine) {
	router.Get("/ping", Ping)

	handler := handlers.NewAddressesHandler(engine.AddressesService)
	router.Route("/v1/ips", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Route("/{IP}", func(r chi.Router) {
			r.With(middleware.IPValidation).Get("/", handler.Get)
		})
		r.Get("/quantity", handler.GetIPQuantityByCountry)
		r.Get("/isps/top", handler.GetTop10ISPByCountry)
	})
	printRoutes(router)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("pong")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func printRoutes(router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		panic(err)
	}
}
