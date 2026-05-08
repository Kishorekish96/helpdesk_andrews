package routes

import (
	_ "helpdesk/docs"
	"helpdesk/handlers"

	_ "helpdesk/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	r.Static("/static", "./static")
	// we should set allow orgin cross domain for frontend to access backend api

	r.Use(func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// IMPORTANT: handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")

	api.POST("/users", handlers.CreateUser)
	api.GET("/users", handlers.GetUsers)
	api.GET("/admins", handlers.GetAdmins)

	api.POST("/tickets", handlers.CreateTicket)
	api.GET("/tickets", handlers.GetTickets)
	api.GET("/tickets/:id", handlers.GetTicketByID)
	api.PUT("/tickets/:id", handlers.UpdateTicket)
	api.DELETE("/tickets/:id", handlers.DeleteTicket)
	api.PUT("/tickets/:id/assign", handlers.AssignTicket) // ✅ add this
}
