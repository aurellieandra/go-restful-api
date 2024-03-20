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
		items  []structs.Item
		result gin.H
	)

	err := db.Find(&orders).Error
	if err != nil {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	var ordersData []map[string]interface{}

	for _, order := range orders {
		err = db.Find(&items).Error
		if err != nil {
			result = gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
				"data":    nil,
			}
			c.JSON(http.StatusBadRequest, result)
			return
		}

		orderData := map[string]interface{}{
			"order_id":      order.Order_Id,
			"customer_name": order.Customer_Name,
			"ordered_at":    order.Ordered_At,
			"items":         items,
		}
		ordersData = append(ordersData, orderData)
	}

	if len(orders) == 0 {
		result = gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
			"data":    nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	}

	result = gin.H{
		"status":  http.StatusOK,
		"message": "Get Orders Successful",
		"data":    ordersData,
	}
	c.JSON(http.StatusOK, result)
}

func CreateOrder(c *gin.Context, db *gorm.DB) {
	var (
		input  map[string]interface{}
		result gin.H
	)

	if err := c.BindJSON(&input); err != nil {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to bind input",
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	var jsonData = make(map[string]interface{})

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
		Ordered_At:    orderedAt,
	}

	// add order data
	if err := db.Create(&order).Error; err != nil {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	// add items data
	var items []structs.Item
	if inputItems, ok := input["items"].([]interface{}); ok {
		for _, item := range inputItems {
			if itemMap, ok := item.(map[string]interface{}); ok {
				itemCode, _ := itemMap["item_code"].(string)
				description, _ := itemMap["description"].(string)
				quantity, _ := strconv.ParseInt(itemMap["quantity"].(string), 10, 64)

				newItems := structs.Item{
					Item_Code:   itemCode,
					Description: description,
					Quantity:    quantity,
					Order_Id:    uint(order.Order_Id),
				}

				if err := db.Create(&newItems).Error; err != nil {
					result = gin.H{
						"status":  http.StatusBadRequest,
						"message": err.Error(),
						"data":    nil,
					}
					c.JSON(http.StatusBadRequest, result)
					return
				}

				items = append(items, newItems)
			}
		}
		jsonData["items"] = items
	}

	jsonData["order"] = order

	result = gin.H{
		"status":  http.StatusOK,
		"message": "Create Order Successful",
		"data":    jsonData,
	}
	c.JSON(http.StatusOK, result)
}

func GetOrder(c *gin.Context, db *gorm.DB) {
	var (
		order  structs.Order
		items  []structs.Item
		result gin.H
	)

	var jsonData = make(map[string]interface{})

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": "ID is Invalid",
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	err = db.Table("orders").Where("order_id = $1", id).Find(&order).Error
	if err != nil {
		result = gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
			"data":    nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	}

	err = db.Table("items").Where("order_id = ?", id).Find(&items).Error
	if err != nil {
		result = gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
			"data":    nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	}

	jsonData["order_id"] = order.Order_Id
	jsonData["customer_name"] = order.Customer_Name
	jsonData["ordered_at"] = order.Ordered_At
	jsonData["items"] = items

	result = gin.H{
		"status":  http.StatusOK,
		"message": "Get Order Successful",
		"data":    jsonData,
	}
	c.JSON(http.StatusOK, result)
}

func UpdateOrder(c *gin.Context, db *gorm.DB) {
	var (
		input  map[string]interface{}
		result gin.H
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": "ID is Invalid",
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	order := db.Table("orders").Where("order_id = ?", id)
	orderUpdateFields := map[string]interface{}{
		"order_id":   id,
		"ordered_at": input["ordered_at"],
	}
	if err := order.Updates(orderUpdateFields).Error; err != nil {
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
	}

	var updatedItems []structs.Item
	if inputItems, ok := input["items"].([]interface{}); ok {
		for _, item := range inputItems {
			if itemMap, ok := item.(map[string]interface{}); ok {
				itemID, _ := strconv.Atoi(itemMap["item_id"].(string))
				itemCode, _ := itemMap["item_code"].(string)
				description, _ := itemMap["description"].(string)
				quantity, _ := strconv.ParseInt(itemMap["quantity"].(string), 10, 64)

				if err := db.Model(&structs.Item{}).Where("item_id = ?", itemID).Updates(map[string]interface{}{
					"item_code":   itemCode,
					"description": description,
					"quantity":    quantity,
				}).Error; err != nil {
					result = gin.H{
						"status":  http.StatusBadRequest,
						"message": err.Error(),
						"data":    nil,
					}
					c.JSON(http.StatusBadRequest, result)
					return
				}

				// Get the updated item
				var item structs.Item
				if err := db.First(&item, itemID).Error; err != nil {
					result = gin.H{
						"status":  http.StatusNotFound,
						"message": "Data Not Found",
						"data":    nil,
					}
					c.JSON(http.StatusNotFound, result)
					return
				}
				updatedItems = append(updatedItems, item)

			}
		}
	}

	responseData := map[string]interface{}{
		"order": updatedOrder,
		"items": updatedItems,
	}

	result = gin.H{
		"status":  http.StatusOK,
		"message": "Update Order Successful",
		"data":    responseData,
	}
	c.JSON(http.StatusOK, result)
}

func DeleteOrder(c *gin.Context, db *gorm.DB) {
	var (
		order  structs.Order
		result gin.H
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		result = gin.H{
			"status":  http.StatusBadRequest,
			"message": "ID is Invalid",
			"data":    nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	err = db.Table("orders").Where("order_id = $1", id).Delete(&order).Error
	if err != nil {
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
			"message": "Delete Order Successful",
			"data":    nil,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}
