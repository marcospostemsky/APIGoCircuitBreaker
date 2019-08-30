package services

import (
	"../utils"
	"../domains"
)

func GetUser(userId int) (*domains.User, *utils.ApiError){

	user := domains.User{
		ID:userId,
	}

	if err := user.Get(); err != nil {
		return nil, err
	}

	return &user, nil
}
