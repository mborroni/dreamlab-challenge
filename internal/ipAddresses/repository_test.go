package ips

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
)

type sqlMock struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func getDB(t *testing.T) *sqlMock {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	if err != nil {
		t.Error(err.Error())
	}
	return &sqlMock{db: db, mock: mock}
}

func TestRepository_List(t *testing.T) {
	db := getDB(t)
	r := NewDBRepository(db.db)

	type fields struct {
		limit   int
		filters map[string]interface{}
	}

	type want struct {
		result []*IP
		err    error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "OK",
			fields: fields{
				limit:   2,
				filters: map[string]interface{}{"country": "Thailand"},
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT ip_from, ip_to, "+
					"country_name, city_name FROM ip2location_px7 WHERE country_name = $1 LIMIT $2")).
					WithArgs(fields.filters["country"], fields.limit).
					WillReturnRows(sqlmock.NewRows(
						[]string{"ip_from", "ip_to", "country_name", "city_name"}).
						AddRow(16778241, 16778241, "Australia", "Melbourne").
						AddRow(16778497, 16778497, "Australia", "Melbourne"),
					)
			},
			want: want{
				result: []*IP{
					{
						From: 16778241,
						To:   16778241,
						Country: Country{
							Name: "Australia",
							City: "Melbourne",
						},
					},
					{
						From: 16778497,
						To:   16778497,
						Country: Country{
							Name: "Australia",
							City: "Melbourne",
						},
					},
				},
				err: nil,
			},
		},
		{name: "error",
			fields: fields{
				limit:   1,
				filters: map[string]interface{}{"country": "Australia"},
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT ip_from, ip_to, country_name, city_name "+
					"FROM ip2location_px7 WHERE country_name = $1 LIMIT $2")).
					WithArgs(fields.filters["country"], fields.limit).
					WillReturnError(sql.ErrConnDone)
			},
			want: want{
				result: nil,
				err:    sql.ErrConnDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := r.List(context.Background(), tt.fields.limit, tt.fields.filters)

			if db.mock != nil {
				if err := db.mock.ExpectationsWereMet(); err != nil {
					t.Error(err.Error())
				}
			}
			if err != tt.want.err {
				t.Errorf("List() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.result) {
				t.Errorf("List() got = %v, want = %v", got, tt.want.result)
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	db := getDB(t)
	r := NewDBRepository(db.db)
	type fields struct {
		decimalIP int64
	}

	type want struct {
		result *IP
		err    error
	}
	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "no content",
			fields: fields{
				decimalIP: 16778497,
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT ip_from, ip_to, proxy_type, " +
					"country_code, country_name, region_name, city_name, isp, domain, usage_type, " +
					"asn, 'as' FROM ip2location_px7 WHERE ip_from <= $1 AND ip_to >= $1")).
					WithArgs(fields.decimalIP).
					WillReturnRows(sqlmock.NewRows(
						[]string{"ip_from", "ip_to", "proxy_type", "country_code",
							"country_name", "region_name", "city_name", "isp", "domain",
							"usage_type", "asn", "as"}))
			},
			want: want{
				result: &IP{},
				err:    sql.ErrNoRows,
			},
		},
		{name: "error",
			fields: fields{
				decimalIP: 16778497,
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT ip_from, ip_to, proxy_type, " +
					"country_code, country_name, region_name, city_name, isp, domain, usage_type, " +
					"asn, 'as' FROM ip2location_px7 WHERE ip_from <= $1 AND ip_to >= $1")).
					WithArgs(fields.decimalIP).
					WillReturnError(sql.ErrConnDone)
			},
			want: want{
				result: &IP{},
				err:    sql.ErrConnDone,
			},
		},
		{name: "OK",
			fields: fields{
				decimalIP: 16778497,
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT ip_from, ip_to, proxy_type, country_code, " +
					"country_name, region_name, city_name, isp, domain, usage_type, asn, 'as' " +
					"FROM ip2location_px7 WHERE ip_from <= $1 AND ip_to >= $1")).
					WithArgs(fields.decimalIP).
					WillReturnRows(sqlmock.NewRows(
						[]string{"ip_from", "ip_to", "proxy_type",
							"country_code", "country_name", "region_name",
							"city_name", "isp", "domain", "usage_type",
							"asn", "as"}).
						AddRow(16778497, 16778498, "PUB", "AU", "Australia",
							"Victoria", "Melbourne", "WirefreeBroadband Pty Ltd",
							"wirefreebroadband.com.au", "ISP", 38803,
							"WirefreeBroadband Pty Ltd"))

			},
			want: want{
				result: &IP{
					From:      16778497,
					To:        16778498,
					ProxyType: "PUB",
					Country: Country{
						Code:   "AU",
						Name:   "Australia",
						Region: "Victoria",
						City:   "Melbourne",
					},
					ISP:    "WirefreeBroadband Pty Ltd",
					Domain: "wirefreebroadband.com.au",
					Usage:  "ISP",
					ASN:    38803,
					AS:     "WirefreeBroadband Pty Ltd",
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := r.Get(context.Background(), tt.fields.decimalIP)

			if db.mock != nil {
				if err := db.mock.ExpectationsWereMet(); err != nil {
					t.Error(err.Error())
				}
			}
			if err != tt.want.err {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.result) {
				t.Errorf("Get() got = %v, want = %v", got, tt.want.result)
			}
		})
	}
}

