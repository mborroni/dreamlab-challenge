package ips

import (
	"context"
	"database/sql"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{
		db: db,
	}
}

func (r *DBRepository) Get(ctx context.Context, decimalIP int64) (*IP, error) {
	ip := &IP{}
	row := r.db.QueryRow("SELECT ip_from, ip_to, proxy_type, country_code, "+
		"country_name, region_name, city_name, isp, domain, usage_type, asn, 'as' "+
		"FROM ip2location_px7 WHERE ip_from <= $1 AND ip_to >= $1", decimalIP)
	err := row.Scan(&ip.From, &ip.To, &ip.ProxyType, &ip.Country.Code, &ip.Country.Name,
		&ip.Country.Region, &ip.Country.City, &ip.ISP, &ip.Domain, &ip.Usage, &ip.ASN, &ip.AS)
	return ip, err
}

func (r *DBRepository) List(ctx context.Context, limit int, filters map[string]interface{}) ([]*IP, error) {
	ips := make([]*IP, 0)
	rows, err := r.db.Query("SELECT ip_from, ip_to, country_name, city_name "+
		"FROM ip2location_px7 WHERE country_name = $1 LIMIT $2", filters["country"], limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		ip := &IP{}
		if err := rows.Scan(&ip.From, &ip.To, &ip.Country.Name, &ip.Country.City); err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func (r *DBRepository) GetIPQuantityByCountry(ctx context.Context, country string) (int, error) {
	var quantity int
	row := r.db.QueryRow("SELECT SUM(ip_to - ip_from + 1) AS quantity FROM ip2location_px7 "+
		"WHERE country_name = $1", country)
	err := row.Scan(&quantity)
	return quantity, err
}

func (r *DBRepository) GetTopNISPByCountry(ctx context.Context, top int, country string) ([]string, error) {
	isps := make([]string, 0)
	rows, err := r.db.Query("SELECT isp, count(isp) + sum(ip_to - ip_from) as total FROM ip2location_px7 "+
		"WHERE country_name = $1 GROUP BY isp ORDER BY total DESC LIMIT 10", country)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var isp string
		var total int
		if err := rows.Scan(&isp, &total); err != nil {
			return nil, err
		}
		isps = append(isps, isp)
	}
	return isps, nil
}
