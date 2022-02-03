package models

type CountryQuantity struct {
	Country  string `json:"country"`
	Quantity int    `json:"quantity"`
}

func ToCountryQuantityModel(country string, quantity int) *CountryQuantity {
	return &CountryQuantity{
		Country:  country,
		Quantity: quantity,
	}
}
