package structs

import (
	"time"
)

type Orders struct {
	Order_Id      uint      `json:"order_id" gorm:"primaryKey;autoIncrement"`
	Customer_Name string    `json:"customer_name" gorm:"column:customer_name"`
	Ordered_At    time.Time `json:"ordered_at" gorm:"column:ordered_at"`
}

type Items struct {
	Item_Id     int64  `json:"item_id" gorm:"column:item_id"`
	Item_Code   string `json:"item_code" gorm:"column:item_code"`
	Description string `json:"description" gorm:"column:description"`
	Quantity    int64  `json:"quantity" gorm:"column:quantity"`
	Order_ID    int64  `json:"order_id" gorm:"column:order_id"`
}
