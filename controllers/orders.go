package controllers

import (
	"assignment2/structs"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrders(c *gin.Context, db *gorm.DB) {
	var (
		orders []structs.Order
		result gin.H
	)

	err := db.Find(&orders).Error
	if err != nil {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": err.Error(),
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	if len(orders) == 0 {
		result = gin.H{
			"status": http.StatusNotFound,
			"message": "Data Not Found",
			"data": nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	} else {
		result = gin.H{
			"status": http.StatusOK,
			"message": "Get Orders Successful",
			"data": orders,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}

func CreateOrder(c *gin.Context, db *gorm.DB) {
	var (
		input map[string]interface{}
		result gin.H
	)

	if err := c.BindJSON(&input); err != nil {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": "Failed to bind input",
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	customerNameInterface, customerNameExists := input["customer_name"]
	orderedAtInterface, orderedAtExists := input["ordered_at"]
	
	if !customerNameExists || !orderedAtExists {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": "Missing required fields",
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}
	var customerName string
	if customerNameInterface != nil {
		customerName, _ = customerNameInterface.(string)
	}
	var orderedAt time.Time
	if orderedAtInterface != nil {
		orderedAt, _ = time.Parse("2006-01-02 03:04:05", orderedAtInterface.(string))
		// orderedAt = theTime.Format(time.RFC3339Nano)
	}

	order := structs.Order{
		Customer_Name: customerName,
		Ordered_At: orderedAt,
	}
	
	// add order data
	if err := db.Create(&order).Error; err != nil {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": err.Error(),
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	// add items data
	if inputItems, ok := input["items"].([]interface{}); ok {
		for _, item := range inputItems {
			if itemMap, ok := item.(map[string]interface{}); ok {
				itemCode, _ := itemMap["item_code"].(string)
				description, _ := itemMap["description"].(string)
				quantity, _ := strconv.ParseInt(itemMap["quantity"].(string), 10, 64)

				items := structs.Item {
					Item_Code: itemCode,
					Description: description,
					Quantity: quantity,
					Order_Id: uint(order.Order_Id),
				}


				if err := db.Create(&items).Error; err != nil {
					result = gin.H{
						"status": http.StatusBadRequest,
						"message": err.Error(),
						"data": nil,
					}
					c.JSON(http.StatusBadRequest, result)
					return
				}
			}
		}
	}

	result = gin.H{
		"status": http.StatusOK,
		"message": "Create Order Successful",
		"data": order,
	}
	c.JSON(http.StatusOK, result)
}

func GetOrder(c *gin.Context, db *gorm.DB) {
	var (
		order structs.Order
		result gin.H
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": "ID is Invalid",
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	err = db.Table("orders").Where("order_id = $1", id).Find(&order).Error
	if err != nil {
		result = gin.H{
			"status": http.StatusNotFound,
			"message": "Data Not Found",
			"data": nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	} else {
		result = gin.H{
			"status": http.StatusOK,
			"message": "Get Order Successful",
			"data": order,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}

func UpdateOrder(c *gin.Context, db *gorm.DB) {
	var (
		input map[string]interface{}
		result gin.H
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": "ID is Invalid",
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": err.Error(),
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	order := db.Table("orders").Where("order_id = ?", id)
	if err := order.Updates(input).Error; err != nil {
		result = gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
			"data":    nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	}

	var updatedOrder structs.Order
	if err := db.First(&updatedOrder, id).Error; err != nil {
		result = gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
			"data":    nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	} else {
		result = gin.H{
			"status":  http.StatusOK,
			"message": "Update Order Successful",
			"data":    updatedOrder, 
		}
		c.JSON(http.StatusOK, result)
	}
}

func DeleteOrder(c *gin.Context, db *gorm.DB) {
	var (
		order structs.Order
		result gin.H
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": "ID is Invalid",
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	err = db.Table("orders").Where("order_id = $1", id).Delete(&order).Error
	if err != nil {
		result = gin.H{
			"status": http.StatusNotFound,
			"message": "Data Not Found",
			"data": nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	} else {
		result = gin.H{
			"status": http.StatusOK,
			"message": "Delete Order Successful",
			"data": nil,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}