package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email       string `json:"email"       validate:"required,email"`
	FullName    string `json:"fullName"    validate:"required"`
	Password    string `json:"password"    validate:"required,min=8"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=8,max=20"`
	Address     string `json:"address"     validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"fullName"`
	PhoneNumber string    `json:"phoneNumber"`
	Address     string    `json:"address"`
	Role        string    `json:"role"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}