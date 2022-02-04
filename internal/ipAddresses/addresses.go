package ips

import (
	"context"
	"database/sql"
	"github.com/mborroni/dreamlab-challenge/internal/conversion"
)

//go:generate mockgen -source=addresses.go -destination=addresses_mock.go -package=ips

type repository interface {
	List(context.Context, int, map[string]interface{}) ([]*IP, error)
	Get(context.Context, int64) (*IP, error)
	GetIPQuantityByCountry(context.Context, string) (int, error)
	GetTop10ISPByCountry(context.Context, string) ([]string, error)
}

type AddressesService struct {
	repository repository
}

func NewAddressesService(repository repository) *AddressesService {
	return &AddressesService{
		repository: repository,
	}
}

func (s *AddressesService) List(ctx context.Context, limit int, filters map[string]interface{}) ([]*IP, error) {
	ips, err := s.repository.List(ctx, limit, filters)
	if err != nil {
		return nil, err
	}
	return split(ips), nil
}

func (s *AddressesService) Get(ctx context.Context, inputIP string) (*IP, error) {
	decimal, err := conversion.IPv4ToDecimal(inputIP)
	if err != nil {
		return nil, err
	}
	ip, err := s.repository.Get(ctx, decimal)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return ip, err
}

func (s *AddressesService) GetIPQuantityByCountry(ctx context.Context, country string) (int, error) {
	quantity, err := s.repository.GetIPQuantityByCountry(ctx, country)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return quantity, nil
}

func (s *AddressesService) GetTop10ISPByCountry(ctx context.Context, country string) ([]string, error) {
	return s.repository.GetTop10ISPByCountry(ctx, country)
}

func split(input []*IP) []*IP {
	ips := make([]*IP, 0)
	for _, ip := range input {
		for i := ip.From; i <= ip.To; i++ {
			ips = append(ips, &IP{
				From: ip.From,
				To:   ip.From,
				Country: Country{
					Name: ip.Country.Name,
					City: ip.Country.City,
				},
			})
		}
	}
	return ips
}
