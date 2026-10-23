package dto

import (
	"time"

	"github.com/google/uuid"
)

type PaymentResponse struct {
	ID       uuid.UUID  `json:"id"`
	BookingID uuid.UUID `json:"bookingId"`
	Amount   int64      `json:"amount"`
	Currency string     `json:"currency"`
	Kind     string     `json:"kind"`
	Status   string     `json:"status"`
	Provider string     `json:"provider"`
	ExpiresAt time.Time `json:"expiresAt"`
	PaidAt    *time.Time `json:"paidAt,omitempty"`

	// Provider-supplied instructions; not columns.
	PaymentURL *string `json:"paymentUrl,omitempty"`
	VANumber   *string `json:"vaNumber,omitempty"`
	QRString   *string `json:"qrString,omitempty"`
}

type RefundResponse struct {
	ID          uuid.UUID `json:"id"`
	PaymentID   uuid.UUID `json:"paymentId"`
	Amount      int64     `json:"amount"`
	Reason      *string   `json:"reason,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}