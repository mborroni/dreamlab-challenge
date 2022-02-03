package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/mborroni/dreamlab-challenge/internal/conversion"
	"net/http"
)

func IPValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := chi.URLParam(r, "IP")
		if _, err := conversion.IPv4ToDecimal(ip); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		next.ServeHTTP(w, r)
	})
}
