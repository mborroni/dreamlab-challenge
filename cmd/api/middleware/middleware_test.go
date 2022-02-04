package middleware

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_ValidIP(t *testing.T) {

	type fields struct {
		ip string
	}

	type want struct {
		statusCode int
		testdata   string
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "ok",
			fields: fields{
				ip: "181.100.10.182",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "invalid IP",
			fields: fields{
				ip: "181.ab.10.182",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := chi.NewRouter()
			handlerMock := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}
			app.With(IPValidation).Get("/v1/ips/{IP}", handlerMock)
			url := fmt.Sprintf("/v1/ips/%s", tc.fields.ip)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
		})
	}
}
