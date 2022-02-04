package application

import (
	"database/sql"
	ips "github.com/mborroni/dreamlab-challenge/internal/ipAddresses"
)

var (
	db      *sql.DB
	configs map[string]string
)

type Engine struct {
	AddressesService *ips.AddressesService
}

func Build() (*Engine, error) {
	buildConfig()
	buildDBConnections()

	return &Engine{
		AddressesService: buildAddressesService(),
	}, nil
}

func buildAddressesService() *ips.AddressesService {
	repository := ips.NewDBRepository(db)
	return ips.NewAddressesService(repository)
}
