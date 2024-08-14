package main

import (
	"dbo-api/config"
	"dbo-api/route"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()
	config.LoadDB()

	r := gin.Default()
	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	route.AuthRouter(api)
	route.CustomerRoute(api)
	route.ProductRoute(api)
	route.OrderRoute(api)

	r.Run(fmt.Sprintf(":%v", config.ENV.PORT)) // listen and serve on 0.0.0.0:8080
}
