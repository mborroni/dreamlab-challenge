package application

import (
	ips "github.com/mborroni/dreamlab-challenge/internal/ipAddresses"
)

func buildAddressesService() *ips.AddressesService {
	repository := ips.NewDBRepository(db)
	return ips.NewAddressesService(repository)
}
