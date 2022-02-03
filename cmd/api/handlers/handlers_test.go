package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/mborroni/dreamlab-challenge/cmd/api/models"
	ips "github.com/mborroni/dreamlab-challenge/internal/ipAddresses"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockAddressesHandler(ctrl *gomock.Controller) *AddressesHandler {
	return NewAddressesHandler(NewMockservice(ctrl))
}

func TestAddressesHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := newMockAddressesHandler(ctrl)

	type fields struct {
		ip string
	}

	type want struct {
		statusCode int
		testdata   string
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{
			name: "ok",
			fields: fields{
				ip: "181.192.10.182",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					Get(gomock.Any(), fields.ip).
					Return(&ips.IP{
						ProxyType: "PUB",
						Country: ips.Country{
							Code:   "AR",
							Name:   "Argentina",
							Region: "Ciudad Autonoma de Buenos Aires",
							City:   "Buenos Aires",
						},
						ISP:    "CTL LATAM",
						Domain: "centurylink.com",
						Usage:  "ISP",
						ASN:    3356,
						AS:     "Level 3 Parent LLC",
					}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				testdata:   "./testdata/get.json",
			},
		},
		{
			name: "not found",
			fields: fields{
				ip: "181.192.10.182",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					Get(gomock.Any(), fields.ip).
					Return(nil, nil)
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "error",
			fields: fields{
				ip: "181.192.10.182",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					Get(gomock.Any(), fields.ip).
					Return(nil, errors.New("error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.fields)
			app := chi.NewRouter()

			app.Get("/v1/ips/{IP}", handler.Get)
			url := fmt.Sprintf("/v1/ips/%s", tc.fields.ip)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)

			if tc.want.testdata != "" {
				bytes, err := ioutil.ReadFile(tc.want.testdata)
				if err != nil {
					t.Fail()
				}
				var expected *models.IP
				if err := json.Unmarshal(bytes, &expected); err != nil {
					t.Fail()
				}
				var output *models.IP
				if err := json.Unmarshal(w.Body.Bytes(), &output); err != nil {
					t.Fail()
				}
				assert.EqualValues(t, expected, output)
			}
			assert.Equal(t, tc.want.statusCode, w.Code)
		})
	}
}

func TestAddressesHandler_GetTop10ISPByCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := newMockAddressesHandler(ctrl)

	type fields struct {
		country string
	}

	type want struct {
		statusCode int
		testdata   string
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{
			name: "ok",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					GetTopNISPByCountry(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]string{"Rook Media GmbH", "RapidSeedbox Ltd", "Sunrise UPC GmbH",
						"Swisscom AG", "Google LLC", "Private Layer Inc", "Datapark AG",
						"Zscaler Inc.", "Bluewin is an LIR and ISP in Switzerland.",
						"Microsoft Corporation"}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				testdata:   "./testdata/top10ISPs.json",
			},
		},
		{
			name: "error",
			fields: fields{
				country: "Argentina",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					GetTopNISPByCountry(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.fields)
			app := chi.NewRouter()
			app.Get("/v1/ips/isps/top", handler.GetTop10ISPByCountry)
			url := fmt.Sprintf("/v1/ips/isps/top?country=%s", tc.fields.country)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)
			if tc.want.testdata != "" {
				bytes, err := ioutil.ReadFile(tc.want.testdata)
				if err != nil {
					t.Fail()
				}
				var expected []string
				if err := json.Unmarshal(bytes, &expected); err != nil {
					t.Fail()
				}
				var output []string
				if err := json.Unmarshal(w.Body.Bytes(), &output); err != nil {
					t.Fail()
				}
				assert.EqualValues(t, expected, output)
			}
			assert.Equal(t, tc.want.statusCode, w.Code)
		})
	}
}

func TestAddressesHandler_GetIPQuantityByCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := newMockAddressesHandler(ctrl)

	type fields struct {
		country string
	}

	type want struct {
		statusCode int
		testdata   string
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{
			name: "ok",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					GetIPQuantityByCountry(gomock.Any(), gomock.Any()).
					Return(1398, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				testdata:   "./testdata/countryQuantity.json",
			},
		},
		{
			name: "error",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					GetIPQuantityByCountry(gomock.Any(), gomock.Any()).
					Return(0, errors.New("error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.fields)
			app := chi.NewRouter()
			app.Get("/v1/ips/quantity", handler.GetIPQuantityByCountry)
			url := fmt.Sprintf("/v1/ips/quantity?country=%s", tc.fields.country)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)
			if tc.want.testdata != "" {
				bytes, err := ioutil.ReadFile(tc.want.testdata)
				if err != nil {
					t.Fail()
				}
				var expected *models.CountryQuantity
				if err := json.Unmarshal(bytes, &expected); err != nil {
					t.Fail()
				}
				var output *models.CountryQuantity
				if err := json.Unmarshal(w.Body.Bytes(), &output); err != nil {
					t.Fail()
				}
				assert.EqualValues(t, expected, output)
			}
			assert.Equal(t, tc.want.statusCode, w.Code)
		})
	}
}

func TestAddressesHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := newMockAddressesHandler(ctrl)

	type fields struct {
		limit   int
		country string
	}

	type want struct {
		statusCode int
		testdata   string
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{
			name: "ok",
			fields: fields{
				limit:   5,
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*ips.IP{
						{
							From: 2151793288,
							Country: ips.Country{
								Name: "Switzerland",
								City: "Carouge",
							},
						},
						{
							From: 2151793430,
							Country: ips.Country{
								Name: "Switzerland",
								City: "Zurich",
							},
						},
						{
							From: 2151793432,
							Country: ips.Country{
								Name: "Switzerland",
								City: "Zurich",
							},
						},
						{
							From: 2151793588,
							Country: ips.Country{
								Name: "Switzerland",
								City: "Zurich",
							},
						},
						{
							From: 2172773163,
							Country: ips.Country{
								Name: "Switzerland",
								City: "Villigen",
							},
						},
					}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				testdata:   "./testdata/list.json",
			},
		},
		{
			name: "no content",
			fields: fields{
				limit:   10,
				country: "Japan",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil)
			},
			want: want{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name: "error",
			fields: fields{
				limit:   10,
				country: "Argentina",
			},
			expectations: func(fields fields) {
				handler.service.(*Mockservice).
					EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.fields)
			app := chi.NewRouter()
			app.Get("/v1/ips", handler.List)
			url := fmt.Sprintf("/v1/ips?country=%s&limit=%d", tc.fields.country, tc.fields.limit)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)
			if tc.want.testdata != "" {
				bytes, err := ioutil.ReadFile(tc.want.testdata)
				if err != nil {
					t.Fail()
				}
				var expected []*models.IP
				if err := json.Unmarshal(bytes, &expected); err != nil {
					t.Fail()
				}
				var output []*models.IP
				if err := json.Unmarshal(w.Body.Bytes(), &output); err != nil {
					t.Fail()
				}
				assert.EqualValues(t, expected, output)
			}
			assert.Equal(t, tc.want.statusCode, w.Code)
		})
	}
}
