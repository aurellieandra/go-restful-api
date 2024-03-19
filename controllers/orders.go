package controllers

import (
	"assignment2/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(c *gin.Context, db *gorm.DB) {
	var (
		orders structs.Orders
		result gin.H
	)

	customer_name := c.PostForm("customer_name")

	orders.Customer_Name = customer_name
	orders.Ordered_At = time.Now()

	err := db.Create(&orders).Error
	if err != nil {
		result = gin.H{
			"error": err.Error(),
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}
	result = gin.H{
		"result": orders,
	}
	c.JSON(http.StatusOK, result)
}

func GetOrders(c *gin.Context, db *gorm.DB) {
	var (
		orders []structs.Orders
		result gin.H
	)

	err := db.Find(&orders).Error
	if err != nil {
		result = gin.H{
			"error": err.Error(),
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	if len(orders) == 0 {
		result = gin.H{
			"result": "The order table is empty",
		}
		c.JSON(http.StatusNotFound, result)
		return
	} else {
		result = gin.H{
			"result": orders,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}
