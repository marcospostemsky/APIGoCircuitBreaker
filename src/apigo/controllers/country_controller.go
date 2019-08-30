package controllers

import (
	"github.com/gin-gonic/gin"
	"../services"
	"net/http"
)

const (
	paramCountryID = "countryId"
)

func GetCountryFromApi (c * gin.Context) {
	countryId := c.Param(paramCountryID)

	response, err := services.GetCountry(countryId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, response)
} 