package main

import (
	"assignment2/config"
	"assignment2/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db := config.DBInit()
	engine := gin.New()

	orders := engine.Group("/orders")

	orders.GET("/", func(c *gin.Context) {
		controllers.GetOrders(c, db)
	})
	orders.POST("/", func(c *gin.Context) {
		controllers.CreateOrder(c, db)
	})
	orders.GET("/:id", func(c *gin.Context) {
		controllers.GetOrder(c, db)
	})
	orders.PUT("/:id", func(c *gin.Context) {
		controllers.UpdateOrder(c, db)
	})
	orders.DELETE("/:id", func(c *gin.Context) {
		controllers.DeleteOrder(c, db)
	})

	engine.Run(":3000")
}
