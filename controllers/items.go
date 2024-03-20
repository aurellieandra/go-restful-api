package controllers

import (
	"assignment2/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetItems(c *gin.Context, db *gorm.DB) {
	var (
		items []structs.Item
		result gin.H
	)

	err := db.Find(&items).Error
	if err != nil {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": err.Error(),
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	if len(items) == 0 {
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
			"message": "Get Items Successful",
			"data": items,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}

func CreateItem(c *gin.Context, db *gorm.DB) {
	var (
		item structs.Item
		result gin.H
	)

	if err := c.Bind(&item); err != nil {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": err.Error(),
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	err := db.Create(&item).Error
	if err != nil {
		result = gin.H{
			"status": http.StatusBadRequest,
			"message": err.Error(),
			"data": nil,
		}
		c.JSON(http.StatusBadRequest, result)
		return
	} else {
		result = gin.H{
			"status": http.StatusOK,
			"message": "Create Item Successful",
			"data": item,
		}
		c.JSON(http.StatusOK, result)
	}
}

func GetItem(c *gin.Context, db *gorm.DB) {
	var (
		item structs.Item
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

	err = db.Table("items").Where("item_id = $1", id).Find(&item).Error
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
			"message": "Get Item Successful",
			"data": item,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}

func UpdateItem(c *gin.Context, db *gorm.DB) {
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

	item := db.Table("items").Where("item_id = ?", id)
	if err := item.Updates(input).Error; err != nil {
		result = gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
			"data":    nil,
		}
		c.JSON(http.StatusNotFound, result)
		return
	}

	var updatedItem structs.Item
	if err := db.First(&updatedItem, id).Error; err != nil {
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
			"message": "Update Item Successful",
			"data":    updatedItem, 
		}
		c.JSON(http.StatusOK, result)
	}
}

func DeleteItem(c *gin.Context, db *gorm.DB) {
	var (
		item structs.Item
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

	err = db.Table("items").Where("item_id = $1", id).Delete(&item).Error
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
			"message": "Delete Item Successful",
			"data": nil,
		}
		c.JSON(http.StatusOK, result)
		return
	}
}