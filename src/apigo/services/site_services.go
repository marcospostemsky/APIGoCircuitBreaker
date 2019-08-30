package services

import (
	"../utils"
	"../domains"
)

func GetSite(siteId string) (*domains.Site, *utils.ApiError){

	site := domains.Site{
		ID:siteId,
	}

	if err := site.Get(); err != nil {
		return nil, err
	}

	return &site, nil
}
