package handlers

import (
	"fmt"
	"net/http"

	"helpdesk/db"
	"helpdesk/models"

	"github.com/gin-gonic/gin"
)

// create user

// @Summary Create user
// @Description Create a new user (admin or student)
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User"
// @Success 200 {object} models.ListUser
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var user models.CreateUserRequest
	var userModel models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userModel.Name = user.Name
	userModel.Email = user.Email
	userModel.Role = user.Role
	userModel.Department = user.Department

	if user.Role != "admin" && user.Role != "student" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Validation
	if user.Role == "admin" && user.Department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department is required for admin"})
		return
	}

	result := db.DB.Create(&userModel)

	if result.Error != nil {
		fmt.Println("DB ERROR:", result.Error)
	}

	c.JSON(http.StatusOK, user)
}

// get all users

// @Summary Get all users
// @Description Get a list of all users
// @Accept json
// @Produce json
// @Success 200 {array} models.ListUser
// @Router /api/users [get]
func GetUsers(c *gin.Context) {
	var users []models.User

	db.DB.Find(&users)

	c.JSON(http.StatusOK, users)
}

// Get Only Admins

// @Summary Get admins
// @Description Get a list of all admins
// @Accept json
// @Produce json
// @Success 200 {array} models.ListUser
// @Router /api/admins [get]
func GetAdmins(c *gin.Context) {
	var admins []models.User

	db.DB.Where("role = ?", "admin").Find(&admins)

	c.JSON(http.StatusOK, admins)
}

// Create Ticket

// @Summary Create ticket
// @Description Create a new ticket
// @Accept json
// @Produce json
// @Param ticket body models.CreateTicketRequest true "Ticket"
// @Success 200 {object} models.CreateTicketRequest
// @Router /api/tickets [post]
func CreateTicket(c *gin.Context) {
	var t models.CreateTicketRequest

	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ticket models.Ticket = models.Ticket{
		Title:       t.Title,
		Description: t.Description,
		Category:    t.Category,
		StudentID:   t.StudentID,
		Status:      t.Status,
		AdminID:     t.AdminID,
	}

	// Validate required fields
	if ticket.Title == "" || ticket.Category == "" || ticket.StudentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Check if student exists
	var student models.User
	if err := db.DB.First(&student, ticket.StudentID).Error; err != nil || student.Role != "student" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student"})
		return
	}

	// Assign admin based on category (simple logic)
	var admin models.User
	if err := db.DB.Where("role = ? AND department = ?", "admin", ticket.Category).First(&admin).Error; err == nil {
		ticket.AdminID = admin.ID
		ticket.Status = "in-progress"
	} else {
		ticket.Status = "open" // fallback
	}

	db.DB.Create(&ticket)

	c.JSON(http.StatusOK, ticket)
}

// Get All Tickets

// @Summary Get all tickets
// @Description Get a list of all tickets
// @Accept json
// @Produce json
// @Success 200 {array} models.TicketResponse
// @Router /api/tickets [get]
func GetTickets(c *gin.Context) {
	var tickets []models.Ticket
	//db.DB.Find(&tickets)

	//db.DB.Preload("Student").Preload("Admin").Find(&tickets)

	status := c.Query("status")

	query := db.DB.Preload("Student").Preload("Admin")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&tickets)

	var ticketResponses []models.TicketResponse
	for _, t := range tickets {
		ticketResponses = append(ticketResponses, models.TicketResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			Category:    t.Category,
			StudentName: t.Student.Name,
			AdminName:   t.Admin.Name,
			StudentID: t.StudentID,
			AdminID: t.AdminID,
			
		})
	}

	c.JSON(http.StatusOK, ticketResponses)
}

// Get Ticket By ID
// @Summary Get ticket by ID
// @Description Get a ticket by its ID
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} models.TicketResponse
// @Router /api/tickets/{id} [get]
func GetTicketByID(c *gin.Context) {
	id := c.Param("id")

	var ticket models.Ticket
	if err := db.DB.Preload("Student").Preload("Admin").First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

// Update Ticket

// @Summary Update ticket
// @Description Update an existing ticket
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param ticket body models.CreateTicketRequest true "Ticket"
// @Success 200 {object} models.CreateTicketRequest
// @Router /api/tickets/{id} [put]
func UpdateTicket(c *gin.Context) {
	id := c.Param("id")

	var ticket models.Ticket
	if err := db.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	var input models.CreateTicketRequest
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		ticket.Title = input.Title
	}
	if input.Description != "" {
		ticket.Description = input.Description
	}
	if input.Status != "" {
		ticket.Status = input.Status
	}

	ticket.Category = input.Category
	ticket.StudentID = input.StudentID
	ticket.AdminID = input.AdminID
	ticket.ID = ticket.ID


	db.DB.Save(&ticket)

	c.JSON(http.StatusOK, ticket)
}

// Delete Ticket
// @Summary Delete ticket
// @Description Delete an existing ticket
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} models.CreateTicketRequest
// @Router /api/tickets/{id} [delete]
func DeleteTicket(c *gin.Context) {
	id := c.Param("id")

	var ticket models.Ticket
	if err := db.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	db.DB.Delete(&ticket)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// Assign Ticket
// @Summary Assign ticket
// @Description Assign a ticket to an admin
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param admin_id path string true "Admin ID"
// @Success 200 {object} models.CreateTicketRequest
// @Router /api/tickets/{id}/assign [post]
func AssignTicket(c *gin.Context) {
	id := c.Param("id")

	var ticket models.Ticket
	if err := db.DB.First(&ticket, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Ticket not found"})
		return
	}

	var input struct {
		AdminID uint `json:"admin_id"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var admin models.User
	if err := db.DB.First(&admin, input.AdminID).Error; err != nil || admin.Role != "admin" {
		c.JSON(400, gin.H{"error": "Invalid admin"})
		return
	}

	ticket.AdminID = input.AdminID
	ticket.Status = "in-progress"

	db.DB.Save(&ticket)

	c.JSON(200, ticket)
}
