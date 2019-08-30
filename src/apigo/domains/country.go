package domains

import (
	"../utils"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"sync"
)

type Country struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Locale             string         `json:"locale"`
	CurrencyID         string         `json:"currency_id"`
	DecimalSeparator   string         `json:"decimal_separator"`
	ThousandsSeparator string         `json:"thousands_separator"`
	TimeZone           string         `json:"time_zone"`
	GeoInformation     GeoInformation `json:"geo_information"`
	States             []States       `json:"states"`
}
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type GeoInformation struct {
	Location Location `json:"location"`
}
type States struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}


func (country *Country) Get() *utils.ApiError {
	if country.ID == "" {
		return &utils.ApiError{
			Message: "Country ID is empty.",
			Status:  http.StatusBadRequest,
		}
	}
	url := utils.UrlCountry + country.ID

	response, err := http.Get(url)
	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if err := json.Unmarshal(data, &country); err != nil { //crea una variable que solo vivi en el if
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (country *Country) GetWg(wg *sync.WaitGroup, ApiError *utils.ApiError) {
	if country.ID == "" {
		ApiError = &utils.ApiError{
			Message: "Country ID is empty.",
			Status:  http.StatusBadRequest,
		}
	}
	url := utils.UrlCountry + country.ID

	response, err := http.Get(url)
	if err != nil {
		ApiError = &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ApiError = &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if err := json.Unmarshal(data, &country); err != nil { //crea una variable que solo vivi en el if
		ApiError = &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	wg.Done()
}

func (country *Country) GetCh(result chan Result) {
	if country.ID == "" {
		ApiError := &utils.ApiError{
			Message: "Country ID is empty.",
			Status:  http.StatusBadRequest,
		}
		result <- Result{ApiError:ApiError}
	}
	url := utils.UrlCountry + country.ID

	response, err := http.Get(url)
	if err != nil {
		ApiError := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		result <- Result{ApiError:ApiError}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ApiError := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		result <- Result{ApiError:ApiError}
	}

	if err := json.Unmarshal(data, &country); err != nil { //crea una variable que solo vivi en el if
		ApiError :=  &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		result <- Result{ApiError:ApiError}
	}

	result <- Result{Country:country}
}