package controllers

import (
	"github.com/gin-gonic/gin"
	"../services"
	"net/http"
	"strconv"
	"../utils"
	"fmt"
	"../circuitbreaker"
)

var (
	Cb *circuitbreaker.CircuitBreaker
)

func GetResultFromApi (c * gin.Context) {
	resultId := c.Param(paramUserID)

	id, error := strconv.Atoi(resultId)
	if error != nil {
		apiErr := utils.ApiError{
			Message:error.Error(),
			Status:http.StatusBadRequest,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}


	if Cb.State != circuitbreaker.STATE_CLOSE {

		c.JSON(580, &utils.ApiError{Status:580,Message:"Está bloqueado amigo"})
		return

	}

	//BORRAR
	Cb.ErrCounts = 0
	for {
		response, err := services.GetResult(id)
		if err == nil || err.Status != 500 {
			if err != nil {
				c.JSON(err.Status, err)
				return
			}
			c.JSON(http.StatusOK, response)
			break
		}


		if Cb.ErrCounts >= 3{
			Cb.State = circuitbreaker.STATE_OPEN
			Cb.ChainStatus <- "Status"
			fmt.Println("Le mando a canal el true")
			c.JSON(580, &utils.ApiError{Status:580,Message:"Está bloqueado amigo"})
			return
		}

		fmt.Println(Cb.ErrCounts)
		Cb.ErrCounts++
	}

}

func GetResultFromApiWg (c * gin.Context) {
	resultId := c.Param(paramUserID)

	id, error := strconv.Atoi(resultId)
	if error != nil {
		apiErr := utils.ApiError{
			Message:error.Error(),
			Status:http.StatusBadRequest,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}

	response, err := services.GetResultWg(id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func GetResultFromApiCh (c * gin.Context) {
	resultId := c.Param(paramUserID)

	id, error := strconv.Atoi(resultId)
	if error != nil {
		apiErr := utils.ApiError{
			Message:error.Error(),
			Status:http.StatusBadRequest,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}



	for {
		if Cb.State != circuitbreaker.STATE_CLOSE {

			c.JSON(580, &utils.ApiError{Status:580,Message:"Está bloqueado amigo"})
			return

		}

		response, err := services.GetResultCh(id)
		if err == nil || err.Status != 500 {
			Cb.ChainCount <- "OK"
			if err != nil {
				c.JSON(err.Status, err)
				return
			}
			c.JSON(http.StatusOK, response)
			return
		}


		Cb.ChainCount <- "ERROR"
	}


	/*response, err := services.GetResultCh(id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, response)*/
}

