package controllers

import (
	"github.com/gin-gonic/gin"
	"../services"
	"net/http"
	"strconv"
	"../utils"
	"fmt"
	"../circuitbreaker"
	"../domains"
	"time"
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

	for {

		if Cb.State != circuitbreaker.STATE_CLOSE {

			c.JSON(http.StatusServiceUnavailable, &utils.ApiError{Status:http.StatusServiceUnavailable ,Message:"Status service unavailable"})
			return

		}

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


	//Se implementa en el for, si el Circuit Breaker (CB) está cerrado para poder enviar request a la API de MELI. Se puede sumar
	// al contador de errores del CB si se activa el timeout (3 segundos) o la API de MELI devuelve error 500.
	// En un pedido correcto se devuelve al canal un "OK", para resetear el contador de eventos.
	for {

		if Cb.State != circuitbreaker.STATE_CLOSE {
			c.JSON(http.StatusServiceUnavailable, &utils.ApiError{Status:http.StatusServiceUnavailable ,Message:"Status service unavailable"})
			return


		}

		// se crean nuevos canales cada bucle, para que las go routines lanzadas no tengan conflictos.
		c1 := make(chan domains.Result, 1)
		c2 := make(chan utils.ApiError, 1)

		go func() {
			response, err := services.GetResultCh(id)
			if err == nil || err.Status != 500 {
				if err != nil {
					c2 <- *err
					return
				}
				c1 <- *response
				return
			}
		}()

		// 	Espera canal c1 de respuesta o c2 de error.
		select {
			case <-time.After(time.Second * 3):
				break
			case res := <-c1:
				c.JSON(http.StatusOK, res)
				Cb.Counter(circuitbreaker.OK)
				return
			case err := <- c2:
				c.JSON(err.Status,err)
				Cb.Counter(circuitbreaker.OK)
				return
		}

		Cb.Counter(circuitbreaker.ERROR)
	}
}

