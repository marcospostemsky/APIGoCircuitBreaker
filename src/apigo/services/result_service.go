package services

import (
	"../domains"
	"../utils"
	"sync"
)

func GetResult(userId int) (*domains.Result, *utils.ApiError){
	user := &domains.User{
		ID:userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}

	country := &domains.Country{
		ID:user.CountryID,
	}

	site := &domains.Site{
		ID:user.SiteID,
	}

	//implementar wait group
	// luego implementar ch

	err = country.Get()
	if err != nil {
		return nil, err
	}

	err = site.Get()
	if err != nil {
		return nil, err
	}


	resp := &domains.Result{
		User: user,
		Site: site,
		Country: country,
	}

	return resp,nil
}

func GetResultWg(userId int) (*domains.Result, *utils.ApiError){
	var wg sync.WaitGroup

	user := &domains.User{
		ID:userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}

	country := &domains.Country{
		ID:user.CountryID,
	}

	site := &domains.Site{
		ID:user.SiteID,
	}

	//implementar wait group
	// luego implementar ch
	ApiError := utils.ApiError{}
	wg.Add(1)
	go country.GetWg(&wg, &ApiError)
	wg.Add(1)
	go site.GetWg(&wg, &ApiError)
	wg.Wait()



	resp := &domains.Result{
		User: user,
		Site: site,
		Country: country,
	}

	return resp,nil
}

func GetResultCh(userId int) (*domains.Result, *utils.ApiError){

	results := make(chan domains.Result, 2)

	user := &domains.User{
		ID:userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}

	country := &domains.Country{
		ID:user.CountryID,
	}

	site := &domains.Site{
		ID:user.SiteID,
	}

	go country.GetCh(results)
	go site.GetCh(results)



	result := domains.Result{User:user,}
	for i := 0; i<2; i++{
		valor := <- results
		if valor.ApiError != nil {
			return nil, valor.ApiError
		}
		if valor.Country != nil {
			result.Country = valor.Country
		}
		if valor.Site != nil {
			result.Site = valor.Site
		}
	}


	return &result,nil

}