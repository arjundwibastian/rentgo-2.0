package domain // or booking, auth, errs, etc.

import "errors"

var (
	ErrNotFound             = errors.New("data not found")
	ErrEmailRegistered      = errors.New("email already registered")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrVehicleUnavailable   = errors.New("vehicle fully booked for selected date")
	ErrUnauthorizedBooking  = errors.New("user not authorized to cancel booking")
	ErrInvalidConfirmInput  = errors.New("invalid cancel confirmation input")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrBookingAlreadyClosed = errors.New("booking completed or already cancelled")
	ErrFailedGenerateUUID = errors.New("Failed to Generate UUID")
	ErrForbidden = errors.New("Forbidden to access")
)