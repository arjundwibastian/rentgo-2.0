package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	VehicleID    uuid.UUID `json:"vehicleId"    validate:"required,uuid"`
	BookingStart time.Time `json:"bookingStart" validate:"required"`
	TotalDays    int       `json:"totalDays"    validate:"required,min=1"`
}

type CancelBookingRequest struct {
	Reason *string `json:"reason"`
}

// nil => "now"
type ReturnBookingRequest struct {
	ReturnedAt *time.Time `json:"returnedAt"`
}

type BookingResponse struct {
	ID            uuid.UUID  `json:"id"`
	Status        string     `json:"status"`
	BookingStart  time.Time  `json:"bookingStart"`
	BookingEnd    time.Time  `json:"bookingEnd"`
	TotalDays     int        `json:"totalDays"`
	DailyRate     int64      `json:"dailyRate"`
	TotalPrice    int64      `json:"totalPrice"`
	HoldExpiresAt *time.Time `json:"holdExpiresAt,omitempty"`
	ReturnedAt    *time.Time `json:"returnedAt,omitempty"`
	ReturnedBy    *uuid.UUID `json:"returnedBy,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`

	Vehicle *VehicleSummary  `json:"vehicle,omitempty"`
	Payment *PaymentResponse `json:"payment,omitempty"`
}

type CreateBookingResponse struct {
	Booking BookingResponse `json:"booking"`
	Payment PaymentResponse `json:"payment"`
}

type AdminBookingResponse struct {
	BookingResponse
	User *UserResponse `json:"user,omitempty"`
}

type BookingFilter struct {
	PageQuery
	Status string    `query:"status" validate:"omitempty,oneof=pending_payment confirmed completed expired cancelled"`
	UserID uuid.UUID `query:"userId" validate:"omitempty,uuid"`
}