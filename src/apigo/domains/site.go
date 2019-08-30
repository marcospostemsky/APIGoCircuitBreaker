package domains

import (
	"../utils"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"sync"
	"fmt"
)

type Site struct {
	ID                 string       `json:"id"`
	Name               string       `json:"name"`
	CountryID          string       `json:"country_id"`
	SaleFeesMode       string       `json:"sale_fees_mode"`
	MercadopagoVersion int          `json:"mercadopago_version"`
	DefaultCurrencyID  string       `json:"default_currency_id"`
	ImmediatePayment   string       `json:"immediate_payment"`
	PaymentMethodIds   []string     `json:"payment_method_ids"`
	Settings           Settings     `json:"settings"`
	Currencies         []Currencies `json:"currencies"`
	Categories         []Categories `json:"categories"`
}
type Rules struct {
	EnabledTaxpayerTypes []string `json:"enabled_taxpayer_types"`
	BeginsWith           string   `json:"begins_with"`
	Type                 string   `json:"type"`
	MinLength            int      `json:"min_length"`
	MaxLength            int      `json:"max_length"`
}
type IdentificationTypesRules struct {
	IdentificationType string  `json:"identification_type"`
	Rules              []Rules `json:"rules"`
}
type Settings struct {
	IdentificationTypes      []string                   `json:"identification_types"`
	TaxpayerTypes            []string                   `json:"taxpayer_types"`
	IdentificationTypesRules []IdentificationTypesRules `json:"identification_types_rules"`
}
type Currencies struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}
type Categories struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (site *Site) Get() *utils.ApiError {
	if site.ID == "" {
		return &utils.ApiError{
			Message: "Site ID is empty.",
			Status:  http.StatusBadRequest,
		}
	}
	url := utils.UrlSite + site.ID

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

	if err := json.Unmarshal(data, &site); err != nil { //crea una variable que solo vivi en el if
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (site *Site) GetWg(wg *sync.WaitGroup, ApiError *utils.ApiError)  {
	if site.ID == "" {
		ApiError = &utils.ApiError{
			Message: "Site ID is empty.",
			Status:  http.StatusBadRequest,
		}
	}
	url := utils.UrlSite + site.ID

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

	if err := json.Unmarshal(data, &site); err != nil { //crea una variable que solo vivi en el if
		ApiError = &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	wg.Done()
}

func (site *Site) GetCh(result chan Result) {
	if site.ID == "" {
		ApiError := &utils.ApiError{
			Message: "Site ID is empty.",
			Status:  http.StatusBadRequest,
		}
		result <- Result{ApiError:ApiError}
	}
	url := utils.UrlSite + site.ID

	response, err := http.Get(url)
	if err != nil {
		ApiError := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		fmt.Println(ApiError)
		result <- Result{ApiError:ApiError}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ApiError := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		fmt.Println(ApiError)
		result <- Result{ApiError:ApiError}
	}

	if err := json.Unmarshal(data, &site); err != nil { //crea una variable que solo vivi en el if
		ApiError := &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		fmt.Println(ApiError)
		result <- Result{ApiError:ApiError}
	}

	result <- Result{Site:site}
}