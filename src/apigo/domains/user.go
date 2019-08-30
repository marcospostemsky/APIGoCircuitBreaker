package domains

import (
	"../utils"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type User struct {
	ID               int              `json:"id"`
	Nickname         string           `json:"nickname"`
	RegistrationDate string           `json:"registration_date"`
	CountryID        string           `json:"country_id"`
	Address          Address          `json:"address"`
	UserType         string           `json:"user_type"`
	Tags             []string         `json:"tags"`
	Logo             interface{}      `json:"logo"`
	Points           int              `json:"points"`
	SiteID           string           `json:"site_id"`
	Permalink        string           `json:"permalink"`
	SellerReputation SellerReputation `json:"seller_reputation"`
	BuyerReputation  BuyerReputation  `json:"buyer_reputation"`
	Status           Status           `json:"status"`
}
type Address struct {
	City  string `json:"city"`
	State string `json:"state"`
}
type Ratings struct {
	Negative int `json:"negative"`
	Neutral  int `json:"neutral"`
	Positive int `json:"positive"`
}
type Transactions struct {
	Canceled  int     `json:"canceled"`
	Completed int     `json:"completed"`
	Period    string  `json:"period"`
	Ratings   Ratings `json:"ratings"`
	Total     int     `json:"total"`
}
type SellerReputation struct {
	LevelID           interface{}  `json:"level_id"`
	PowerSellerStatus interface{}  `json:"power_seller_status"`
	Transactions      Transactions `json:"transactions"`
}
type BuyerReputation struct {
	Tags []interface{} `json:"tags"`
}
type Status struct {
	SiteStatus string `json:"site_status"`
}

func (user *User) Get() *utils.ApiError {
	if user.ID == 0 {
		return &utils.ApiError{
			Message: "User ID is empty.",
			Status:  http.StatusBadRequest,
		}
	}
	url := utils.UrlUser + strconv.Itoa(user.ID)


	response, err := http.Get(url)
	if err != nil{
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

	if err := json.Unmarshal(data, &user); err != nil { //crea una variable que solo vivi en el if
		return &utils.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}