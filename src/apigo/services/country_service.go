package services

import (
	"../utils"
	"../domains"
)

func GetCountry(countryId string) (*domains.Country, *utils.ApiError){

	country := domains.Country{
		ID:countryId,
	}

	if err := country.Get(); err != nil {
		return nil, err
	}

	return &country, nil
}
