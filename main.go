package main

import (
	"helpdesk/db"
	"helpdesk/models"
	"helpdesk/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db.ConnectDB()

	// Auto migrate
	db.DB.AutoMigrate(&models.Ticket{})

	db.DB.AutoMigrate(&models.User{}, &models.Ticket{})

	routes.SetupRoutes(r)

	r.Run(":8080")
}
