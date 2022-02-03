package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mborroni/dreamlab-challenge/cmd/api/models"
	ips "github.com/mborroni/dreamlab-challenge/internal/ipAddresses"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//go:generate mockgen -source=handlers.go -destination=handlers_mock.go -package=handlers

type service interface {
	List(context.Context, int, map[string]interface{}) ([]*ips.IP, error)
	Get(context.Context, string) (*ips.IP, error)
	GetTopNISPByCountry(context.Context, int, string) ([]string, error)
	GetIPQuantityByCountry(context.Context, string) (int, error)
}

type AddressesHandler struct {
	service service
}

func NewAddressesHandler(service service) *AddressesHandler {
	return &AddressesHandler{
		service: service,
	}
}

func (h *AddressesHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := obtainLimit(r.URL)
	filters, err := obtainFilters(r.URL)
	if err != nil {
		log.WithContext(r.Context()).
			WithFields(log.Fields{"event": "error getting filters"}).
			Error(err)
		_ = RespondJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	ipAddresses, err := h.service.List(r.Context(), limit, filters)
	if err != nil {
		log.WithContext(r.Context()).
			WithFields(log.Fields{"event": "error listing"}).
			Error(err)
		_ = RespondJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(ipAddresses) == 0 {
		_ = RespondJSON(w, nil, http.StatusNoContent)
		return
	}
	_ = RespondJSON(w, models.ToIPsModel(ipAddresses), http.StatusOK)
	return
}

func (h *AddressesHandler) Get(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "IP")
	ip, err := h.service.Get(r.Context(), input)
	if err != nil {
		log.WithContext(r.Context()).
			WithFields(log.Fields{"event": "get ip address"}).
			Error(err)
		_ = RespondJSON(w, err, http.StatusInternalServerError)
		return
	}
	if ip == nil {
		_ = RespondJSON(w, nil, http.StatusNotFound)
		return
	}
	_ = RespondJSON(w, models.ToIPModel(input, ip), http.StatusOK)
	return
}

func (h *AddressesHandler) GetTop10ISPByCountry(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	isps, err := h.service.GetTopNISPByCountry(r.Context(), 10, strings.Title(country))
	if err != nil {
		log.WithContext(r.Context()).
			WithFields(log.Fields{"event": "get top 10 ISPs"}).
			Error(err)
		_ = RespondJSON(w, err, http.StatusInternalServerError)
		return
	}
	_ = RespondJSON(w, isps, http.StatusOK)
	return
}

func (h *AddressesHandler) GetIPQuantityByCountry(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	quantity, err := h.service.GetIPQuantityByCountry(r.Context(), strings.Title(country))
	if err != nil {
		log.WithContext(r.Context()).
			WithFields(log.Fields{"event": "get ip quantity by country"}).
			Error(err)
		_ = RespondJSON(w, err, http.StatusInternalServerError)
		return
	}
	_ = RespondJSON(w, models.ToCountryQuantityModel(country, quantity), http.StatusOK)
	return
}

func RespondJSON(w http.ResponseWriter, v interface{}, code int) error {
	if code == http.StatusNoContent || v == nil {
		w.WriteHeader(code)
		return nil
	}
	var jsonData []byte
	var err error
	switch v := v.(type) {
	case []byte:
		jsonData = v
	case io.Reader:
		jsonData, err = ioutil.ReadAll(v)
	default:
		jsonData, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(jsonData); err != nil {
		return err
	}
	return nil
}

func obtainLimit(u *url.URL) int {
	limit := 100
	if l := u.Query().Get("limit"); l != "" {
		if l, err := strconv.Atoi(l); err == nil {
			limit = l
		}
	}
	return limit
}

func obtainFilters(u *url.URL) (map[string]interface{}, error) {
	filters := make(map[string]interface{})
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	if country := q.Get("country"); country != "" {
		filters["country"] = strings.Title(country)
	}
	return filters, nil
}
