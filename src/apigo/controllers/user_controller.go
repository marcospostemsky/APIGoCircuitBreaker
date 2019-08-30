package controllers

import (
	"github.com/gin-gonic/gin"
	"../services"
	"net/http"
	"strconv"
	"../utils"
)

const (
	paramUserID = "userId"
)

func GetUserFromApi (c * gin.Context) {
	userId := c.Param(paramUserID)

	id, error := strconv.Atoi(userId)
	if error != nil {
		apiErr := utils.ApiError{
			Message:error.Error(),
			Status:http.StatusBadRequest,
		}
		c.JSON(apiErr.Status, apiErr)
		return
	}

	response, err := services.GetUser(id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, response)
} 