func TestRepository_GetIPQuantityByCountry(t *testing.T) {
	db := getDB(t)
	r := NewDBRepository(db.db)
	type fields struct {
		country string
	}

	type want struct {
		result int
		err    error
	}
	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				country: "Argentina",
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(ip_to - ip_from + 1),0) AS " +
					"quantity FROM ip2location_px7 WHERE country_name = $1")).
					WithArgs(fields.country).
					WillReturnRows(sqlmock.NewRows(
						[]string{"quantity"}).AddRow(75684))
			},
			want: want{
				result: 75684,
				err:    nil,
			},
		},
		{name: "no content",
			fields: fields{
				country: "Argentina",
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(ip_to - ip_from + 1),0) AS " +
					"quantity FROM ip2location_px7 WHERE country_name = $1")).
					WithArgs(fields.country).
					WillReturnRows(sqlmock.NewRows(
						[]string{"quantity"}))
			},
			want: want{
				result: 0,
				err:    sql.ErrNoRows,
			},
		},
		{name: "error",
			fields: fields{
				country: "Argentina",
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(ip_to - ip_from + 1),0) AS " +
					"quantity FROM ip2location_px7 WHERE country_name = $1")).
					WithArgs(fields.country).
					WillReturnError(sql.ErrConnDone)
			},
			want: want{
				result: 0,
				err:    sql.ErrConnDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := r.GetIPQuantityByCountry(context.Background(), tt.fields.country)

			if db.mock != nil {
				if err := db.mock.ExpectationsWereMet(); err != nil {
					t.Error(err.Error())
				}
			}
			if err != tt.want.err {
				t.Errorf("GetIPQuantityByCountry() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.result) {
				t.Errorf("GetIPQuantityByCountry() got = %v, want = %v", got, tt.want.result)
			}
		})
	}
}

func TestRepository_GetTop10ISPByCountry(t *testing.T) {
	db := getDB(t)
	r := NewDBRepository(db.db)
	type fields struct {
		country string
		limit   int
	}

	type want struct {
		result []string
		err    error
	}
	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT isp, count(isp) + " +
					"sum(ip_to - ip_from) as total FROM ip2location_px7 " +
					"WHERE country_name = $1 GROUP BY isp ORDER BY total DESC LIMIT 10")).
					WithArgs(fields.country).
					WillReturnRows(sqlmock.NewRows(
						[]string{"isp", "total"}).
						AddRow("Rook Media GmbH", 180).
						AddRow("RapidSeedbox Ltd", 139),
					)
			},
			want: want{
				result: []string{
					"Rook Media GmbH",
					"RapidSeedbox Ltd",
				},
				err: nil,
			},
		},

		{
			name: "error",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				db.mock.ExpectQuery(regexp.QuoteMeta("SELECT isp, count(isp) + " +
					"sum(ip_to - ip_from) as total FROM ip2location_px7 " +
					"WHERE country_name = $1 GROUP BY isp ORDER BY total DESC LIMIT 10")).
					WithArgs(fields.country).
					WillReturnError(sql.ErrConnDone)
			},
			want: want{
				result: nil,
				err:    sql.ErrConnDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := r.GetTop10ISPByCountry(context.Background(), tt.fields.country)

			if db.mock != nil {
				if err := db.mock.ExpectationsWereMet(); err != nil {
					t.Error(err.Error())
				}
			}
			if err != tt.want.err {
				t.Errorf("GetTop10ISPByCountry() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.result) {
				t.Errorf("GetTop10ISPByCountry() got = %v, want = %v", got, tt.want.result)
			}
		})
	}
}
