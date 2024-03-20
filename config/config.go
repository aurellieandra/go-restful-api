package config

import (
	"assignment2/structs"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit() *gorm.DB {
	host := "localhost"
	port := "5432"
	user := "postgres"
	password := "admin123"
	dbname := "orders_by"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database...")
	}

	db.AutoMigrate(&structs.Item{}, &structs.Order{})
	return db
}
