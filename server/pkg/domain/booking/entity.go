package booking

import (
	"errors"
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

type Booking struct {
	ID uint `json:"id"`

	Reference string `json:"reference"`

	UserID     uint   `json:"userId"`
	RoomID     uint   `json:"roomId"`
	RoomNumber string `json:"roomNumber"`

	CheckInDate  time.Time `json:"checkInDate"`
	CheckOutDate time.Time `json:"checkOutDate"`

	NightCount int    `json:"nightCount"`
	UnitPrice  int64  `json:"unitPrice"`
	TaxAmount  int64  `json:"taxAmount"`
	Discount   int64  `json:"discount"`
	TotalPrice int64  `json:"totalPrice"`
	Currency   string `json:"currency"`

	Status        BookingStatus `json:"status"`
	PaymentStatus PaymentStatus `json:"paymentStatus"`

	ExpiresAt   time.Time  `json:"expiresAt"`
	CancelledAt *time.Time `json:"cancelledAt,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Booking) ValidateDates() error {
	if b.CheckInDate.After(b.CheckOutDate) {
		return errors.New("check-in date must be before check-out date")
	}
	if b.CheckInDate.Equal(b.CheckOutDate) {
		return errors.New("check-in and check-out dates cannot be the same")
	}
	return nil
}

func (b *Booking) Cancel(){
	now := time.Now()
	b.Status = BookingStatusCancelled
	b.CancelledAt = &now
}
