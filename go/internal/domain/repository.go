package domain

import (
	"context"

	"uuid"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	Register(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}


type EmailService interface {
    SendBookingConfirmation(booking *Booking, vehicle *Vehicle, user *User) error
    SendBookingCancellation(booking *Booking, user *User) error 
}