package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:helpdesk@tcp(localhost:3306)/helpdesk?charset=utf8mb4&parseTime=True&loc=Local" // replace with your own credentials for example: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect DB")
	}

	DB = database
	fmt.Println("Database connected")
}
