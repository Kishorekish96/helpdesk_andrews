package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
}

type Ticket struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Category    string `json:"category"`

	StudentID uint `json:"student_id"`
	Student   User `gorm:"foreignKey:StudentID"`

	AdminID uint `json:"admin_id"`
	Admin   User `gorm:"foreignKey:AdminID"`
}


type CreateUserRequest struct {
	Name       string `json:"name" example:"John"`
	Email      string `json:"email" example:"john@gmail.com"`
	Role       string `json:"role" example:"student"`
	Department string `json:"department" example:"IT"`
}

type CreateTicketRequest struct {
	Title       string `json:"title" example:"WiFi issue"`
	Description string `json:"description" example:"Internet slow"`
	Category    string `json:"category" example:"IT"`
	StudentID   uint   `json:"student_id" example:"1"`
	Status      string `json:"status" example:"open"`
	AdminID     uint   `json:"admin_id" example:"0"`
}


type TicketResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	StudentName string `json:"student_name"`
	AdminName   string `json:"admin_name"`
	StudentID   uint   `json:"student_id"`
	AdminID     uint   `json:"admin_id"`
}

type ListUser struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
}