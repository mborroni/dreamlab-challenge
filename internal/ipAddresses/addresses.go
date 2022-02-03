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
	GetTopNISPByCountry(context.Context, int, string) ([]string, error)
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
	return s.repository.List(ctx, limit, filters)
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

func (s *AddressesService) GetTopNISPByCountry(ctx context.Context, top int, country string) ([]string, error) {
	return s.repository.GetTopNISPByCountry(ctx, top, country)
}
