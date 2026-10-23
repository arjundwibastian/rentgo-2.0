package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateVehicleRequest struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description" validate:"required"`
	Category    string `json:"category"    validate:"required"`
	Quantity    int    `json:"quantity"    validate:"required,min=0"`
	DailyRate   int64  `json:"dailyRate"   validate:"required,gt=0"`
}

// Pointers so PATCH can distinguish "not sent" from "set to zero".
type UpdateVehicleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Quantity    *int    `json:"quantity"    validate:"omitempty,min=0"`
	DailyRate   *int64  `json:"dailyRate"   validate:"omitempty,gt=0"`
	IsActive    *bool   `json:"isActive"`
}

type VehicleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Quantity    int       `json:"quantity"`
	DailyRate   int64     `json:"dailyRate"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type VehicleSummary struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
}

type VehicleAvailabilityResponse struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Category          string    `json:"category"`
	DailyRate         int64     `json:"dailyRate"`
	AvailableQuantity int       `json:"availableQuantity"`
}