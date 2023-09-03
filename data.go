package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define the Users table structure
type User struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Latitude  float64  `json:"latitude,omitempty"`
	Longitude float64 `json:"logitude,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Define the Products table structure
type Product struct {
	ProductID               uint     `gorm:"primary_key,em"`
	UserID                  uint     `json:"user_id,omitempty"`
	ProductName             string   `json:"product_name"`
	ProductDescription      string   `json:"product_description"`
	ProductImages           []string `json:"product_images"`
	ProductPrice            float64  `json:"product_price"`
	CompressedProductImages []string `json:"compressed_product_images"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

var dsn = "dbuser:dbpassword@tcp(localhost:3306)/apidatabase?charset=utf8mb4&parseTime=True&loc=Local"

func ProductDatabase() (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the Product table
	err = db.AutoMigrate(&Product{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func UserDatabase() (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the User and Product tables
	err = db.AutoMigrate(&User{}, &Product{})
	if err != nil {
		return nil, err
	}

	return db, nil

}

user := User{
	ID: 1,
	Name: "chandan",

}

