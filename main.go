package main

import (
	"assignment2/config"
	"assignment2/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db := config.DBInit()

	router := gin.Default()

	router.POST("/orders", func(c *gin.Context) {
		controllers.CreateOrder(c, db)
	})

	router.GET("/orders", func(c *gin.Context) {
		controllers.GetOrders(c, db)
	})

	router.Run(":3000")
}
