package main

import (
	"github.com/gin-gonic/gin"
	"../controllers"
	"log"
	"time"
	"../circuitbreaker"
)

const (
	port = ":9090"
)

var (
	router = gin.Default()
)


func main() {
	timeout := 20
	controllers.Cb = circuitbreaker.NewCircuitBreaker("cb",3,time.Second*time.Duration(timeout),0)

	go controllers.Cb.SetState()


	router.GET("/users/:userId", controllers.GetUserFromApi)
	router.GET("/countries/:countryId", controllers.GetCountryFromApi)
	router.GET("/sites/:siteId", controllers.GetSiteFromApi)
	router.GET("/results/:userId", controllers.GetResultFromApi)
	router.GET("/resultswg/:userId", controllers.GetResultFromApiWg)
	router.GET("/resultsch/:userId", controllers.GetResultFromApiCh)
	log.Fatal(router.Run(port))
}
